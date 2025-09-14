package models

import (
	"time"
)

type Cat struct {
	ID                string    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name              string    `json:"name" gorm:"not null"`
	YearsOfExperience int       `json:"yearsOfExperience" gorm:"not null"`
	Breed             string    `json:"breed" gorm:"not null"`
	Salary            float64   `json:"salary" gorm:"not null"`
	Status            string    `json:"status" gorm:"not null;default:'Available'"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
