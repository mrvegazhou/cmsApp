package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppFavorites struct {
	postgresqlx.BaseModle
	Id         uint64    `gorm:"column:id;primary_key;not null" json:"id" form:"id" label:"收藏夹标识"`
	Name       string    `gorm:"column:name;not null" json:"name" form:"name" label:"收藏夹名称"`
	UserId     string    `gorm:"column:user_id;not null" json:"userId" form:"userId" label:"用户标识"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime"`
}

type AppFavoritesItem struct {
	Id        uint64 `json:"id" form:"id" label:"收藏夹标识"`
	Name      string `json:"name" form:"name" label:"收藏夹名称"`
	IsChecked bool   `json:"isChecked" form:"isChecked" label:"是否选择"`
}

func (article *AppFavorites) TableName() string {
	return "cms_app.app_favorites"
}

func (article *AppFavorites) FillData(db *gorm.DB) {

}

func (article *AppFavorites) GetConnName() string {
	return "default"
}
