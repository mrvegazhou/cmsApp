package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleComment struct {
	postgresqlx.BaseModle
	ArticleCommentFields
}

func (comment *AppArticleComment) TableName() string {
	return "cms_app.app_article_comment"
}

func (comment *AppArticleComment) FillData(db *gorm.DB) {

}

func (comment *AppArticleComment) GetConnName() string {
	return "default"
}

type ArticleCommentFields struct {
	Id           uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	UserId       uint64    `gorm:"column:user_id;not null" json:"userId" form:"userId" label:"评论人" name:"user_id"`
	ArticleId    uint64    `gorm:"column:article_id;not null" json:"articleId" form:"articleId" label:"文章标识" name:"article_id"`
	Content      string    `gorm:"column:content;not null" json:"content" form:"content" label:"评论内容" name:"content"`
	ReplyCount   uint      `gorm:"column:reply_count" json:"replyCount" form:"replyCount" label:"回复总量" name:"reply_count"`
	LikeCount    uint      `gorm:"column:like_count" json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	DislikeCount uint      `gorm:"column:dislike_count" json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量" name:"dislike_count"`
	ReportCount  uint      `gorm:"column:report_count" json:"reportCount" form:"reportCount" label:"举报总数" name:"report_count"`
	Score        int       `gorm:"column:score" json:"score" form:"score" label:"评分" name:"score"`
	Status       uint      `gorm:"column:status" json:"status,omitempty" form:"status" label:"状态" name:"status"`
	Ip           string    `gorm:"column:ip" json:"ip" form:"ip" label:"ip" name:"ip"`
	TopReplyIds  string    `gorm:"column:top_reply_ids" json:"topReplyIds,omitempty" form:"topReplyIds" label:"前两条回复" name:"top_reply_ids"`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"评论创建时间" name:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"评论修改时间" name:"update_time"`
}

type ArticleCommentResp struct {
	Id           uint64    `json:"id" form:"id" name:"id"`
	UserId       uint64    `json:"userId" form:"userId" label:"评论人" name:"user_id"`
	ArticleId    uint64    `json:"articleId" form:"articleId" label:"文章标识" name:"article_id"`
	Content      string    `json:"content" form:"content" label:"评论内容" name:"content"`
	ReplyCount   uint      `json:"replyCount" form:"replyCount" label:"回复总量" name:"reply_count"`
	LikeCount    uint      `json:"likeCount" form:"likeCount" label:"点赞总量" name:"like_count"`
	DislikeCount uint      `json:"dislikeCount" form:"dislikeCount" label:"不喜欢总量" name:"dislike_count"`
	Ip           string    `json:"ip" form:"ip" label:"ip" name:"ip"`
	CreateTime   time.Time `json:"createTime" label:"评论创建时间" name:"create_time"`
	UpdateTime   time.Time `json:"updateTime,omitempty" label:"评论修改时间" name:"update_time"`
}

type ArticleCommentPost struct {
	UserId    uint64 `json:"userId" form:"userId" label:"评论人" binding:"required" name:"user_id"`
	ArticleId uint64 `json:"articleId" form:"articleId" label:"文章标识" binding:"required" name:"article_id"`
	Content   string `json:"content" form:"content" label:"评论内容" binding:"required" name:"content"`
}

type ArticleCommentListPost struct {
	ArticleId   uint64 `json:"articleId" form:"articleId" label:"文章标识" binding:"required"`
	Page        int    `json:"page" form:"page" label:"分页标识"`
	CurrentTime int64  `json:"currentTime" form:"currentTime" label:"时间标识"`
	OrderBy     string `json:"orderBy" form:"orderBy" label:"排序"`
}

type ArticleCommentWithUserInfo struct {
	ArticleCommentResp
	UserInfo AppUserInfo `json:"userInfo"`
}

// 评论列表内附加回复列表
type ArticleCommentReplies struct {
	Comment ArticleCommentWithUserInfo              `json:"comment"`
	Replies []ArticleReplyWithUserAndToReplyContent `json:"replies"`
	HasNext bool                                    `json:"hasNext" label:"是否有数据"` // 这里需要通过comment的回复总数来判断
}

// 评论列表 上拉分页
type ArticleCommentListResp struct {
	CommentList []ArticleCommentReplies `json:"commentList" label:"评论列表"`
	Page        int                     `json:"page" label:"页码"`
	CurrentTime int64                   `json:"currentTime" form:"currentTime" label:"时间标识"`
	HasNext     bool                    `json:"hasNext" label:"是否有数据"` // 这里需要通过article的评论总数来判断
}
