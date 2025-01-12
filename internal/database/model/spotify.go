package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Spotify struct {
	gorm.Model
	ID                    uuid.UUID `gorm:"type:uuid"`
	Username              string    `gorm:"unique;not null" json:"username"`
	SpotifyTrackEmbedURI  string    `gorm:"unique" json:"spotifyTrackEmbedUri"`
	SpotifyTrackEmbedHtml string    `gorm:"unique" json:"spotifyTrackEmbedHtml"`
	SpotifyAccessToken    string    `gorm:"unique" json:"SpotifyAccessToken"`
}

func (spotify *Spotify) BeforeCreate(tx *gorm.DB) (err error) {
	spotify.ID = uuid.New()
	return
}
