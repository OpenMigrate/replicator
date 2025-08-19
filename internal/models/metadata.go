package models

import "time"

// --- servers (metadata) ---
type Metadata struct {
	ID              string `json:"id" gorm:"primaryKey;Size:64;not null"`
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Arch            string `json:"arch"`
	NumCPU          int    `json:"num_cpu"`
	Kernel          string `json:"kernel"`
	Uptime          string `json:"uptime"`
	TotalMemoryMB   uint64 `json:"total_memory_mb"`
	TotalDiskSizeGB string `json:"total_disk_size_gb"`
	MountedCount    int    `json:"mounted_count"`
	TimestampUTC    string `json:"timestamp_utc" gorm:"index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Apps []App `json:"apps" gorm:"many2many:app_servers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Apps []App `json:"apps" gorm:"many2many:app_servers"`
}
