package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppArticleCommentDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceArticleComment *AppArticleCommentDao
	onceArticleCommentDao  sync.Once
)

func NewAppArticleCommentDao() *AppArticleCommentDao {
	onceArticleCommentDao.Do(func() {
		instanceArticleComment = &AppArticleCommentDao{DB: postgresqlx.GetDB(&models.AppArticleComment{})}
	})
	return instanceArticleComment
}

func (dao *AppArticleCommentDao) CreateArticleComment(articleComment models.AppArticleComment) (uint64, error) {
	if err := dao.DB.Create(&articleComment).Error; err != nil {
		return 0, err
	}
	return articleComment.Id, nil
}

func (dao *AppArticleCommentDao) GetArticleComment(conditions map[string]interface{}) (article models.AppArticleComment, err error) {
	err = dao.DB.Where(conditions).First(&article).Error
	return
}

func (dao *AppArticleCommentDao) GetArticleCommentList(conditions map[string][]interface{}, pageParam int, pageSizeParam int) ([]models.AppArticleComment, int, int, error) {
	comments := []models.AppArticleComment{}
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleCommentFields{})

	total, err := dao.GetArticleCommentTotal(conditions)
	if err != nil {
		return comments, 1, 0, err
	}
	page, totalPage, pageSize, offset := dao.Page(pageParam, pageSizeParam, total)
	Db = Db.Scopes(dao.Order("create_time desc")).Offset(offset).Limit(pageSize)
	if err := Db.Find(&comments).Error; err != nil {
		return comments, page, totalPage, err
	}
	return comments, page, totalPage, nil
}

func (dao *AppArticleCommentDao) GetArticleCommentTotal(conditions map[string][]interface{}) (int64, error) {
	Db := dao.DB.Model(&models.AppArticleComment{})
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleCommentFields{})
	var count int64
	err := Db.Count(&count).Error
	return count, err
}
