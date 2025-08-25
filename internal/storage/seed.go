package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	"replicator/internal/models"
)

// SeedSampleData inserts sample apps and servers for testing.
func (s *Store) SeedSampleData(ctx context.Context) error {
	now := time.Now()

	// Sample servers
	servers := []models.Metadata{
		{
			ID:              uuid.NewString(),
			Hostname:        "srv-payments-01",
			OS:              "linux",
			Arch:            "amd64",
			NumCPU:          4,
			Kernel:          "5.15.0",
			Uptime:          "12h",
			TotalMemoryMB:   8192,
			TotalDiskSizeGB: "100",
			MountedCount:    3,
			TimestampUTC:    now.UTC().Format(time.RFC3339),
		},
		{
			ID:              uuid.NewString(),
			Hostname:        "srv-analytics-01",
			OS:              "linux",
			Arch:            "amd64",
			NumCPU:          8,
			Kernel:          "5.15.0",
			Uptime:          "3h",
			TotalMemoryMB:   16384,
			TotalDiskSizeGB: "200",
			MountedCount:    4,
			TimestampUTC:    now.UTC().Format(time.RFC3339),
		},
	}

	// Insert servers
	if err := s.DB.Create(&servers).Error; err != nil {
		return err
	}

	// Sample apps
	apps := []models.App{
		{
			ID:          uuid.NewString(),
			Name:        "Payments",
			Description: "Prod payment service",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.NewString(),
			Name:        "Analytics",
			Description: "Analytics and BI service",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	if err := s.DB.Create(&apps).Error; err != nil {
		return err
	}

	// Link first server to Payments app
	appServer := models.AppServer{
		AppID:      apps[0].ID,
		MetadataID: servers[0].ID,
	}
	if err := s.DB.Create(&appServer).Error; err != nil {
		return err
	}

	return nil
}
