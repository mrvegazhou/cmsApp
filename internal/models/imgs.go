package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type Imgs struct {
	postgresqlx.BaseModle
	Id   uint64 `gorm:"column:id;primary_key" json:"id" form:"id"`
	Name string `gorm:"column:name" json:"name" form:"name"`
	Path string `gorm:"column:path" json:"path" form:"path"`
	// 1. 文章图片
	Type       uint      `gorm:"column:type" json:"type" form:"type"`
	Tags       string    `gorm:"column:tags" json:"tags" form:"tags"`
	Width      int       `gorm:"column:widht" json:"widht" form:"widht"`
	Height     int       `gorm:"column:height" json:"height" form:"height"`
	UserId     uint64    `gorm:"column:user_id" json:"userId" form:"userId"`
	ResourceId uint64    `gorm:"column:resource_id" json:"resourceId" form:"resourceId"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}

func (article *Imgs) TableName() string {
	return "cms_app.app_imgs"
}

func (article *Imgs) FillData(db *gorm.DB) {

}

func (article *Imgs) GetConnName() string {
	return "default"
}
