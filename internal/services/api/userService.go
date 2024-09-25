package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/utils/arrayx"
	"github.com/jinzhu/copier"
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

func (ser *apiUserService) GetUserInfoRes(condition map[string]interface{}) (user models.AppUserInfo, err error) {
	userInfo, err := ser.Dao.GetAppUser(condition)
	if err == gorm.ErrRecordNotFound {
		return models.AppUserInfo{}, nil
	}
	userInfoRes := models.AppUserInfo{
		Id:         userInfo.Id,
		Nickname:   userInfo.Nickname,
		Email:      userInfo.Email,
		About:      userInfo.About,
		AvatarUrl:  userInfo.AvatarUrl,
		CreateTime: userInfo.CreateTime,
	}
	return userInfoRes, err
}

func (ser *apiUserService) SearchUserList(name string, pageParam int, pageSizeParam int, all bool) (userList []models.AppUserInfo, page int, totalPage int, hasNextPage bool, err error) {
	userList, page, totalPage, err = ser.Dao.SearchUserList(name, pageParam, pageSizeParam, all)
	if err == gorm.ErrRecordNotFound {
		return userList, page, totalPage, false, nil
	}
	if page < totalPage {
		hasNextPage = true
	} else {
		hasNextPage = false
	}
	return userList, page, totalPage, hasNextPage, err
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

// 通过in userIds查询出用户列表后，进行去重ids，返回user map
func (ser *apiUserService) GetUserMapListByIds(userIds []uint64) (userMap map[uint64]models.AppUserInfo, err error) {
	userIds = arrayx.RemoveRepeatedElement(userIds)
	userList, err := NewApiUserService().GetUserList(userIds)
	if err != nil {
		return userMap, err
	}
	// 获取user list 的info
	userMap = make(map[uint64]models.AppUserInfo, len(userList))
	for _, userInfo := range userList {
		userModel := models.AppUserInfo{}
		copier.CopyWithOption(&userModel, userInfo, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		userMap[userInfo.Id] = userModel
	}
	return userMap, nil
}
