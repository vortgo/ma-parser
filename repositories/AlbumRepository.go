package repositories

import (
	"github.com/vortgo/ma-parser/models"
)

type AlbumRepository struct {
	*DbAndElasticRepository
}

func MakeAlbumRepository() *AlbumRepository {
	return &AlbumRepository{&DbAndElasticRepository{PostgresDB}}
}

func (repo *AlbumRepository) FindAlbumByPlatformId(platformId int) *models.Album {
	album := models.Album{PlatformID: platformId}

	repo.Where(&album).First(&album)
	return &album
}
