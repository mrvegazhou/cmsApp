package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppTypeDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceAppTypeDao *AppTypeDao
	onceAppTypeDao     sync.Once
)

func NewAppTypeDao() *AppTypeDao {
	onceAppTypeDao.Do(func() {
		instanceAppTypeDao = &AppTypeDao{DB: postgresqlx.GetDB(&models.AppType{})}
	})
	return instanceAppTypeDao
}

func (dao *AppTypeDao) CreateAppType(appType *models.AppType) (uint64, error) {
	if err := dao.DB.Create(appType).Error; err != nil {
		return 0, err
	}
	return appType.Id, nil
}

func (dao *AppTypeDao) GetTypeList(conditions map[string][]interface{}) ([]models.AppType, error) {
	appTypes := []models.AppType{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.BaseDao.ConditionWhere(Db, conditions, models.AppTypeFields{})
	}
	Db = Db.Scopes(dao.Order("create_time desc"))
	if err := Db.Find(&appTypes).Error; err != nil {
		return appTypes, err
	}
	return appTypes, nil
}

func (dao *AppTypeDao) GetTypeInfo(conditions map[string][]interface{}) (models.AppType, error) {
	appType := models.AppType{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.BaseDao.ConditionWhere(Db, conditions, models.AppTypeFields{})
	}
	Db = Db.Scopes(dao.Order("create_time desc"))
	if err := Db.First(&appType).Error; err != nil {
		return appType, err
	}
	return appType, nil
}
