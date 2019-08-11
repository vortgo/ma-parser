package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
	"strings"
)

type LabelRepository struct {
	*gorm.DB
}

func MakeLabelRepository() *LabelRepository {
	return &LabelRepository{PostgresDB}
}

func (repo *LabelRepository) FindOrCreateLabelByName(name string) *models.Label {
	label := models.Label{Name: strings.ToUpper(name)}
	repo.Where(&label).FirstOrCreate(&label)
	return &label
}
