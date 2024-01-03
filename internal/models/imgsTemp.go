package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
)

type ImgsTemp struct {
	postgresqlx.BaseModle
	ImgsFields
}

func (article *ImgsTemp) TableName() string {
	return "cms_app.app_imgs_temp"
}

func (article *ImgsTemp) FillData(db *gorm.DB) {

}

func (article *ImgsTemp) GetConnName() string {
	return "default"
}
