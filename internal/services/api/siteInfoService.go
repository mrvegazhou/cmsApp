package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"sync"
)

type apiSiteInfoService struct {
	Dao *dao.SiteInfoDao
}

var (
	instanceApiSiteInfoService *apiSiteInfoService
	onceApiSiteInfoService     sync.Once
)

func NewApiSiteInfoService() *apiSiteInfoService {
	onceApiSiteInfoService.Do(func() {
		instanceApiSiteInfoService = &apiSiteInfoService{
			Dao: dao.NewSiteInfoDao(),
		}
	})
	return instanceApiSiteInfoService
}

func (ser *apiSiteInfoService) GetSiteInfo(condition map[string]interface{}) (user models.SiteInfo, err error) {
	return ser.Dao.GetSiteInfo(condition)
}

func (ser *apiSiteInfoService) GetSiteInfoType() map[string]uint {
	return map[string]uint{"LOGIN_PAGE_INFO": 1, "REGISTER_PAGE_INFO": 2}
}
