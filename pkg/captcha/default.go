package captcha

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/pkg/redisClient"
	"context"
	"errors"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"strconv"
	"time"
)

var (
	ctx = context.Background()
)

type RedisStore struct {
}

// 实现设置 captcha 的方法
func (r RedisStore) Set(id string, value string) error {
	TTL := time.Second * time.Duration(configs.App.Captcha.IntervalTime)
	key := GenCaptchaCodeKey(id, "LOGIN")
	err := redisClient.GetRedisClient().Set(ctx, key, value, TTL).Err()
	return err
}

// 实现获取 captcha 的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := GenCaptchaCodeKey(id, "LOGIN")
	val, err := redisClient.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		err := redisClient.GetRedisClient().Del(ctx, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证 captcha 的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)
	return v == answer
}

var store base64Captcha.Store = RedisStore{} //base64Captcha.DefaultMemStore

// 获取验证码
func MakeCaptcha(codeLen int) (string, string, string, error) {
	//定义一个driver
	var driver base64Captcha.Driver
	driverDigit := &base64Captcha.DriverDigit{
		Height:   56,  //高度
		Width:    100, //宽度
		MaxSkew:  0.7,
		Length:   codeLen, //数字个数
		DotCount: 80,
	}
	driver = driverDigit
	//生成验证码
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	code := store.Get(id, false)
	return code, b64s, id, err
}

// tag:唯一标记，如phone username等
// from: 标记是哪个业务申请的验证码
func GenCaptchaCodeKey(id, from string) string {
	return "CAPTCHA:" + from + ":" + id
}

func GenCaptchaCountKey(id, from string) string {
	return "CAPTCHA_COUNT" + from + ":" + id
}

func GenintervalKey(id, from string) string {
	return "CAPTCHA_INTERVAL" + from + ":" + id
}

func CreateCaptcha(id, from string, codeLen int) (string, string, string, error) {
	// 间隔
	interValKey := GenintervalKey(id, from)
	res, err := redisClient.GetRedisClient().Get(ctx, interValKey).Result()
	if res != "" {
		errStr := fmt.Sprintf(constant.CAPTCHA_KEY_INTERVAL_ERR, strconv.Itoa(configs.App.Captcha.IntervalKeyTime))
		return "", "", "", errors.New(errStr)
	}
	err = redisClient.GetRedisClient().Set(ctx, interValKey, "1", time.Second*time.Duration(configs.App.Captcha.IntervalKeyTime)).Err()

	code, b64s, codeId, err := MakeCaptcha(codeLen)
	if err != nil {
		return "", "", "", err
	}

	// 判断间隔内使用次数
	countKey := GenCaptchaCountKey(id, from)
	count, err := redisClient.GetRedisClient().Get(ctx, countKey).Result()
	if count == "" {
		recoveryTime := time.Second * time.Duration(configs.App.Captcha.RecoveryTime)
		err = redisClient.GetRedisClient().Set(ctx, countKey, 1, recoveryTime).Err()
	}
	c, err := strconv.Atoi(count)
	if c <= configs.App.Captcha.Count {
		err = redisClient.GetRedisClient().Incr(ctx, countKey).Err()
	} else {
		return "", "", "", errors.New(constant.CAPTCHA_COUNT_ERR)
	}

	// 如果code没有过期，是不允许再生成的
	//key := GenCaptchaCodeKey(id, from)
	//TTL := time.Minute * time.Duration(configs.App.Captcha.IntervalTime)
	//bools, err := redisClient.GetRedisClient().SetNX(ctx, key, code, TTL).Result()
	//fmt.Println(bools, err, "----expire----")
	//if err != nil {
	//	return "", "", "", errors.New(constant.CAPTCHA_NO_EXPIRE)
	//}

	return code, b64s, codeId, err
}

func Verify(id, capt string) bool {
	if store.Verify(id, capt, false) {
		return true
	} else {
		return false
	}
}
