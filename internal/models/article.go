package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticle struct {
	postgresqlx.BaseModle
	ArticleFields
}

type ArticleFields struct {
	Id              uint64 `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Title           string `gorm:"column:title;not null" json:"title" form:"title" label:"文章标题" name:"title"`
	Description     string `gorm:"column:description" json:"description" form:"description" label:"文章描述" name:"description"`
	AuthorId        uint64 `gorm:"column:author_id;not null" json:"authorId" form:"authorId" label:"作者标识" name:"author_id"`
	Content         string `gorm:"column:content;not null" json:"content" form:"content" label:"文章内容" name:"content"`
	CoverUrl        string `gorm:"column:cover_url;not null" json:"coverUrl" form:"coverUrl" label:"封面" name:"cover_url"`
	ViewCount       uint   `gorm:"column:view_count" json:"viewCount" form:"viewCount" label:"阅读总量" name:"view_count"`
	CommentCount    uint   `gorm:"column:comment_count;not null" json:"commentCount" form:"commentCount" label:"评论总量" name:"comment_count"`
	ReplyCount      uint   `gorm:"column:reply_count;not null" json:"replyCount" form:"replyCount" label:"评论总量" name:"reply_count"`
	CollectionCount uint   `gorm:"column:collection_count" json:"collectionCount" form:"collectionCount" label:"收藏总量" name:"collection_count"`
	LikeCount       uint   `gorm:"column:like_count" json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	ShareCount      uint   `gorm:"column:share_count" json:"shareCount" form:"shareCount" label:"分享总量" name:"share_count"`
	// 1 草稿 2 正常
	State        uint      `gorm:"column:state" json:"state" form:"state" label:"文章状态" name:"state"`
	TypeId       uint      `gorm:"column:type_id" json:"typeId" form:"typeId" label:"文章类别" name:"type_id"`
	Tags         string    `gorm:"column:tags" json:"tags" form:"tags" label:"文章TAG" name:"tags"`
	IsSetCatalog uint      `gorm:"column:is_set_catalog" json:"isSetCatalog" form:"isSetCatalog" label:"设置目录" name:"is_set_catalog"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"文章创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"文章修改时间" name:"update_time"`
	DeleteTime   time.Time `gorm:"column:delete_time;default:(-);" json:"-" name:"delete_time"`
	Marks        string    `gorm:"column:marks" json:"marks" form:"marks" label:"标注" name:"marks"`
}

type Article struct {
	Id           uint64    `json:"id" form:"id" name:"id"`
	Title        string    `form:"title" label:"文章标题" json:"title"`
	Description  string    `json:"description" form:"description" label:"文章描述"`
	AuthorId     uint64    `json:"authorId" form:"author_id" label:"作者标识"`
	Content      string    `json:"content" form:"content" label:"文章内容"`
	CoverUrl     string    `json:"coverUrl" form:"cover_url" label:"封面"`
	State        uint      `json:"state" form:"state" label:"文章状态"`
	TypeId       uint      `json:"typeId" form:"type_id" label:"文章类别"`
	Tags         []uint64  `json:"tags" form:"tags" label:"文章TAG"`
	IsSetCatalog int       `json:"isSetCatalog" form:"is_set_catalog" label:"设置目录"`
	SaveType     int       `json:"saveType" form:"save_type" label:"客户端保存类型"` // 1.手动保存；2.自动保存
	CreateTime   time.Time `json:"createTime" label:"文章创建时间"`
	Marks        string    `json:"marks" form:"marks" label:"标注"`
}

type ArticleDraft struct {
	ArticleId    uint64    `json:"articleId" form:"articleId" name:"articleId"`
	Title        string    `form:"title" label:"文章标题" json:"title"`
	Description  string    `json:"description" form:"description" label:"文章描述"`
	AuthorId     uint64    `json:"authorId" form:"authorId" label:"作者标识"`
	Content      string    `json:"content" form:"content" label:"文章内容"`
	CoverUrl     string    `json:"coverUrl" form:"coverUrl" label:"封面"`
	State        uint      `json:"state" form:"state" label:"文章状态"`
	TypeId       uint      `json:"typeId" form:"typeId" label:"文章类别"`
	Tags         []uint64  `json:"tags" form:"tags" label:"文章TAG"`
	IsSetCatalog int       `json:"isSetCatalog" form:"isSetCatalog" label:"设置目录"`
	SaveType     int       `json:"saveType" form:"saveType" label:"客户端保存类型"` // 1.手动保存；2.自动保存
	CreateTime   time.Time `json:"createTime" form:"createTime" label:"文章创建时间"`
	Marks        string    `json:"marks" form:"marks" label:"标注"`
}

type AppArticleInfoReq struct {
	ArticleId string `form:"articleId" label:"文章标识" binding:"required" json:"articleId"`
}

type AppArticleReq struct {
	ArticleId uint64 `form:"articleId" label:"文章标识" binding:"required" json:"articleId"`
}

type AppDraftInfoReq struct {
	Id uint64 `form:"id" label:"草稿标识" binding:"required" json:"id"`
}

type AppArticleImgsReq struct {
	Page int `form:"page" label:"页码" json:"page"`
}

type AppArticleCollabInviteReq struct {
	UserIds    []uint64 `json:"userIds" binding:"required" form:"userIds" label:"接收用户"`
	ArticleId  uint64   `json:"articleId" form:"articleId" label:"文章标识"`
	ExpireName string   `json:"expireName" form:"expireName" label:"过期时间"`
}

func (article *AppArticle) TableName() string {
	return "cms_app.app_article"
}

func (article *AppArticle) FillData(db *gorm.DB) {

}

func (article *AppArticle) GetConnName() string {
	return "default"
}
