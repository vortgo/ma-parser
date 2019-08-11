package models

import (
	"time"
)

type LatestBandUpdate struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	BandID    uint
	Band      Band
}
