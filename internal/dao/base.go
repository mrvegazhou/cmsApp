package dao

import (
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"math"
	"reflect"
)

type BaseDao struct {
	DB *gorm.DB
}

func (dao *BaseDao) GetFields(obj interface{}) []string {
	var fields []string
	getType := reflect.TypeOf(obj)
	for i := 0; i < getType.NumField(); i++ {
		fieldType := getType.Field(i)
		fields[i] = fieldType.Name
	}
	return fields
}

func (dao *BaseDao) Page(pageParam interface{}, pageSizeParam interface{}, totalParam interface{}) (int, int, int, int) {
	page := cast.ToInt(pageParam)
	pageSize := cast.ToInt(pageSizeParam)
	total := cast.ToInt(totalParam)
	totalPage := cast.ToInt(math.Ceil(float64(total) / float64(pageSize)))
	if page == 0 {
		page = 1
	}
	switch {
	case page == 0:
		page = 1
	case page > totalPage:
		page = totalPage
	}

	if pageSize <= 0 {
		pageSize = 8
	}
	offset := (page - 1) * pageSize
	return page, totalPage, pageSize, offset
}

func (dao *BaseDao) Paginate(pageParam interface{}, pageSizeParam interface{}, totalParam interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		_, _, pageSize, offset := dao.Page(pageParam, pageSizeParam, totalParam)
		return db.Offset(offset).Limit(pageSize)
	}
}

func (dao *BaseDao) Order(sort string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(sort) == 0 {
			sort = "id desc"
		}
		return db.Order(sort)
	}
}
