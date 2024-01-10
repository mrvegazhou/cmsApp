package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type apiTypeService struct {
	TypeDao *dao.AppTypeDao
}

var (
	instanceApiTypeService *apiTypeService
	onceApiTypeService     sync.Once
)

func NewApiTypeService() *apiTypeService {
	onceApiTypeService.Do(func() {
		instanceApiTypeService = &apiTypeService{
			TypeDao: dao.NewAppTypeDao(),
		}
	})
	return instanceApiTypeService
}

func (ser *apiTypeService) GetTypeList(name string) (typeList []models.AppType, err error) {
	conditions := map[string][]interface{}{}
	name = strings.TrimSpace(name)
	if name != "" {
		conditions = map[string][]interface{}{
			"name": {"like ?", "%" + name + "%"},
		}
	}
	typeList, err = ser.TypeDao.GetTypeList(conditions)
	if err == gorm.ErrRecordNotFound {
		return typeList, nil
	}
	return typeList, err
}

func (ser *apiTypeService) GetTypeListByPid(pid uint64) (typeList []models.AppType, err error) {
	condition := map[string][]interface{}{
		"pid": []interface{}{"= ?", pid},
	}
	typeList, err = ser.TypeDao.GetTypeList(condition)
	if err == gorm.ErrRecordNotFound {
		return typeList, nil
	}
	return typeList, err
}

func (ser *apiTypeService) GetTypeInfoById(id uint64) (typeInfo models.AppTypeByIdInfo, err error) {
	condition := map[string][]interface{}{
		"id":     []interface{}{"= ?", id},
		"status": []interface{}{"= ?", 1},
	}
	appType := models.AppType{}
	appType, err = ser.TypeDao.GetTypeInfo(condition)
	if err == gorm.ErrRecordNotFound {
		return typeInfo, nil
	}
	typeInfo.Id = appType.Id
	typeInfo.Name = appType.Name

	if appType.Pid != 0 {
		var appPType models.AppType
		condition = map[string][]interface{}{
			"id":     []interface{}{"= ?", appType.Pid},
			"status": []interface{}{"= ?", 1},
		}
		appPType, err = ser.TypeDao.GetTypeInfo(condition)
		if err == nil {
			typeInfo.Pid = appPType.Id
			typeInfo.Pname = appPType.Name
		}
	}
	return typeInfo, err
}
