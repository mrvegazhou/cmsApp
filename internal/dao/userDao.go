package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"cmsApp/pkg/utils/stringx"
	"gorm.io/gorm"
	"sync"
)

type AppUserDao struct {
	DB *gorm.DB
	BaseDao
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

func (dao *AppUserDao) SearchUserList(name string, pageParam int, pageSizeParam int) ([]models.AppUser, int, int, error) {
	users := []models.AppUser{}
	if name == "" {
		return users, 0, 0, nil
	}
	Db := dao.DB.Model(&models.AppUser{})
	if stringx.CheckEmail(name) {
		Db = Db.Where("email=?", name)
	} else if stringx.CheckMobile(name) {
		Db = Db.Where("phone=?", name)
	} else {
		name = "%" + name + "%"
		Db = Db.Where("nickname like ?", name)
	}

	total, err := dao.SearchUserListTotal(name)
	if err != nil {
		return users, 1, 0, err
	}

	page, totalPage, pageSize, offset := dao.Page(pageParam, pageSizeParam, total)
	Db = Db.Scopes(dao.Order("id desc")).Offset(offset).Limit(pageSize)
	if err := Db.Find(&users).Error; err != nil {
		return users, page, totalPage, err
	}
	return users, page, totalPage, nil
}

func (dao *AppUserDao) GetUserList(conditions map[string][]interface{}) ([]models.AppUser, error) {
	appUser := []models.AppUser{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.UserFields{})
	}
	Db = Db.Scopes(dao.Order("id desc"))
	if err := Db.Find(&appUser).Error; err != nil {
		return appUser, err
	}
	return appUser, nil
}

func (dao *AppUserDao) SearchUserListTotal(name string) (int64, error) {
	Db := dao.DB.Model(&models.AppUser{})
	if stringx.CheckEmail(name) {
		Db = Db.Where("email=?", name)
	} else if stringx.CheckMobile(name) {
		Db = Db.Where("phone=?", name)
	} else {
		name = "%" + name + "%"
		Db = Db.Where("nickname like ?", name)
	}
	var count int64
	err := Db.Count(&count).Error
	return count, err
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
