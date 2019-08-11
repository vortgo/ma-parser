package repositories

import (
	"github.com/vortgo/ma-parser/models"
)

type SongRepository struct {
	*DbAndElasticRepository
}

func MakeSongRepository() *SongRepository {
	return &SongRepository{&DbAndElasticRepository{PostgresDB}}
}

func (repo *SongRepository) LoadByPlatformId(platformId int) *models.Song {
	song := models.Song{PlatformID: platformId}

	repo.Where(&song).FirstOrInit(&song)
	return &song
}
