package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
	"strings"
)

type LyricalThemeRepository struct {
	*gorm.DB
}

func MakeLyricalThemeRepository() *LyricalThemeRepository {
	return &LyricalThemeRepository{PostgresDB}
}

func (repo *LyricalThemeRepository) FindOrCreatLyricalThemeByName(name string) *models.LyricalTheme {
	lyricalTheme := models.LyricalTheme{}

	lyricalTheme.Name = strings.ToUpper(name)
	repo.Where(&lyricalTheme).FirstOrCreate(&lyricalTheme)
	return &lyricalTheme
}
