package models

import "time"

// --- apps ---
type App struct {
	ID          string `json:"id" gorm:"primaryKey;size:64;not null"`
	Name        string `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Servers []Metadata `json:"servers" gorm:"many2many:app_servers"`
}

type AppServer struct {
	AppID      string `json:"app_id" gorm:"size:64;not null;primaryKey;column:app_id"`
	MetadataID string `json:"metadata_id" gorm:"size:64;not null;primaryKey;column:metadata_id"`
	CreatedAt  time.Time

	// Optional FKs (good for cascades)
	App      App      `gorm:"foreignKey:AppID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Metadata Metadata `gorm:"foreignKey:MetadataID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
