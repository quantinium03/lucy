package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Mouse struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid"`
	Username string `gorm:"unique;not null;primaryKey" json:"username"`
	LeftClick uint64 `json:"leftClick"`	
	RightClick uint64 `json:"rightClick"`
	MouseTravel uint64 `json:"mouseTravel"`
}
