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
	Id   uint64 `gorm:"column:id;primary_key" json:"id" form:"id" name:"id"`
	Name string `gorm:"column:name" json:"name" form:"name" name:"name"`
	Path string `gorm:"column:path" json:"path" form:"path" name:"path"`
	// 1. 文章图片
	Type       uint      `gorm:"column:type" json:"type" form:"type" name:"type"`
	Tags       string    `gorm:"column:tags" json:"tags" form:"tags" name:"tags"`
	Width      int       `gorm:"column:width" json:"widht" form:"widht" name:"width"`
	Height     int       `gorm:"column:height" json:"height" form:"height" name:"height"`
	UserId     uint64    `gorm:"column:user_id" json:"userId" form:"userId" name:"user_id"`
	ResourceId uint64    `gorm:"column:resource_id" json:"resourceId" form:"resourceId" name:"resource_id"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime" name:"update_time"`
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

type ImgReq struct {
	Name string `label:"图片名称" json:"name"`
}

func (img *Imgs) TableName() string {
	return "cms_app.app_imgs"
}

func (img *Imgs) FillData(db *gorm.DB) {

}

func (img *Imgs) GetConnName() string {
	return "default"
}
