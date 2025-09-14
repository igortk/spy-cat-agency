package models

import (
	"github.com/lib/pq"
	"time"
)

type Target struct {
	ID          string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	MissionID   string         `json:"missionId" gorm:"type:uuid;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Country     string         `json:"country" gorm:"not null"`
	Notes       pq.StringArray `json:"notes" gorm:"type:jsonb"`
	IsCompleted bool           `json:"isCompleted" gorm:"not null;default:false"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}
