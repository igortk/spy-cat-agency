package models

import (
	"time"
)

type Mission struct {
	ID          string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CatID       string    `json:"catId" gorm:"type:uuid;unique"`
	IsCompleted bool      `json:"isCompleted" gorm:"not null;default:false"`
	Targets     []Target  `json:"targets" gorm:"foreignKey:MissionID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
