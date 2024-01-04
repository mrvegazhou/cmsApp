package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type AppArticleTypeDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceAppArticleTypeDao *AppArticleTypeDao
	onceAppArticleTypeDao     sync.Once
)

func NewAppArticleTypeDao() *AppArticleTypeDao {
	onceAppArticleTypeDao.Do(func() {
		instanceAppArticleTypeDao = &AppArticleTypeDao{DB: postgresqlx.GetDB(&models.AppArticleType{})}
	})
	return instanceAppArticleTypeDao
}

func (dao *AppArticleTypeDao) CreateAppArticleType(articleType *models.AppArticleType) (uint64, error) {
	if err := dao.DB.Create(articleType).Error; err != nil {
		return 0, err
	}
	return articleType.Id, nil
}

func (dao *AppArticleTypeDao) GetArticleTypeList(conditions map[string][]interface{}) ([]models.AppArticleType, error) {
	articleTypes := []models.AppArticleType{}
	Db := dao.DB
	fmt.Println(conditions, "--conditions--")
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.AppArticleTypeFields{})
	Db = Db.Scopes(dao.Order("create_time desc"))
	if err := Db.Find(&articleTypes).Error; err != nil {
		return articleTypes, err
	}
	return articleTypes, nil
}
