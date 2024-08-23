package models

import (
	"cmsApp/pkg/postgresqlx"
	"time"
)

type AppArticleComment struct {
	postgresqlx.BaseModle
	ArticleCommentFields
}

type ArticleCommentFields struct {
	Id           uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	UserId       uint64    `gorm:"column:user_id;not null" json:"userId" form:"userId" label:"评论人" name:"user_id"`
	ArticleId    uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" label:"文章标识" name:"article_id"`
	Content      string    `gorm:"column:content;not null" json:"content" form:"content" label:"评论内容" name:"content"`
	LikeCount    uint      `gorm:"column:like_count" json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	DislikeCount uint      `gorm:"column:dislike_count" json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量" name:"dislike_count"`
	Status       uint      `gorm:"column:status" json:"status" form:"status" label:"状态" name:"status"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"评论创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"评论修改时间" name:"update_time"`
}

type ArticleCommentPost struct {
	UserId    uint64 `json:"userId" form:"userId" label:"评论人" binding:"required" name:"user_id"`
	ArticleId uint64 `json:"articleId" form:"articleId" label:"文章标识" binding:"required" name:"article_id"`
	Content   string `json:"content" form:"content" label:"评论内容" binding:"required" name:"content"`
}

type ArticleCommentListPost struct {
	ArticleId  uint64 `json:"articleId" form:"articleId" label:"文章标识" binding:"required" name:"article_id"`
	OffSetPage int    `json:"offsetPage" form:"offsetPage" label:"分页标识" name:"offset_page"`
}
