package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
)

type UpcomingAlbumRepository struct {
	*gorm.DB
}

func MakeUpcomingAlbumRepository() *UpcomingAlbumRepository {
	return &UpcomingAlbumRepository{PostgresDB}
}

func (repo *UpcomingAlbumRepository) FindByAlbumId(albumId uint) *models.UpcomingAlbum {
	upcomingAlbum := models.UpcomingAlbum{AlbumID: albumId}

	repo.Where(&upcomingAlbum).First(&upcomingAlbum)
	return &upcomingAlbum
}
