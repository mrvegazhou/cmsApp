package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type ReportReason struct {
	postgresqlx.BaseModle
	ReportReasonFields
}

func (reportReason *ReportReason) TableName() string {
	return "cms_app.app_report_reason"
}

func (reportReason *ReportReason) FillData(db *gorm.DB) {

}

func (reportReason *ReportReason) GetConnName() string {
	return "default"
}

type ReportReasonFields struct {
	Id         uint      `gorm:"primary_key;not null" json:"id" form:"id" name:"id"`
	Name       string    `gorm:"column:name;not null" json:"name" form:"name" label:"举报原因" name:"name"`
	Conditions string    `gorm:"column:conditions;not null" json:"conditions" form:"conditions" label:"举报需要填写的条件" name:"conditions"`
	Pid        uint      `gorm:"column:pid;not null" json:"pid" form:"pid" label:"父级标识" name:"pid"`
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" label:"创建时间" name:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"updateTime,omitempty" label:"修改时间" name:"update_time"`
}

type ReportReasonConstraint struct {
	IsRequired bool   `json:"isRequired" form:"isRequired"`
	IsShow     bool   `json:"isShow" form:"isShow"`
	Max        uint   `json:"max" form:"max"`
	Min        uint   `json:"min" form:"min"`
	Tips       string `json:"tips" form:"tips"`
}

type ReportReasonCondition struct {
	Url      ReportReasonConstraint `json:"url" form:"url"`
	Desc     ReportReasonConstraint `json:"desc" form:"desc"`
	Pictures ReportReasonConstraint `json:"pictures" form:"pictures"`
}

type ReportReasonResp struct {
	Id        uint                  `json:"id" form:"id"`
	Name      string                `json:"name" form:"name"`
	Condition ReportReasonCondition `json:"condition" form:"condition"`
	Nodes     []ReportReasonResp    `json:"nodes" form:"nodes"`
}
