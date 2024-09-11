package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleReply struct {
	postgresqlx.BaseModle
	ArticleReplyFields
}

func (reply *AppArticleReply) TableName() string {
	return "cms_app.app_article_reply"
}

func (reply *AppArticleReply) FillData(db *gorm.DB) {

}

func (reply *AppArticleReply) GetConnName() string {
	return "default"
}

type ArticleReplyFields struct {
	Id           uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	ArticleId    uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" label:"文章标识" name:"article_id"`
	CommentId    uint64    `gorm:"column:comment_id;not null" json:"commentId" form:"commentId" label:"评论标识" name:"comment_id"`
	ReplyId      uint64    `gorm:"column:reply_id;not null" json:"replyId" form:"replyId" label:"回复标识" name:"reply_id"`
	Pids         string    `gorm:"column:pids;not null" json:"pids,omitempty" form:"pids" label:"父级别帖子" name:"pids"`
	FromUid      uint64    `gorm:"column:from_uid;not null" json:"fromUid" form:"fromUid" label:"评论人" name:"from_uid"`
	ToUid        uint64    `gorm:"column:to_uid;not null" json:"toUid" form:"toUid" label:"回复人" name:"to_uid"`
	Content      string    `gorm:"column:content;not null" json:"content" form:"content" label:"评论内容" name:"content"`
	ReplyCount   uint      `gorm:"column:reply_count" json:"replyCount" form:"replyCount" label:"回复总量" name:"reply_count"`
	LikeCount    uint      `gorm:"column:like_count" json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	DislikeCount uint      `gorm:"column:dislike_count" json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量" name:"dislike_count"`
	ReportCount  uint      `gorm:"column:report_count" json:"reportCount" form:"reportCount" label:"举报总数" name:"report_count"`
	Score        int       `gorm:"column:score" json:"score" form:"score" label:"评分" name:"score"`
	Ip           string    `gorm:"column:ip" json:"ip" form:"ip" label:"ip" name:"ip"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"回复创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"回复修改时间" name:"update_time"`
	Status       uint      `gorm:"column:status" json:"status,omitempty" form:"status" label:"状态" name:"status"`
}

type ArticleReplyPost struct {
	CommentId uint64 `json:"commentId" form:"replyId" label:"评论标识" binding:"required" name:"comment_id"`
	Content   string `json:"content" form:"content" label:"评论内容" binding:"required" name:"content"`
	ReplyId   uint   `json:"replyId" form:"replyId" label:"回复标识" name:"reply_id"`
}

type ArticleReplyListPost struct {
	CommentId   uint64 `json:"commentId" form:"commentId" label:"评论标识" binding:"required"`
	ReplyId     uint64 `json:"replyId" form:"replyId" label:"回复标识"`
	Page        int    `json:"page" form:"page" label:"分页标识"`
	CurrentTime int64  `json:"currentTime" form:"currentTime" label:"时间标识"`
	OrderBy     string `json:"orderBy" form:"orderBy" label:"排序"`
}

type ArticleReplyListResp struct {
	ReplyList   []ArticleReplyWithUserAndToReplyContent `json:"replyList" label:"评论列表"`
	Page        int                                     `json:"page" label:"页码"`
	HasNext     bool                                    `json:"hasNext" label:"是否有数据"`
	CurrentTime int64                                   `json:"currentTime" form:"currentTime" label:"时间标识"`
}

type ArticleReplyWithUserAndToReplyContent struct {
	ArticleReplyWithUserResp
	ToReplyContent string `json:"toReplyContent" form:"toReplyContent"`
}

type ArticleReplyWithUserResp struct {
	Id           uint64      `json:"id" form:"id"`
	ArticleId    uint64      `json:"articleId" form:"articleId" label:"文章标识"`
	CommentId    uint64      `json:"commentId" form:"commentId" label:"评论标识"`
	FromUser     AppUserInfo `json:"fromUser" form:"fromUser" label:"评论人"`
	ToUser       AppUserInfo `json:"toUser" form:"toUser" label:"回复人"`
	Content      string      `json:"content" form:"content" label:"评论内容"`
	ReplyCount   uint        `json:"replyCount" form:"replyCount" label:"回复总量"`
	LikeCount    uint        `json:"likeCount" form:"likeCount" label:"点赞总量"`
	DislikeCount uint        `json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量"`
	Ip           string      `json:"ip" form:"ip" label:"ip"`
	CreateTime   time.Time   `json:"createTime" label:"回复创建时间"`
	UpdateTime   time.Time   `json:"updateTime,omitempty" label:"回复修改时间"`
}
