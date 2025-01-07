package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Keyboard struct {
	ID uuid.UUID `gorm:"type:uuid"`
	Username string `gorm:"unique;not null" json:"username"`
	Keypress uint64 `gorm:"default:0" json:"keypress"`
}

func (keyboard * Keyboard) BeforeCreate(tx *gorm.DB) (err error) {
	keyboard.ID = uuid.New()
	return
}
