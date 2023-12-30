package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type SiteInfoDao struct {
	DB *gorm.DB
}

var (
	instanceSiteInfo *SiteInfoDao
	onceSiteInfoDao  sync.Once
)

func NewSiteInfoDao() *SiteInfoDao {
	onceSiteInfoDao.Do(func() {
		instanceSiteInfo = &SiteInfoDao{DB: postgresqlx.GetDB(&models.SiteInfo{})}
	})
	return instanceSiteInfo
}

func (dao *SiteInfoDao) GetSiteInfo(conditions map[string]interface{}) (info models.SiteInfo, err error) {
	err = dao.DB.Where(conditions).First(&info).Error
	return
}
