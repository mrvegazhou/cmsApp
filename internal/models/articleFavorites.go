package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleFavorites struct {
	postgresqlx.BaseModle
	Id          uint64    `gorm:"column:id;primary_key;not null" json:"id" form:"id"`
	FavoritesId uint64    `gorm:"column:favorites_id;not null" json:"favoritesId" form:"favoritesId" label:"收藏夹标识"`
	ArticleId   uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" label:"文章标识"`
	CreateTime  time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"文章创建时间"`
}

func (article *AppArticleFavorites) TableName() string {
	return "cms_app.app_article_favorites"
}

func (article *AppArticleFavorites) FillData(db *gorm.DB) {

}

func (article *AppArticleFavorites) GetConnName() string {
	return "default"
}
