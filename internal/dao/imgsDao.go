package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type ImgsDao struct {
	DB *gorm.DB
}

var (
	instanceImgs *ImgsDao
	onceImgsDao  sync.Once
)

func NewImgsDao() *ImgsDao {
	onceImgsDao.Do(func() {
		instanceImgs = &ImgsDao{DB: postgresqlx.GetDB(&models.Imgs{})}
	})
	return instanceImgs
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
