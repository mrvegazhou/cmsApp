package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"sync"
)

type apiArticleTypeService struct {
	Dao *dao.AppArticleTypeDao
}

var (
	instanceApiArticleTypeService *apiArticleTypeService
	onceApiArticleTypeService     sync.Once
)

func NewApiArticleTypeService() *apiArticleTypeService {
	onceApiArticleTypeService.Do(func() {
		instanceApiArticleTypeService = &apiArticleTypeService{
			Dao: dao.NewAppArticleTypeDao(),
		}
	})
	return instanceApiArticleTypeService
}

func (ser *apiArticleTypeService) GetArticleTypeList(name string) (articleTypeList []models.AppArticleType, err error) {
	condition := map[string][]interface{}{
		"name": []interface{}{"like ?", name + "%"},
	}
	articleTypeList, err = ser.Dao.GetArticleTypeList(condition)
	if err == gorm.ErrRecordNotFound {
		return articleTypeList, nil
	}
	return articleTypeList, err
}
