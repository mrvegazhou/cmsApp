package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"gorm.io/gorm"
	"sync"
)

type apiUserService struct {
	Dao *dao.AppUserDao
}

var (
	instanceApiUserService *apiUserService
	onceApiUserService     sync.Once
)

func NewApiUserService() *apiUserService {
	onceApiUserService.Do(func() {
		instanceApiUserService = &apiUserService{
			Dao: dao.NewAppUserDao(),
		}
	})
	return instanceApiUserService
}

func (ser *apiUserService) GetUserInfoRes(condition map[string]interface{}) (user models.AppUserRes, err error) {
	userInfo, err := ser.Dao.GetAppUser(condition)
	if err == gorm.ErrRecordNotFound {
		return models.AppUserRes{}, nil
	}
	userInfoRes := models.AppUserRes{
		Id:         userInfo.Id,
		Nickname:   userInfo.Nickname,
		Email:      userInfo.Email,
		CreateTime: userInfo.CreateTime,
	}
	return userInfoRes, err
}
