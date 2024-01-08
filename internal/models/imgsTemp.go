package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
)

type ImgsTemp struct {
	postgresqlx.BaseModle
	ImgsFields
}

func (temp *ImgsTemp) TableName() string {
	return "cms_app.app_imgs_temp"
}

func (temp *ImgsTemp) FillData(db *gorm.DB) {

}

func (temp *ImgsTemp) GetConnName() string {
	return "default"
}
