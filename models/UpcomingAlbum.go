package models

import (
	"time"
)

type UpcomingAlbum struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AlbumID   uint
	Album     Album
}
