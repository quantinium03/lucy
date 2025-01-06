package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Coding struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid"`
	Username         string    `gorm:"unique;not null;primaryKey" json:"username"`
	TimeCodingDaily  time.Time `gorm:"not null" json:"timeCodingDaily"`
	TimeCodingWeekly time.Time `gorm:"not null" json:"timeCodingWeekly"`
}

type CodingActivity struct {
	gorm.Model
	ID                 uuid.UUID  `gorm:"type:uuid"`
	Date               string     `gorm:"type:date;not null;unique" json:"date"`
	TotalMinutePerDate int        `gorm:"not null" json:"totalMinutesPerDate"`
	Languages          []Language `json:"languages"`
}

type Language struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid"`
	LanguageName string    `json:"language"`
}
