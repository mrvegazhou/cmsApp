package models

import (
	"cmsApp/pkg/postgresqlx"
)

type AppRport struct {
	postgresqlx.BaseModle
	Id         uint64 `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Reason     string `gorm:"column:reason;not null" json:"reason" form:"reason" label:"举报原因" name:"reason"`
	Content    string `gorm:"column:content" json:"content" form:"content" label:"补充内容" name:"content"`
	ResourceId string `gorm:"column:resource_id" json:"resourceId" form:"resource_id" label:"补充内容" name:"resource_id"`
}
