// storage/apps_store.go
package storage

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"

	"replicator/internal/models"
)

func (s *Store) CreateApp(in AppCreate) (*models.App, error) {
	app := &models.App{
		ID:          in.ID,
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.DB.Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// DeleteApp deletes the app row and detaches all related servers (rows in app_servers).
func (s *Store) DeleteApp(sel AppSelector) error {
	app, err := s.FindApp(sel)
	if err != nil {
		return err
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		// detach memberships
		if err := tx.Where("app_id = ?", app.ID).Delete(&models.AppServer{}).Error; err != nil {
			return err
		}
		// delete the app itself
		res := tx.Delete(&models.App{}, "id = ?", app.ID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (s *Store) FindApp(sel AppSelector) (*models.App, error) {
	var app models.App
	tx := s.DB.Model(&models.App{})
	switch {
	case sel.ID != nil && *sel.ID != "":
		if err := tx.First(&app, "id = ?", *sel.ID).Error; err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("empty selector")
	}
	return &app, nil
}

func (s *Store) ListApps(afterID string, limit int) ([]models.App, string, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	q := s.DB.Model(&models.App{})
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

func (s *Store) ModifyAppServers(sel AppSelector, serverIDs []string, op MembershipOp) error {
	app, err := s.FindApp(sel)
	if err != nil {
		return err
	}
	switch op {
	case MembershipAdd:
		return s.addAppServers(app.ID, serverIDs)
	case MembershipRemove:
		return s.removeAppServers(app.ID, serverIDs)
	default:
		return errors.New("invalid membership op")
	}
}

func (s *Store) ListAppServers(sel AppSelector, cur Cursor) ([]models.Metadata, int64, string, error) {
	app, err := s.FindApp(sel)
	if err != nil {
		return nil, 0, "", err
	}
	if cur.Limit <= 0 || cur.Limit > 500 {
		cur.Limit = 50
	}
	var total int64
	if err := s.DB.Model(&models.AppServer{}).Where("app_id = ?", app.ID).Count(&total).Error; err != nil {
		return nil, 0, "", err
	}
	sub := s.DB.Model(&models.AppServer{}).Select("metadata_id").Where("app_id = ?", app.ID)
	q := s.DB.Model(&models.Metadata{}).Where("id IN (?)", sub)
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

func (s *Store) addAppServers(appID string, serverIDs []string) error {
	if len(serverIDs) == 0 {
		return nil
	}
	serverIDs = unique(serverIDs)
	var existing []string
	if err := s.DB.Model(&models.Metadata{}).Where("id IN ?", serverIDs).Pluck("id", &existing).Error; err != nil {
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
	return s.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "app_id"}, {Name: "metadata_id"}},
		DoNothing: true,
	}).CreateInBatches(&links, 500).Error
}

func (s *Store) removeAppServers(appID string, serverIDs []string) error {
	if len(serverIDs) == 0 {
		return nil
	}
	serverIDs = unique(serverIDs)
	return s.DB.Where("app_id = ? AND metadata_id IN ?", appID, serverIDs).Delete(&models.AppServer{}).Error
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
