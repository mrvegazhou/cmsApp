package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppTagDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceAppTagDao *AppTagDao
	onceAppTagDao     sync.Once
)

func NewAppTagDao() *AppTagDao {
	onceAppTagDao.Do(func() {
		instanceAppTagDao = &AppTagDao{DB: postgresqlx.GetDB(&models.AppTag{})}
	})
	return instanceAppTagDao
}

func (dao *AppTagDao) CreateAppTag(appTag *models.AppTag) (uint64, error) {
	if err := dao.DB.Create(appTag).Error; err != nil {
		return 0, err
	}
	return appTag.Id, nil
}

func (dao *AppTagDao) GetTagList(conditions map[string][]interface{}) ([]models.AppTag, error) {
	appTag := []models.AppTag{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.AppTagFields{})
	}
	Db = Db.Scopes(dao.Order("create_time desc"))
	if err := Db.Find(&appTag).Error; err != nil {
		return appTag, err
	}
	return appTag, nil
}
