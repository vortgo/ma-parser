package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/elasticsearch"
	"github.com/vortgo/ma-parser/models"
)

type AlbumRepository struct {
	*gorm.DB
}

func MakeAlbumRepository() *AlbumRepository {
	return &AlbumRepository{PostgresDB}
}

func (repo *AlbumRepository) FindAlbumByPlatformId(platformId int) *models.Album {
	album := models.Album{PlatformID: platformId}

	repo.Where(&album).First(&album)
	return &album
}

func (repo *AlbumRepository) SaveToElastic(album models.Album) {
	var band models.Band
	repo.DB.Where("id = ?", album.BandID).First(&band)
	elasticDocument := album.GetIndexJson(band)
	elasticsearch.IndexDataToElastic(&album, elasticDocument)
}
