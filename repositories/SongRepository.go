package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
)

type SongRepository struct {
	*gorm.DB
}

func MakeSongRepository() *SongRepository {
	return &SongRepository{PostgresDB}
}

func (repo *SongRepository) LoadByPlatformId(platformId int) *models.Song {
	song := models.Song{PlatformID: platformId}

	repo.Where(&song).FirstOrInit(&song)
	return &song
}
