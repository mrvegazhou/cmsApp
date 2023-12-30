package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"sync"
)

type apiSiteConfigService struct {
	Dao *dao.SiteConfigDao
}

var (
	instanceApiSiteConfigService *apiSiteConfigService
	onceApiSiteConfigService     sync.Once
)

func NewApiSiteConfigService() *apiSiteConfigService {
	onceApiSiteConfigService.Do(func() {
		instanceApiSiteConfigService = &apiSiteConfigService{
			Dao: dao.NewSiteConfigDao(),
		}
	})
	return instanceApiSiteConfigService
}

func (ser *apiSiteConfigService) GetSiteConfigInfo(condition map[string]interface{}) (user models.SiteConfig, err error) {
	return ser.Dao.GetSiteConfigInfo(condition)
}
