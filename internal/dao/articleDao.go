package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type AppArticleDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceArticle *AppArticleDao
	onceArticleDao  sync.Once
)

func NewAppArticleDao() *AppArticleDao {
	onceArticleDao.Do(func() {
		instanceArticle = &AppArticleDao{DB: postgresqlx.GetDB(&models.AppArticle{})}
	})
	return instanceArticle
}

func (dao *AppArticleDao) GetAppArticle(conditions map[string]interface{}) (article models.AppArticle, err error) {
	err = dao.DB.Where(conditions).First(&article).Error
	return
}

func (dao *AppArticleDao) UpdateColumnCount(articleId uint64, column, op string, tx *gorm.DB) error {
	opStr := fmt.Sprintf("%s %s ?", column, op)
	if tx != nil {
		return tx.Model(&models.AppArticle{}).Where("id = ?", articleId).UpdateColumn(column, gorm.Expr(opStr, 1)).Error
	}
	return dao.DB.Model(&models.AppArticle{}).Where("id = ?", articleId).UpdateColumn(column, gorm.Expr(opStr, 1)).Error
}

func (dao *AppArticleDao) CreateAppArticle(articleDraft models.AppArticle) (uint64, error) {
	if err := dao.DB.Create(&articleDraft).Error; err != nil {
		return 0, err
	}
	return articleDraft.Id, nil
}

func (dao *AppArticleDao) UpdateArticle(id uint64, column models.AppArticle) (int64, error) {
	var draft models.AppArticle
	modelDB := dao.DB.Model(&draft)
	result := modelDB.Where("id = ?", id).Updates(column)
	return result.RowsAffected, modelDB.Error
}

func (dao *AppArticleDao) GetArticleList(conditions map[string][]interface{}) ([]models.AppArticle, error) {
	appArticle := []models.AppArticle{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.ArticleFields{})
	}
	Db = Db.Scopes(dao.Order("id desc"))
	if err := Db.Find(&appArticle).Error; err != nil {
		return appArticle, err
	}
	return appArticle, nil
}
