package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Mouse struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid"`
	Username    string    `gorm:"unique;not null" json:"username"`
	LeftClick   uint64    `gorm:"default:0" json:"leftClick"`
	RightClick  uint64    `gorm:"default:0" json:"rightClick"`
	MouseTravel float64   `gorm:"default:0" json:"mouseTravel"`
}

func (mouse *Mouse) BeforeCreate(tx *gorm.DB) (err error) {
	mouse.ID = uuid.New()
	return
}
