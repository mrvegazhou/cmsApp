package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type ImgsTemp struct {
	postgresqlx.BaseModle
	ImgsFields
}

type ImgsTempFields struct {
	Name string `gorm:"column:name" json:"name" form:"name" name:"name"`
	Path string `gorm:"column:path" json:"path" form:"path" name:"path"`
	// 1. 文章图片 2. 封面
	Type       uint      `gorm:"column:type" json:"type" form:"type" name:"type"`
	Tags       string    `gorm:"column:tags" json:"tags" form:"tags" name:"tags"`
	Width      int       `gorm:"column:width" json:"widht" form:"widht" name:"width"`
	Height     int       `gorm:"column:height" json:"height" form:"height" name:"height"`
	UserId     uint64    `gorm:"column:user_id" json:"userId" form:"userId" name:"user_id"`
	ResourceId uint64    `gorm:"column:resource_id" json:"resourceId" form:"resourceId" name:"resource_id"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime" name:"update_time"`
}

func (temp *ImgsTemp) TableName() string {
	return "cms_app.app_imgs_temp"
}

func (temp *ImgsTemp) FillData(db *gorm.DB) {

}

func (temp *ImgsTemp) GetConnName() string {
	return "default"
}
