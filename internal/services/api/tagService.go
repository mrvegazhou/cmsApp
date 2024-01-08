package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"sync"
)

type apiTagService struct {
	Dao *dao.AppTagDao
}

var (
	instanceApiTagService *apiTagService
	onceApiTagService     sync.Once
)

func NewApiTagService() *apiTagService {
	onceApiTypeService.Do(func() {
		instanceApiTagService = &apiTagService{
			Dao: dao.NewAppTagDao(),
		}
	})
	return instanceApiTagService
}

func (ser *apiTagService) GetTagList(name string) (appTagList []models.AppTag, err error) {
	condition := map[string][]interface{}{
		"name": []interface{}{"like ?", "%" + name + "%"},
	}
	appTagList, err = ser.Dao.GetTagList(condition)
	if err == gorm.ErrRecordNotFound {
		return appTagList, nil
	}
	return appTagList, err
}
