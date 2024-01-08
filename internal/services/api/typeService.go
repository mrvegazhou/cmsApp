package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"sync"
)

type apiTypeService struct {
	Dao *dao.AppTypeDao
}

var (
	instanceApiTypeService *apiTypeService
	onceApiTypeService     sync.Once
)

func NewApiTypeService() *apiTypeService {
	onceApiTypeService.Do(func() {
		instanceApiTypeService = &apiTypeService{
			Dao: dao.NewAppTypeDao(),
		}
	})
	return instanceApiTypeService
}

func (ser *apiTypeService) GetTypeList(name string) (typeList []models.AppType, err error) {
	var condition map[string][]interface{}
	if name != "" {
		condition = map[string][]interface{}{
			"name": []interface{}{"like ?", "%" + name + "%"},
		}
	}
	typeList, err = ser.Dao.GetTypeList(condition)
	if err == gorm.ErrRecordNotFound {
		return typeList, nil
	}
	return typeList, err
}

func (ser *apiTypeService) GetTypeListByPid(pid uint64) (typeList []models.AppType, err error) {
	condition := map[string][]interface{}{
		"pid": []interface{}{"= ?", pid},
	}
	typeList, err = ser.Dao.GetTypeList(condition)
	if err == gorm.ErrRecordNotFound {
		return typeList, nil
	}
	return typeList, err
}

func (ser *apiTypeService) GetTypeInfoById(id uint64) (typeList []models.AppType, err error) {
	condition := map[string][]interface{}{
		"id": []interface{}{"= ?", id},
	}
	typeList, err = ser.Dao.GetTypeList(condition)
	if err == gorm.ErrRecordNotFound {
		return typeList, nil
	}
	return typeList, err
}
