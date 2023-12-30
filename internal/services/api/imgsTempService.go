package api

import (
	"cmsApp/internal/dao"
	"sync"
)

type apiImgsTempService struct {
	Dao *dao.ImgsTempDao
}

var (
	instanceApiImgsTempService *apiImgsTempService
	onceApiImgsTempService     sync.Once
)

func NewApiImgsTempService() *apiImgsTempService {
	onceApiImgsTempService.Do(func() {
		instanceApiImgsTempService = &apiImgsTempService{
			Dao: dao.NewImgsTempDao(),
		}
	})
	return instanceApiImgsTempService
}

func (ser *apiImgsTempService) SaveImgs() {

}
