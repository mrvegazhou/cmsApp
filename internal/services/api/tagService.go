package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"strings"
	"sync"
)

type apiTagService struct {
	TagDao *dao.AppTagDao
}

var (
	instanceApiTagService *apiTagService
	onceApiTagService     sync.Once
)

func NewApiTagService() *apiTagService {
	onceApiTagService.Do(func() {
		instanceApiTagService = &apiTagService{
			TagDao: dao.NewAppTagDao(),
		}
	})
	return instanceApiTagService
}

func (ser *apiTagService) GetTagList(name string) (appTagList []models.AppTag, err error) {
	conditions := map[string][]interface{}{}
	name = strings.TrimSpace(name)
	if name == "" {
		// 应该返回热门的tag列表
		return appTagList, nil
	} else {
		conditions = map[string][]interface{}{
			"name": {"like ?", "%" + name + "%"},
		}
		appTagList, err = ser.TagDao.GetTagList(conditions)
		if err == gorm.ErrRecordNotFound {
			return appTagList, nil
		}
		return appTagList, err
	}
}

func (ser *apiTagService) GetTagListByIds(ids []uint64) (appTagList []models.AppTag, err error) {
	conditions := map[string][]interface{}{}
	if len(ids) == 0 {
		return appTagList, nil
	} else {
		conditions = map[string][]interface{}{
			"id": {"in (?)", ids},
		}
		appTagList, err = ser.TagDao.GetTagList(conditions)
		if err == gorm.ErrRecordNotFound {
			return appTagList, nil
		}
		return appTagList, err
	}
}

func (ser *apiTagService) GetTagInfo(id uint64) (appTagInfo models.AppTag, err error) {
	return ser.TagDao.GetTagInfo(map[string]interface{}{"id": id})
}
