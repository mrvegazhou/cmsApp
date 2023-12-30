package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppUserDao struct {
	DB *gorm.DB
}

var (
	instanceUser *AppUserDao
	onceUserDao  sync.Once
)

func NewAppUserDao() *AppUserDao {
	onceUserDao.Do(func() {
		instanceUser = &AppUserDao{DB: postgresqlx.GetDB(&models.AppUser{})}
	})
	return instanceUser
}

func (dao *AppUserDao) GetAppUser(conditions map[string]interface{}) (user models.AppUser, err error) {
	err = dao.DB.Where(conditions).First(&user).Error
	return
}

func (dao *AppUserDao) UpdateColumns(conditions, field map[string]interface{}, tx *gorm.DB) error {
	if tx != nil {
		return tx.Model(&models.AppUser{}).Where(conditions).UpdateColumns(field).Error
	}
	return dao.DB.Model(&models.AppUser{}).Where(conditions).UpdateColumns(field).Error
}

func (dao *AppUserDao) CreateAppUser(user models.AppUser) (uint64, error) {
	if err := dao.DB.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}
