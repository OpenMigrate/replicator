package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"replicator/internal/models"
)

type Store struct {
	DB *gorm.DB
}

// Init opens SQLite and runs migrations.
// dbURL examples:
//
//	file:replicator.db?cache=shared&_busy_timeout=5000
//	:memory:
func Init(dbUrl string) (*Store, error) {
	cfg := &gorm.Config{}
	db, err := gorm.Open(sqlite.Open(dbUrl), cfg)
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Metadata{},
		&models.App{},
		&models.AppServer{},
	); err != nil {
		return nil, err
	}

	if err := db.SetupJoinTable(&models.App{}, "Servers", &models.AppServer{}); err != nil {
		return nil, err
	}
	if err := db.SetupJoinTable(&models.Metadata{}, "Apps", &models.AppServer{}); err != nil {
		return nil, err
	}

	return &Store{DB: db}, nil
}

func (s *Store) SaveServer(md models.Metadata) error {
	return s.DB.Create(&md).Error
}

func (s *Store) ListServers() (res []models.Metadata, err error) {
	err = s.DB.Find(&res).Error
	return
}

func (s *Store) GetServer(id string) (models.Metadata, error) {
	var md models.Metadata
	return md, s.DB.First(&md, "id = ?", id).Error
}

func (s *Store) DeleteServer(id string) error {
	return s.DB.Delete(&models.Metadata{}, "id = ?", id).Error
}
