package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type AppArticle struct {
	postgresqlx.BaseModle
	Id              uint64 `gorm:"primary_key;not null" json:"id" form:"id"`
	Title           string `gorm:"column:title;not null" json:"title" form:"title" label:"文章标题"`
	Description     string `gorm:"column:description" json:"description" form:"description" label:"文章描述"`
	AuthorId        uint64 `gorm:"column:author_id;not null" json:"authorId" form:"author_id" label:"作者标识"`
	Content         string `gorm:"column:content;not null" json:"content" form:"content" label:"文章内容"`
	CoverUrl        string `gorm:"column:cover_url;not null" json:"coverUrl" form:"cover_url" label:"封面"`
	ViewCount       uint   `gorm:"column:view_count" json:"viewCount" form:"view_count" label:"阅读总量"`
	CommentCount    uint   `gorm:"column:comment_count;not null" json:"commentCount" form:"comment_count" label:"评论总量"`
	CollectionCount uint   `gorm:"column:collection_count" json:"collectionCount" form:"collection_count" label:"收藏总量"`
	LikeCount       uint   `gorm:"column:like_count" json:"likeCount" form:"like_count" label:"点赞总量"`
	ShareCount      uint   `gorm:"column:share_count" json:"shareCount" form:"share_count" label:"分享总量"`
	// 1 草稿 2 正常
	State      uint      `gorm:"column:state" json:"state" form:"state" label:"文章状态"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"文章创建时间"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"文章修改时间"`
	DeleteTime time.Time `gorm:"column:delete_time;default:(-);" json:"-"`
}

type AppArticleReq struct {
	ArticleId uint64 `form:"articleId" binding:"required" label:"文章标识" json:"articleId"`
}

type AppArticleUploadImage struct {
	File      *multipart.FileHeader `form:"file0" label:"文件" binding:"required"`
	ArticleId uint64                `form:"articleId" label:"文件标识"`
	Tags      string                `form:"tags" label:"文件标签"`
}

func (article *AppArticle) TableName() string {
	return "cms_app.app_article"
}

func (article *AppArticle) FillData(db *gorm.DB) {

}

func (article *AppArticle) GetConnName() string {
	return "default"
}
