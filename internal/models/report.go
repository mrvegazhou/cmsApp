package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppReport struct {
	postgresqlx.BaseModle
	AppReportFields
}

func (report *AppReport) TableName() string {
	return "cms_app.app_article_comment"
}

func (report *AppReport) FillData(db *gorm.DB) {

}

func (report *AppReport) GetConnName() string {
	return "default"
}

type AppReportFields struct {
	Id         uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Type       string    `gorm:"column:type;not null" json:"type" form:"type" label:"举报类型" name:"type"`
	Reason     string    `gorm:"column:reason;not null" json:"reason" form:"reason" label:"举报原因" name:"reason"`
	Content    string    `gorm:"column:content" json:"content" form:"content" label:"补充内容" name:"content"`
	Imgs       string    `gorm:"column:imgs" json:"imgs" form:"imgs" label:"图片" name:"imgs"`
	ResourceId string    `gorm:"column:resource_id" json:"resourceId" form:"resource_id" label:"补充内容" name:"resource_id"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"举报创建时间" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"举报修改时间" name:"update_time"`
	DeleteTime time.Time `gorm:"column:delete_time" json:"deleteTime,omitempty" label:"举报删除时间" name:"delete_time"`
}

type ArticleCommentReportReq struct {
	Id      uint64 `json:"id" binding:"required" form:"id"`
	Type    string `json:"type" binding:"required" form:"type"`
	Content string `json:"content" form:"content"`
	Reason  string `json:"reason" binding:"required" form:"reason"`
	Imgs    string `json:"imgs" form:"imgs"`
}
