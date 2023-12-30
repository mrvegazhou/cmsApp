package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type SiteInfo struct {
	postgresqlx.BaseModle
	Id         uint      `gorm:"column:id;primary_key" json:"id" form:"id"`
	Title      string    `gorm:"column:title" json:"title" form:"title"`
	Content    string    `gorm:"column:content" json:"content" form:"content"`
	Type       uint      `gorm:"column:type" json:"type" form:"type"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime"`
}
type SiteInfoReq struct {
	Type string `json:"type" form:"type"`
}
type SiteInfoRes struct {
	Id         uint      `json:"id,omitempty" form:"id"`
	Title      string    `json:"title" form:"title"`
	Content    string    `json:"content" form:"content"`
	Type       uint      `json:"type" form:"type"`
	CreateTime time.Time `json:"create_time,omitempty" form:"createTime"`
}

func (info *SiteInfo) TableName() string {
	return "cms_app.site_info"
}

func (info *SiteInfo) FillData(db *gorm.DB) {

}

func (info *SiteInfo) GetConnName() string {
	return "default"
}
