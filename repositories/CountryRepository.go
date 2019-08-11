package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
	"strings"
)

type CountryRepository struct {
	*gorm.DB
}

func MakeCountryRepository() *CountryRepository {
	return &CountryRepository{PostgresDB}
}

func (repo *CountryRepository) FindOrCreatCountryByName(name string) *models.Country {
	country := models.Country{Name: strings.ToUpper(name)}
	repo.Where(&country).FirstOrCreate(&country)
	return &country
}
