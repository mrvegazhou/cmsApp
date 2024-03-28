package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppArticleHistoryDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceArticleHistory *AppArticleHistoryDao
	onceArticleHistoryDao  sync.Once
)

func NewAppArticleHistoryDao() *AppArticleHistoryDao {
	onceArticleHistoryDao.Do(func() {
		instanceArticleHistory = &AppArticleHistoryDao{DB: postgresqlx.GetDB(&models.AppArticleHistory{})}
	})
	return instanceArticleHistory
}

func (dao *AppArticleHistoryDao) CreateArticleHistory(article models.AppArticleHistory) (uint64, error) {
	if err := dao.DB.Create(&article).Error; err != nil {
		return 0, err
	}
	return article.Id, nil
}

func (dao *AppArticleHistoryDao) GetArticleHistoryInfo(conditions map[string]interface{}) (articleHistory models.AppArticleHistory, err error) {
	err = dao.DB.Where(conditions).First(&articleHistory).Error
	return
}

// 获取最后一条数据
func (dao *AppArticleHistoryDao) GetLastArticleHistoryInfo() (articleHistory models.AppArticleHistory, err error) {
	err = dao.DB.Last(&articleHistory).Error
	return
}

func (dao *AppArticleHistoryDao) UpdateColumns(conditions, field map[string]interface{}, tx *gorm.DB) error {
	if tx != nil {
		return tx.Model(&models.AppArticleHistory{}).Where(conditions).UpdateColumns(field).Error
	}
	return dao.DB.Model(&models.AppArticleHistory{}).Where(conditions).UpdateColumns(field).Error
}

func (dao *AppArticleHistoryDao) GetArticleHistoryList(conditions map[string][]interface{}) ([]models.AppArticleHistory, error) {
	articleHistory := []models.AppArticleHistory{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.ArticleHistoryFields{}).Scopes(dao.Order("id desc"))
		if err := Db.Find(&articleHistory).Error; err != nil {
			return articleHistory, err
		}
	}
	return articleHistory, nil
}

func (dao *AppArticleHistoryDao) DelArticleHistory(conditions map[string][]interface{}) (bool, error) {
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.ArticleFields{})
		err := Db.Delete(&models.AppArticleHistory{}).Error
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
