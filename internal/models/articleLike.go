package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleLike struct {
	postgresqlx.BaseModle
	Id         uint64    `gorm:"column:id;primary_key;not null" json:"id" form:"id"`
	ArticleId  uint64    `gorm:"column:article_id;not null" json:"articleId" form:"article_id"`
	UserId     uint64    `gorm:"column:user_id;not null" json:"userId" form:"user_id"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
}

type AppArticleLikeReq struct {
	ArticleId uint64 `form:"article_id" binding:"required" label:"文章标识" json:"articleId"`
}

type AppArticleToolBarDataResp struct {
	IsLiked         bool                        `label:"是否点赞" json:"isLiked"  form:"isLiked"`
	IsCollected     bool                        `label:"是否收藏" json:"isCollected" form:"isCollected"`
	Favorites       map[uint64]AppFavoritesItem `label:"收藏夹" json:"favorites" form:"favorites"`
	IsReport        bool                        `label:"举报" json:"isReport" form:"isReport"`
	LikeCount       uint                        `label:"喜欢数" json:"likeCount" form:"likeCount"`
	CommentCount    uint                        `label:"评论数" json:"commentCount" form:"commentCount"`
	CollectionCount uint                        `label:"收藏数" json:"collectionCount" form:"collectionCount"`
	ShareCount      uint                        `label:"分享数" json:"shareCount" form:"shareCount"`
}

func (articleLike *AppArticleLike) TableName() string {
	return "cms_app.app_article_like"
}

func (articleLike *AppArticleLike) FillData(db *gorm.DB) {

}

func (articleLike *AppArticleLike) GetConnName() string {
	return "default"
}
