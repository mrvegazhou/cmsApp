package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type AppArticleReplyDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceArticleReply *AppArticleReplyDao
	onceArticleReplyDao  sync.Once
)

func NewAppArticleReplyDao() *AppArticleReplyDao {
	onceArticleReplyDao.Do(func() {
		instanceArticleReply = &AppArticleReplyDao{DB: postgresqlx.GetDB(&models.AppArticleReply{})}
	})
	return instanceArticleReply
}

func (dao *AppArticleReplyDao) CreateArticleReply(articleReply models.AppArticleReply) (uint64, error) {
	if err := dao.DB.Create(&articleReply).Error; err != nil {
		return 0, err
	}
	return articleReply.Id, nil
}

func (dao *AppArticleReplyDao) GetArticleReply(conditions map[string]interface{}) (reply models.AppArticleReply, err error) {
	err = dao.DB.Where(conditions).First(&reply).Error
	return
}

func (dao *AppArticleReplyDao) GetArticleReplies(conditions map[string][]interface{}, orderBy string) (replies []models.AppArticleReply, err error) {
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleReplyFields{})
	if orderBy == "" {
		orderBy = "create_time desc"
	}
	Db = Db.Scopes(dao.Order(orderBy))
	if err := Db.Find(&replies).Error; err != nil {
		return replies, err
	}
	return
}

func (dao *AppArticleReplyDao) GetArticleReplyList(conditions map[string][]interface{}, pageParam int, pageSizeParam int) ([]models.AppArticleReply, int, int, error) {
	replies := []models.AppArticleReply{}
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleReplyFields{})

	total, err := dao.GetArticleReplyTotal(conditions)
	if err != nil {
		return replies, 1, 0, err
	}
	page, totalPage, pageSize, offset := dao.Page(pageParam, pageSizeParam, total)
	Db = Db.Scopes(dao.Order("create_time desc")).Offset(offset).Limit(pageSize)
	if err := Db.Find(&replies).Error; err != nil {
		return replies, page, totalPage, err
	}
	return replies, page, totalPage, nil
}

func (dao *AppArticleReplyDao) GetArticleReplyTotal(conditions map[string][]interface{}) (int64, error) {
	Db := dao.DB.Model(&models.AppArticleReply{})
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleReplyFields{})
	var count int64
	err := Db.Count(&count).Error
	return count, err
}

func (dao *AppArticleReplyDao) GetArticleReplyListNoTotal(conditions map[string][]interface{}, pageParam int, pageSizeParam int, orderBy string) ([]models.AppArticleReply, error) {
	replies := []models.AppArticleReply{}
	Db := dao.DB
	Db = dao.BaseDao.ConditionWhere(Db, conditions, models.ArticleReplyFields{})
	if pageParam == 0 {
		pageParam = 1
	}
	offset := (pageParam - 1) * pageSizeParam

	if orderBy == "" {
		orderBy = "create_time"
	}
	orderBy = fmt.Sprintf("%s desc", orderBy)

	Db = Db.Scopes(dao.Order(orderBy)).Offset(offset).Limit(pageSizeParam)
	if err := Db.Find(&replies).Error; err != nil {
		return replies, err
	}
	return replies, nil
}
