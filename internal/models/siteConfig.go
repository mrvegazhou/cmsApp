package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type SiteConfig struct {
	postgresqlx.BaseModle
	Id         uint       `gorm:"column:id;primary_key" json:"id" form:"id"`
	Version    uint       `gorm:"column:version" json:"version" form:"version"`
	Email      string     `gorm:"column:email" json:"email" form:"email"`
	Jwt        string     `gorm:"column:jwt;not null" json:"-" form:"jwt"`
	Phone      string     `gorm:"column:phone" json:"phone" form:"phone"`
	Qq         string     `gorm:"column:qq" json:"qq" form:"qq"`
	Wechat     string     `gorm:"column:wechat" json:"wechat" form:"wechat"`
	Weibo      string     `gorm:"column:weibo" json:"weibo" form:"weibo"`
	Type       string     `gorm:"column:type" json:"type" form:"type"`
	CreateTime *time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
}

type SiteConfigReq struct {
	Uid uint `json:"uid" form:"uid" label:"用户ID"`
}

type SiteConfigRes struct {
	SiteConfig SiteConfig `json:"siteConfig"`
	UserInfo   AppUserRes `json:"userInfo"`
}

func (site *SiteConfig) TableName() string {
	return "cms_app.site_config"
}

func (site *SiteConfig) FillData(db *gorm.DB) {

}

func (site *SiteConfig) GetConnName() string {
	return "default"
}
