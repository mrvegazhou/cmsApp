package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type ImgsTempDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceImgsTemp *ImgsTempDao
	onceImgsTempDao  sync.Once
)

func NewImgsTempDao() *ImgsTempDao {
	onceImgsTempDao.Do(func() {
		instanceImgsTemp = &ImgsTempDao{DB: postgresqlx.GetDB(&models.ImgsTemp{})}
	})
	return instanceImgsTemp
}

func (dao *ImgsTempDao) CreateImgsTemp(imgsTemp models.ImgsTemp) (uint64, error) {
	if err := dao.DB.Create(&imgsTemp).Error; err != nil {
		return 0, err
	}
	return imgsTemp.Id, nil
}

func (dao *ImgsTempDao) GetImgsTempTotal(conditions map[string][]interface{}) (int64, error) {
	Db := dao.DB.Model(&models.ImgsTemp{})
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})
	var count int64
	err := Db.Count(&count).Error
	return count, err
}

func (dao *ImgsTempDao) GetImgInfo(conditions map[string]interface{}) (info models.ImgsTemp, err error) {
	err = dao.DB.Where(conditions).First(&info).Error
	return
}

func (dao *ImgsTempDao) GetImgs(conditions map[string][]interface{}) ([]models.ImgsTemp, error) {
	imgs := []models.ImgsTemp{}
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})
	if err := Db.Find(&imgs).Error; err != nil {
		return imgs, err
	}
	return imgs, nil
}

func (dao *ImgsTempDao) DeleteImage(conditions map[string][]interface{}) error {
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ImgsFields{})
	err := Db.Delete(&models.ImgsTemp{}).Error
	if err != nil {
		return err
	}
	return nil
}
