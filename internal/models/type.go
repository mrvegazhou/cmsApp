package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppType struct {
	postgresqlx.BaseModle
	AppTypeFields
}

type AppTypeFields struct {
	Id         uint64    `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Name       string    `gorm:"name;not null" json:"name" form:"name" name:"name"`
	Status     int       `gorm:"status;not null" json:"status" form:"status" name:"status"`
	Pid        uint64    `gorm:"pid;not null" json:"pid" form:"pid" name:"pid"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"类型创建时间" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"类型修改时间" name:"update_time"`
}

type AppTypeReq struct {
	Name string `label:"类别名称" json:"name"`
}

type AppTypeByPidReq struct {
	Pid uint64 `label:"父类别标识" json:"pid"`
}

type AppTypeByIdReq struct {
	Id uint64 `label:"类别标识" json:"id"`
}

type AppTypeByIdInfo struct {
	Id    uint64 `json:"id" form:"id" name:"id"`
	Name  string `json:"name" form:"name" name:"name"`
	Pid   uint64 `json:"pid" form:"pid" name:"pid"`
	Pname string `json:"pname" form:"pname" name:"pname"`
}

func (appType *AppType) TableName() string {
	return "cms_app.app_type"
}

func (appType *AppType) FillData(db *gorm.DB) {

}

func (appType *AppType) GetConnName() string {
	return "default"
}
