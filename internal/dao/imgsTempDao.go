package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"fmt"
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
	fmt.Println(imgsTemp, "--imgsTemp-")
	if err := dao.DB.Create(&imgsTemp).Error; err != nil {
		fmt.Println(err, "---er--")
		return 0, err
	}
	return imgsTemp.Id, nil
}
