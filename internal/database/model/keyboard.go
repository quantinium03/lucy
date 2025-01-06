package model

import "github.com/google/uuid"

type Keyboard struct {
	ID uuid.UUID `gorm:"type:uuid"`
	Username string `gorm:"unique;not null;primaryKey" json:"username"`
	Keypress uint64 `json:"keypress"`
}
