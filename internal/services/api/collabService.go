package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/jwt"
	"cmsApp/pkg/redisClient"
	"cmsApp/pkg/utils/number"
	"cmsApp/pkg/utils/snowflake"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"strings"
	"sync"
	"time"
)

type apiCollabService struct {
	Dao *dao.AppArticleDao
}

var (
	instanceApiCollabService *apiCollabService
	onceApiCollabService     sync.Once
)

func NewApiCollabService() *apiCollabService {
	onceApiCollabService.Do(func() {
		instanceApiCollabService = &apiCollabService{
			Dao: dao.NewAppArticleDao(),
		}
	})
	return instanceApiCollabService
}

func (ser *apiCollabService) JoinCollab(uid uint64, articleId uint64, userIds []uint64, expireTime int64) (string, error) {
	var ctx = context.Background()
	rdb := redisClient.GetRedisClient()
	// 检查用户id是否合法
	userList, err := NewApiUserService().GetUserList(userIds)
	if err != nil {
		return "", err
	}
	var ids = make([]interface{}, 0)
	for _, item := range userList {
		ids = append(ids, fmt.Sprintf("%d", item.Id))
	}
	setKey := fmt.Sprintf("%s:%d:%d", constant.REDIS_COLLAB_USER, uid, articleId)
	// 生成token
	tokenKey := fmt.Sprintf("%s:%d:%d", constant.REDIS_COLLAB_TOKEN, uid, articleId)
	tokenStr, tokenExp, err := ser.GenURlToken(uid, articleId, userIds, expireTime)
	if err != nil {
		return "", fmt.Errorf("%s transaction(token): %w", constant.COLLAB_TOKEN_ERR, err)
	}
	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			_, err = pipe.SRem(ctx, setKey, ids).Result()
			if err == nil {
				if _, err := pipe.SAdd(ctx, setKey, ids).Result(); err != nil {
					return fmt.Errorf("%s transaction(add): %w", constant.COLLAB_INVITE_USER_ERR, err)
				} else {
					// TOKEN
					if tokenExp.IsZero() {
						pipe.Set(ctx, tokenKey, tokenStr, 0)
					} else {
						pipe.Set(ctx, tokenKey, tokenStr, time.Duration(expireTime)*time.Second)
					}
					if expireTime != -1 {
						pipe.Expire(ctx, setKey, time.Duration(expireTime)*time.Second)
					}
				}
			} else {
				return fmt.Errorf("%s transaction(srem): %w", constant.COLLAB_INVITE_USER_ERR, err)
			}
			return nil
		})
		return err
	})
	if err != nil {
		return "", fmt.Errorf("%s transaction: %w", constant.COLLAB_INVITE_USER_ERR, err)
	}
	return tokenStr, nil
}

func (ser *apiCollabService) KickOutCollab(uid uint64, articleId uint64, userIds []uint64) error {
	var ctx = context.Background()
	setKey := fmt.Sprintf("%s:%d:%d", constant.REDIS_COLLAB_USER, uid, articleId)
	var ids = make([]interface{}, 0)
	for _, id := range userIds {
		ids = append(ids, fmt.Sprintf("%d", id))
	}
	_, err := redisClient.GetRedisClient().SRem(ctx, setKey, ids).Result()
	if err != nil {
		return fmt.Errorf("%s transaction(srem): %w", constant.COLLAB_KICKOUT_USER_ERR, err)
	}
	return nil
}

func (ser *apiCollabService) ExitCollab(uid, articleId uint64) error {
	var ctx = context.Background()
	setKey := fmt.Sprintf("%s:%d:%d", constant.REDIS_COLLAB_USER, uid, articleId)
	_, err := redisClient.GetRedisClient().Del(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("%s transaction(del): %w", constant.COLLAB_EXIT_ERR, err)
	}
	return err
}

// 文章id => { {userList:用户info, articleInfo:文章info} }
func (ser *apiCollabService) ShowKeysCollab(uid uint64) (map[uint64]interface{}, error) {
	var ctx = context.Background()
	setKey := fmt.Sprintf("%s:%d:*", constant.REDIS_COLLAB_USER, uid)
	res, err := redisClient.GetRedisClient().Keys(ctx, setKey).Result()
	if err != nil {
		return map[uint64]interface{}{}, fmt.Errorf("%s transaction(keys): %w", constant.COLLAB_LIST_ERR, err)
	}
	userIds := make([]uint64, 0)
	articleIds := make([]uint64, 0)
	for _, key := range res {
		listIds, err := redisClient.GetRedisClient().SMembers(ctx, key).Result()
		if err == nil {
			for _, id := range listIds {
				userIds = append(userIds, cast.ToUint64(id))
			}
		}
		strSplit := strings.Split(key, ":")
		articleId := cast.ToUint64(strSplit[len(strSplit)-1])
		if articleId != 0 {
			articleIds = append(articleIds, articleId)
		}
	}
	// 去重
	setIds := number.RemoveRepeatedInArr(userIds)
	// in 查询用户列表
	userList, err := NewApiUserService().GetUserList(setIds)
	userMap := make(map[uint64]models.AppUser, len(userList))
	for i := 0; i < len(userList); i++ {
		userList[i].Phone = ""
		userMap[userList[i].Id] = userList[i]
	}

	setArticleIds := number.RemoveRepeatedInArr(articleIds)
	// in 查询文章列表
	articleList, err := NewApiArticleService().GetArticleList(setArticleIds)
	articleMap := make(map[uint64]models.CollabArticleInfo, len(userList))
	for i := 0; i < len(articleList); i++ {
		collab := models.CollabArticleInfo{
			TokenUrl: "",
			Info:     articleList[i],
		}
		articleMap[articleList[i].Id] = collab
	}

	// 组合返回结构
	resMap := make(map[uint64]interface{}, len(res))
	for _, key := range res {
		strSplit := strings.Split(key, ":")
		articleId := cast.ToUint64(strSplit[len(strSplit)-1])
		listIds, err := redisClient.GetRedisClient().SMembers(ctx, key).Result()
		if err == nil {
			tmpArr := []models.AppUser{}
			for _, uid := range listIds {
				value, ok := userMap[cast.ToUint64(uid)]
				if ok {
					tmpArr = append(tmpArr, value)
				}
			}
			// 获取共享地址
			tokenKey := fmt.Sprintf("%s:%d:%d", constant.REDIS_COLLAB_TOKEN, uid, articleId)
			tokenUrl, err := redisClient.GetRedisClient().Get(ctx, tokenKey).Result()
			collab := models.CollabArticleInfo{}
			collab = articleMap[articleId]
			if err == nil {
				collab.TokenUrl = tokenUrl
			}
			fmt.Println(tokenKey, collab.TokenUrl, "---s---")
			resMap[articleId] = map[string]interface{}{"userList": tmpArr, "articleInfo": collab}
		}
	}
	return resMap, nil
}

// 生成邀请token
func (ser *apiCollabService) GenURlToken(uid, articleId uint64, userIds []uint64, expireTime int64) (string, time.Time, error) {
	var payload jwt.Payload
	var exp time.Time

	// 房间名称
	payload.Name = snowflake.GenIDString()
	// 本人id
	payload.Iss = cast.ToString(uid)
	// 主题
	payload.Sub = cast.ToString(articleId)

	var temp = make([]string, len(userIds)) //为了使传参类型适用于strings.join函数
	for k, v := range userIds {
		temp[k] = fmt.Sprintf("%d", v)
	}
	// 接收人
	payload.Aud = strings.Join(temp, ",")

	if expireTime != -1 {
		exp = time.Now().Add(time.Duration(expireTime) * time.Second)
	} else {
		exp = time.Time{}
	}
	payload.Exp = exp

	token, err := jwt.Generate(configs.App.Article.JwtSecret, "HS256", payload)
	return token, exp, err
}
