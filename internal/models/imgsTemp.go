package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type ImgsTemp struct {
	postgresqlx.BaseModle
	Id      uint64 `gorm:"column:id;primary_key" json:"id" form:"id"`
	Url     string `gorm:"column:url" json:"url" form:"url"`
	BaseDir string `gorm:"column:base_dir" json:"baseDir" form:"baseDir"`
	// 1. 文章
	Type       uint      `gorm:"column:type" json:"type" form:"type"`
	ResourceId uint      `gorm:"column:resource_id" json:"resourceId" form:"resourceId"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
}

func (article *ImgsTemp) TableName() string {
	return "cms_app.app_imgs_temp"
}

func (article *ImgsTemp) FillData(db *gorm.DB) {

}

func (article *ImgsTemp) GetConnName() string {
	return "default"
}
