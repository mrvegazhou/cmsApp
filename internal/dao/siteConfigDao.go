package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type SiteConfigDao struct {
	DB *gorm.DB
}

var (
	instanceSiteConfig *SiteConfigDao
	onceSiteConfigDao  sync.Once
)

func NewSiteConfigDao() *SiteConfigDao {
	onceSiteConfigDao.Do(func() {
		instanceSiteConfig = &SiteConfigDao{DB: postgresqlx.GetDB(&models.SiteConfig{})}
	})
	return instanceSiteConfig
}

func (dao *SiteConfigDao) GetSiteConfigInfo(conditions map[string]interface{}) (siteConfig models.SiteConfig, err error) {
	if len(conditions) == 0 {
		err = dao.DB.Order("create_time desc").Limit(1).Take(&siteConfig).Error
	} else {
		err = dao.DB.Where(conditions).Order("create_time desc").Limit(1).Take(&siteConfig).Error
	}
	return
}
