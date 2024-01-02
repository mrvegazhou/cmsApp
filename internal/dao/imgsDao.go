package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"cmsApp/pkg/utils/arrayx"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"sync"
)

type ImgsDao struct {
	DB *gorm.DB
	BaseDao
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

func (dao *ImgsDao) GetImgs(conditions map[string][]interface{}, pageParam int, pageSizeParam int) ([]models.Imgs, int, int, error) {
	imgs := []models.Imgs{}
	Db := dao.DB
	fields := dao.GetFields(models.ImgsFields{})

	for key, cond := range conditions {
		if reflect.TypeOf(cond).Kind() == reflect.Slice && len(cond) >= 2 {
			op := cond[0]
			val := cond[1]
			if arrayx.IsContain(fields, key) {
				opStr := fmt.Sprintf("%s %s", key, op)
				Db = Db.Where(opStr, val)
			}
		}
	}
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
	Db := dao.DB
	fields := dao.GetFields(models.ImgsFields{})

	for key, cond := range conditions {
		if reflect.TypeOf(cond).Kind() == reflect.Slice && len(cond) >= 2 {
			op := cond[0]
			val := cond[1]
			if arrayx.IsContain(fields, key) {
				opStr := fmt.Sprintf("%s %s", key, op)
				Db = Db.Where(opStr, val)
			}
		}
	}
	var count *int64
	err := Db.Scopes(dao.Order("create_time desc")).Count(count).Error
	return *count, err
}
