package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type Imgs struct {
	postgresqlx.BaseModle
	ImgsFields
}

type ImgsFields struct {
	Id   uint64 `gorm:"column:id;primary_key" json:"id" form:"id"`
	Name string `gorm:"column:name" json:"name" form:"name"`
	Path string `gorm:"column:path" json:"path" form:"path"`
	// 1. 文章图片
	Type       uint      `gorm:"column:type" json:"type" form:"type"`
	Tags       string    `gorm:"column:tags" json:"tags" form:"tags"`
	Width      int       `gorm:"column:width" json:"widht" form:"widht"`
	Height     int       `gorm:"column:height" json:"height" form:"height"`
	UserId     uint64    `gorm:"column:user_id" json:"userId" form:"userId"`
	ResourceId uint64    `gorm:"column:resource_id" json:"resourceId" form:"resourceId"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

type ImgsListResp struct {
	Id         uint64    `label:"序列" json:"id"`
	Name       string    `label:"图片名称" json:"name"`
	Tags       string    `label:"名称" json:"tags"`
	Type       uint      `label:"类型" json:"type"`
	Width      int       `label:"宽" json:"width"`
	Height     int       `label:"高" json:"height"`
	CreateTime time.Time `label:"创建时间" json:"createTime"`
}

func (img *Imgs) TableName() string {
	return "cms_app.app_imgs"
}

func (img *Imgs) FillData(db *gorm.DB) {

}

func (img *Imgs) GetConnName() string {
	return "default"
}
