package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppArticleType struct {
	postgresqlx.BaseModle
	AppArticleTypeFields
}

type AppArticleTypeFields struct {
	Id         uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Name       string    `gorm:"name;not null" json:"name" form:"name" name:"name"`
	Status     int       `gorm:"status;not null" json:"status" form:"status" name:"status"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"文章类型创建时间" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"文章类型修改时间" name:"update_time"`
}

type AppArticleTypeReq struct {
	Name string `label:"名称" json:"name"`
}

func (articleType *AppArticleType) TableName() string {
	return "cms_app.app_article_type"
}

func (articleType *AppArticleType) FillData(db *gorm.DB) {

}

func (articleType *AppArticleType) GetConnName() string {
	return "default"
}
