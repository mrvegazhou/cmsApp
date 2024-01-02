package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppArticleLikeDao struct {
	DB *gorm.DB
}

var (
	instanceArticleLike *AppArticleLikeDao
	onceArticleLikeDao  sync.Once
)

func NewAppArticleLikeDao() *AppArticleLikeDao {
	onceArticleLikeDao.Do(func() {
		instanceArticleLike = &AppArticleLikeDao{DB: postgresqlx.GetDB(&models.AppArticleLike{})}
	})
	return instanceArticleLike
}

func (dao *AppArticleLikeDao) DoArticleLike(articleId, userId uint64) error {
	tx := dao.DB.Begin()
	art := models.AppArticleLike{
		ArticleId: articleId,
		UserId:    userId,
	}
	if err := tx.Create(&art).Error; err != nil {
		tx.Rollback()
		return err
	}
	err := NewAppArticleDao().UpdateColumnCount(articleId, "like_count", "+", tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (dao *AppArticleLikeDao) DoArticleUnlike(articleId, userId uint64) error {
	tx := dao.DB.Begin()
	err := tx.Delete(&models.AppArticleLike{}, "article_id = ? and user_id = ?", articleId, userId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = NewAppArticleDao().UpdateColumnCount(articleId, "like_count", "-", tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (dao *AppArticleLikeDao) CheckArticleLike(articleId, userId uint64) (art models.AppArticleLike, err error) {
	err = dao.DB.Where("article_id = ? and user_id = ?", articleId, userId).Take(&art).Error
	return
}
