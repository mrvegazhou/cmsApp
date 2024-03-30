package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"fmt"
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

func (ser *apiUserService) SearchUserList(name string, pageParam int, pageSizeParam int) (userList []models.AppUser, page int, totalPage int, err error) {
	fmt.Println(name, pageParam, pageSizeParam, "==s===")
	userList, page, totalPage, err = ser.Dao.SearchUserList(name, pageParam, pageSizeParam)
	if err == gorm.ErrRecordNotFound {
		return userList, page, totalPage, nil
	}
	for i := 0; i < len(userList); i++ {
		userList[i].Phone = ""
	}
	return userList, page, totalPage, err
}

func (ser *apiUserService) GetUserList(userIds []uint64) (userList []models.AppUser, err error) {
	conditions := map[string][]interface{}{}
	if len(userIds) > 0 {
		conditions = map[string][]interface{}{
			"id": {"IN ?", userIds},
		}
		return ser.Dao.GetUserList(conditions)
	} else {
		return userList, err
	}
}
