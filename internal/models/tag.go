package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppTag struct {
	postgresqlx.BaseModle
	AppTagFields
}

type AppTagFields struct {
	Id         uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Name       string    `gorm:"name;not null" json:"name" form:"name" name:"name"`
	Status     int       `gorm:"status;not null" json:"status" form:"status" name:"status"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"标签创建时间" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"标签修改时间" name:"update_time"`
}

type AppTagInfo struct {
	Id   uint64 `label:"标识" json:"id"`
	Name string `label:"名称" json:"name"`
}

type AppTagReq struct {
	Name string `label:"名称" json:"name"`
}

func (appTag *AppTag) TableName() string {
	return "cms_app.app_tag"
}

func (appTag *AppTag) FillData(db *gorm.DB) {

}

func (appTag *AppTag) GetConnName() string {
	return "default"
}
