package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleHistory struct {
	postgresqlx.BaseModle
	ArticleHistoryFields
}

type ArticleHistoryFields struct {
	Id           uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	ArticleId    uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" name:"article_id"`
	Title        string    `gorm:"column:title;not null" json:"title" form:"title" label:"文章标题" name:"title"`
	Description  string    `gorm:"column:description" json:"description" form:"description" label:"文章描述" name:"description"`
	AuthorId     uint64    `gorm:"column:author_id;not null" json:"authorId" form:"author_id" label:"作者标识" name:"author_id"`
	Content      string    `gorm:"column:content;not null" json:"content" form:"content" label:"文章内容" name:"content"`
	CoverUrl     string    `gorm:"column:cover_url;not null" json:"coverUrl" form:"cover_url" label:"封面" name:"cover_url"`
	TypeId       uint      `gorm:"column:type_id" json:"typeId" form:"type_id" label:"文章类别" name:"type_id"`
	SaveType     string    `gorm:"column:save_type" json:"saveType" form:"SaveType" label:"保存方式" name:"save_type"`
	SourceType   string    `gorm:"column:source_type" json:"sourceType" form:"sourceType" label:"保存来源" name:"save_type"`
	Tags         string    `gorm:"column:tags" json:"tags" form:"tags" label:"文章TAG" name:"tags"`
	IsSetCatalog uint      `gorm:"column:is_set_catalog" json:"isSetCatalog" form:"isSetCatalog" label:"设置目录" name:"isSetCatalog"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"文章创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"文章修改时间" name:"update_time"`
}

type ArticleHistoryResp struct {
	Id           uint64       `json:"id" form:"id"`
	ArticleId    uint64       `json:"articleId" form:"articleId"`
	Title        string       `json:"title" form:"title" label:"文章标题"`
	Description  string       `json:"description" form:"description" label:"文章描述"`
	AuthorId     uint64       `json:"authorId" form:"author_id" label:"作者标识"`
	Content      string       `json:"content" form:"content" label:"文章内容"`
	CoverUrl     string       `json:"coverUrl" form:"cover_url" label:"封面"`
	TypeId       uint         `json:"typeId" form:"type_id" label:"文章类别"`
	SaveType     string       `json:"saveType" form:"SaveType" label:"保存方式"`
	SourceType   string       `json:"sourceType" form:"sourceType" label:"保存来源"`
	Tags         []AppTagInfo `json:"tags" form:"tags" label:"文章TAGS"`
	IsSetCatalog uint         `json:"isSetCatalog" form:"isSetCatalog" label:"设置目录"`
	CreateTime   time.Time    `json:"createTime" form:"createTime" label:"文章创建时间"`
}

func (article *AppArticleHistory) TableName() string {
	return "cms_app.app_article_history"
}

func (article *AppArticleHistory) FillData(db *gorm.DB) {

}

func (article *AppArticleHistory) GetConnName() string {
	return "default"
}
