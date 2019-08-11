package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
	"strings"
)

type GenreRepository struct {
	*gorm.DB
}

func MakeGenreRepository() *GenreRepository {
	return &GenreRepository{PostgresDB}
}

func (repo *GenreRepository) FindOrCreatGenreByName(name string) *models.Genre {
	genre := models.Genre{}

	genre.Name = strings.ToUpper(name)
	repo.Where(&genre).FirstOrCreate(&genre)
	return &genre
}
