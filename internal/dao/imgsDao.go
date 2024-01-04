package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type ImgsDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceImgsDao *ImgsDao
	onceImgsDao     sync.Once
)

func NewImgsDao() *ImgsDao {
	onceImgsDao.Do(func() {
		instanceImgsDao = &ImgsDao{DB: postgresqlx.GetDB(&models.Imgs{})}
	})
	return instanceImgsDao
}

func (dao *ImgsDao) GetImgInfo(conditions map[string]interface{}) (info models.Imgs, err error) {
	err = dao.DB.Where(conditions).First(&info).Error
	return
}

func (dao *ImgsDao) CreateImage(image models.Imgs) (uint64, error) {
	if err := dao.DB.Create(&image).Error; err != nil {
		return 0, err
	}
	return image.Id, nil
}

func (dao *ImgsDao) GetImgs(conditions map[string][]interface{}, pageParam int, pageSizeParam int) ([]models.Imgs, int, int, error) {
	imgs := []models.Imgs{}
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})

	total, err := dao.GetImgsTotal(conditions)
	if err != nil {
		return imgs, 1, 0, err
	}
	page, totalPage, pageSize, offset := dao.Page(pageParam, pageSizeParam, total)
	Db = Db.Scopes(dao.Order("create_time desc")).Offset(offset).Limit(pageSize)
	if err := Db.Find(&imgs).Error; err != nil {
		return imgs, page, totalPage, err
	}
	return imgs, page, totalPage, nil
}

func (dao *ImgsDao) GetImgsTotal(conditions map[string][]interface{}) (int64, error) {
	Db := dao.DB.Model(&models.Imgs{})
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})
	var count int64
	err := Db.Count(&count).Error
	return count, err
}

func (dao *ImgsDao) DeleteImage(conditions map[string][]interface{}) error {
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})
	err := Db.Delete(&models.Imgs{}).Error
	if err != nil {
		return err
	}
	return nil
}
