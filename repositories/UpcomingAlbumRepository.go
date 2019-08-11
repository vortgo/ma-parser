package repositories

import (
	"github.com/jinzhu/gorm"
)

type UpcomingAlbumRepository struct {
	*gorm.DB
}

func MakeUpcomingAlbumRepository() *UpcomingAlbumRepository {
	return &UpcomingAlbumRepository{PostgresDB}
}
