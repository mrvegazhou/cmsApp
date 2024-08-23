package models

import (
	"cmsApp/pkg/postgresqlx"
	"time"
)

type AppArticleReply struct {
	postgresqlx.BaseModle
	ArticleReplyFields
}

type ArticleReplyFields struct {
	Id           uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	ArticleId    uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" label:"文章标识" name:"article_id"`
	Pids         string    `gorm:"column:pids;not null" json:"pids" form:"pids" label:"父级别帖子" name:"pids"`
	FromUid      uint64    `gorm:"column:user_id;not null" json:"fromUid" form:"fromUid" label:"评论人" name:"from_uid"`
	ToUid        uint64    `gorm:"column:to_uid;not null" json:"userId" form:"userId" label:"回复人" name:"to_uid"`
	Content      string    `gorm:"column:content;not null" json:"content" form:"content" label:"评论内容" name:"content"`
	LikeCount    uint      `gorm:"column:like_count" json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	DislikeCount uint      `gorm:"column:dislike_count" json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量" name:"dislike_count"`
	PostId       uint64    `gorm:"column:post_id" json:"postId" form:"postId" label:"评论标识" name:"post_id"`
	Type         uint      `gorm:"column:type" json:"type" form:"type" label:"评论类型" name:"type"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"回复创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"回复修改时间" name:"update_time"`
	DeleteTime   time.Time `gorm:"column:delete_time" json:"deleteTime,omitempty" label:"删除时间" name:"delete_time"`
}

type ArticleReplyPost struct {
	PostId  uint64 `json:"postId" form:"postId" label:"评论标识" binding:"required" name:"post_id"`
	Content string `json:"content" form:"content" label:"评论内容" binding:"required" name:"content"`
	Type    uint   `json:"type" form:"type" label:"评论类型" binding:"required" name:"type"`
}

type ArticleReplyListPost struct {
	CommentId  uint64 `json:"commentId" form:"commentId" label:"评论标识" binding:"required" name:"comment_id"`
	OffSetPage int    `json:"offsetPage" form:"offsetPage" label:"分页标识" name:"offset_page"`
}
