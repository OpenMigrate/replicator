// storage/apps_store.go
package storage

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"replicator/internal/models"
)

const defaultTimeout = 5 * time.Second

func withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultTimeout)
}

func (s *Store) CreateApp(ctx context.Context, in AppCreate) (*models.App, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	app := &models.App{
		ID:          in.ID,
		Name:        in.Name,
		Identifier:  in.Identifier,
		Description: in.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.DB.WithContext(ctx).Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

func (s *Store) FindApp(ctx context.Context, sel AppSelector) (*models.App, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	var app models.App
	tx := s.DB.WithContext(ctx).Model(&models.App{})
	switch {
	case sel.ID != nil && *sel.ID != "":
		if err := tx.First(&app, "id = ?", *sel.ID).Error; err != nil {
			return nil, err
		}
	case sel.Identifier != nil && *sel.Identifier != "":
		if err := tx.First(&app, "identifier = ?", *sel.Identifier).Error; err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("empty selector")
	}
	return &app, nil
}

func (s *Store) ListApps(ctx context.Context, afterID string, limit int) ([]models.App, string, error) {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	q := s.DB.WithContext(ctx).Model(&models.App{})
	if afterID != "" {
		q = q.Where("id > ?", afterID)
	}
	var apps []models.App
	if err := q.Order("id ASC").Limit(limit).Find(&apps).Error; err != nil {
		return nil, "", err
	}
	var next string
	if len(apps) == limit {
		next = apps[len(apps)-1].ID
	}
	return apps, next, nil
}

func (s *Store) ModifyAppServers(ctx context.Context, sel AppSelector, serverIDs []string, op MembershipOp) error {
	app, err := s.FindApp(ctx, sel)
	if err != nil {
		return err
	}
	switch op {
	case MembershipAdd:
		return s.addAppServers(ctx, app.ID, serverIDs)
	case MembershipRemove:
		return s.removeAppServers(ctx, app.ID, serverIDs)
	case MembershipReplace:
		return s.replaceAppServersDelta(ctx, app.ID, serverIDs)
	default:
		return errors.New("invalid membership op")
	}
}

func (s *Store) ListAppServers(ctx context.Context, sel AppSelector, cur Cursor) ([]models.Metadata, int64, string, error) {
	app, err := s.FindApp(ctx, sel)
	if err != nil {
		return nil, 0, "", err
	}
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	if cur.Limit <= 0 || cur.Limit > 500 {
		cur.Limit = 50
	}
	var total int64
	if err := s.DB.WithContext(ctx).Model(&models.AppServer{}).Where("app_id = ?", app.ID).Count(&total).Error; err != nil {
		return nil, 0, "", err
	}
	sub := s.DB.WithContext(ctx).Model(&models.AppServer{}).Select("metadata_id").Where("app_id = ?", app.ID)
	q := s.DB.WithContext(ctx).Model(&models.Metadata{}).Where("id IN (?)", sub)
	if cur.AfterID != "" {
		q = q.Where("id > ?", cur.AfterID)
	}
	var servers []models.Metadata
	if err := q.Order("id ASC").Limit(cur.Limit).Find(&servers).Error; err != nil {
		return nil, 0, "", err
	}
	var next string
	if len(servers) == cur.Limit {
		next = servers[len(servers)-1].ID
	}
	return servers, total, next, nil
}

func (s *Store) addAppServers(ctx context.Context, appID string, serverIDs []string) error {
	if len(serverIDs) == 0 {
		return nil
	}
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	serverIDs = unique(serverIDs)
	var existing []string
	if err := s.DB.WithContext(ctx).Model(&models.Metadata{}).Where("id IN ?", serverIDs).Pluck("id", &existing).Error; err != nil {
		return err
	}
	if len(existing) == 0 {
		return nil
	}
	now := time.Now()
	links := make([]models.AppServer, 0, len(existing))
	for _, sid := range existing {
		links = append(links, models.AppServer{AppID: appID, MetadataID: sid, CreatedAt: now})
	}
	tx := s.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "app_id"}, {Name: "metadata_id"}},
		DoNothing: true,
	})
	const batch = 500
	return tx.CreateInBatches(&links, batch).Error
}

func (s *Store) removeAppServers(ctx context.Context, appID string, serverIDs []string) error {
	if len(serverIDs) == 0 {
		return nil
	}
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	serverIDs = unique(serverIDs)
	return s.DB.WithContext(ctx).Where("app_id = ? AND metadata_id IN ?", appID, serverIDs).Delete(&models.AppServer{}).Error
}

func (s *Store) replaceAppServersDelta(ctx context.Context, appID string, desired []string) error {
	ctx, cancel := withTimeout(ctx)
	defer cancel()
	desired = unique(desired)
	var current []string
	if err := s.DB.WithContext(ctx).Model(&models.AppServer{}).Where("app_id = ?", appID).Pluck("metadata_id", &current).Error; err != nil {
		return err
	}
	var existing []string
	if len(desired) > 0 {
		if err := s.DB.WithContext(ctx).Model(&models.Metadata{}).Where("id IN ?", desired).Pluck("id", &existing).Error; err != nil {
			return err
		}
	}
	addSet := diff(desired, current)
	if len(addSet) > 0 {
		addSet = intersect(addSet, existing)
	}
	delSet := diff(current, desired)
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(delSet) > 0 {
			if err := tx.Where("app_id = ? AND metadata_id IN ?", appID, delSet).Delete(&models.AppServer{}).Error; err != nil {
				return err
			}
		}
		if len(addSet) > 0 {
			now := time.Now()
			links := make([]models.AppServer, 0, len(addSet))
			for _, sid := range addSet {
				links = append(links, models.AppServer{AppID: appID, MetadataID: sid, CreatedAt: now})
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "app_id"}, {Name: "metadata_id"}},
				DoNothing: true,
			}).Create(&links).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func unique(in []string) []string {
	m := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if v == "" {
			continue
		}
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func diff(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, v := range b {
		mb[v] = struct{}{}
	}
	out := make([]string, 0, len(a))
	for _, v := range a {
		if _, ok := mb[v]; !ok && v != "" {
			out = append(out, v)
		}
	}
	return out
}

func intersect(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, v := range b {
		mb[v] = struct{}{}
	}
	out := make([]string, 0, len(a))
	for _, v := range a {
		if _, ok := mb[v]; ok {
			out = append(out, v)
		}
	}
	return out
}
