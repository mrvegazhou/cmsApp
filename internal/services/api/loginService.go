package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/errorx"
	"cmsApp/internal/models"
	"cmsApp/pkg/emailx"
	"cmsApp/pkg/jwt"
	"cmsApp/pkg/limit"
	"cmsApp/pkg/redisClient"
	captchaService "cmsApp/pkg/slideCaptcha"
	"cmsApp/pkg/utils/random"
	"cmsApp/pkg/utils/regexpx"
	"cmsApp/pkg/utils/strings"
	"context"
	"errors"
	"fmt"
	jwtLib "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type apiLoginService struct {
	Dao *dao.AppUserDao
}

var (
	instanceApiLoginService *apiLoginService
	onceApiLoginService     sync.Once
)

func NewApiLoginService() *apiLoginService {
	onceApiUserService.Do(func() {
		instanceApiLoginService = &apiLoginService{
			Dao: dao.NewAppUserDao(),
		}
	})
	return instanceApiLoginService
}

/**
* 用户注册
**/
func (ser *apiLoginService) RegisterByEmail(req models.AppUserRegisterReq, ip string) error {
	var (
		user    models.AppUser
		err     error
		seconds = configs.App.Register.LimitDuration
		quota   = configs.App.Register.LimitRate
	)

	flag := ser.GetCaptchaIdentifying(req.Email, constant.CAPTCHA_REG_MODE)
	// 需要滑动验证码
	if flag == "1" {
		return errors.New(constant.CAPTCHA_DO_ERR)
	}

	// 限制注册频率
	key := constant.REGISTER_LIMIT_RATE
	l := limit.NewPeriodLimit(seconds, quota, key)
	val, err := l.Take(ip)
	if val == 3 {
		return errors.New(constant.REGISTER_LIMIT_RATE_ERR)
	}

	// 判断验证码正确
	maps := ser.SendEmailCodeType()
	ty, _ := maps[2]
	code := redisClient.GetRedisClient().Get(context.Background(), fmt.Sprintf(constant.SEND_EMAIL_CODE, req.Email, ty)).Val()
	if code == "" {
		return errors.New(constant.SEND_EMAIL_CODE_ERR)
	}
	if code != req.Code {
		return errors.New(constant.SEND_EMAIL_CODE_NOT_EQUAL_ERR)
	}

	user, err = ser.Dao.GetAppUser(map[string]interface{}{"email": req.Email})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if user.Id != 0 {
		return errorx.NewCustomError(errorx.HTTP_UNKNOW_ERR, constant.EMAIL_ALREADY_EXISTS_ERR)
	}
	// 密码解码
	password, err := NewApiKeyService().DecryptPasswordPrivateKey(req.Password)
	confirmPassword, err := NewApiKeyService().DecryptPasswordPrivateKey(req.ConfirmPassword)

	if !regexpx.RegPassword(password) {
		return errors.New(constant.IS_PASSWORD_ERR)
	}
	if !regexpx.RegPassword(confirmPassword) {
		return errors.New(constant.IS_PASSWORD_CONFIRM_ERR)
	}
	if password != confirmPassword {
		return errors.New(constant.TWO_NEW_PASSWORD_NOT_EQUAL_ERR)
	}

	user.Email = req.Email
	salt := random.RandString(6)
	passwordSalt := strings.Encryption(password, salt)
	user.Password = passwordSalt
	user.Salt = salt
	user.ExpirTime = time.Now().AddDate(1, 0, 0)
	user.UpdateTime = time.Now()
	user.CreateTime = time.Now()

	_, err = ser.Dao.CreateAppUser(user)
	if err == nil {
		ser.SetCaptchaIdentifying(req.Email, constant.CAPTCHA_REG_MODE, "0")
	}
	return err
}

// 登录1小时内失败6次账号自动锁定，1小时之后自动解锁
func (ser *apiLoginService) CheckLoginTimes(id string, times uint, ip string) uint {
	ctx := context.Background()
	key := fmt.Sprintf(constant.LOGIN_EMAIL_TIMES, id, ip)
	if times == 0 {
		redisClient.GetRedisClient().Set(ctx, key, 1, 1*time.Hour)
	} else {
		redisClient.GetRedisClient().Incr(ctx, key).Result()
	}
	return times + 1
}

// 清除登录次数记录
func (ser *apiLoginService) ClearLoginTimes(id, ip string) {
	ctx := context.Background()
	key := fmt.Sprintf(constant.LOGIN_EMAIL_TIMES, id, ip)
	redisClient.GetRedisClient().Del(ctx, key)
}

func (ser *apiLoginService) GetLoginTimes(id, ip string) uint {
	key := fmt.Sprintf(constant.LOGIN_EMAIL_TIMES, id, ip)
	times, err := redisClient.GetRedisClient().Get(context.Background(), key).Result()
	if err != nil || times == "" {
		return 0
	}
	t, err := strconv.Atoi(times)
	return uint(t)
}

func (ser *apiLoginService) SetCaptchaIdentifying(email, mode, val string) {
	key := fmt.Sprintf(constant.CAPTCHA_IDENTIFYING, mode, email)
	redisClient.GetRedisClient().Set(context.Background(), key, val, time.Minute*30)
}

func (ser *apiLoginService) GetCaptchaIdentifying(email, mode string) string {
	key := fmt.Sprintf(constant.CAPTCHA_IDENTIFYING, mode, email)
	return redisClient.GetRedisClient().Get(context.Background(), key).Val()
}

/**
* 验证用户登录
 */
func (ser *apiLoginService) LoginByEmail(req models.AppUserLoginReq, ip string, t uint) (userInfo models.AppUser, token string, refreshToken string, times uint, err error) {
	var (
		myClaims         jwt.MyClaims
		email            = req.Email
		originalPassword = req.Password
	)
	// 解密password
	password, err := NewApiKeyService().DecryptPasswordPrivateKey(originalPassword)
	if err != nil {
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_PASSWORD_DECRYPT_ERR)
	}
	if !regexpx.RegPassword(password) {
		return userInfo, token, refreshToken, times, errors.New(constant.IS_PASSWORD_ERR)
	}

	//times = ser.GetLoginTimes(email, ip)
	configTimes := configs.App.Login.Times
	if t > configTimes {
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_TIMES_ERR)
	}

	userInfo, err = ser.Dao.GetAppUser(map[string]interface{}{"email": email})
	if err != nil {
		times = ser.CheckLoginTimes(email, t, ip)
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_EMAIL_NO_EXIST)
	}

	if userInfo.Id == 0 {
		times = ser.CheckLoginTimes(email, t, ip)
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_ACCOUNT_ERR)
	}

	//校验密码
	passwordSalt := strings.Encryption(password, userInfo.Salt)
	if passwordSalt != userInfo.Password {
		times = ser.CheckLoginTimes(email, t, ip)
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_PASSWORD_ERR)
	}

	// 根据用户表里的过期时间设置token的失效期
	durationHours := int(userInfo.ExpirTime.Sub(userInfo.UpdateTime).Hours())

	// 规划暂停账号
	if durationHours <= 100 || userInfo.ExpirTime.String() == "0001-01-01 00:00:00 +0000 UTC" {
		return userInfo, token, refreshToken, times, errors.New(constant.LOGIN_TIME_EXPIRE)
	}

	//生成token
	myClaims = jwt.MyClaims{}
	myClaims.Name = userInfo.Email
	myClaims.ID = cast.ToString(userInfo.Id)
	fmt.Println(durationHours, time.Now().Add(time.Hour*time.Duration(durationHours)), jwtLib.NewNumericDate(time.Now().Add(time.Hour*time.Duration(durationHours))), "===s====")
	myClaims.ExpiresAt = jwtLib.NewNumericDate(time.Now().Add(time.Hour * time.Duration(durationHours)))
	token, err = jwt.Generate(myClaims, configs.App.Login.JwtSecret)
	if err != nil {
		return userInfo, token, refreshToken, times, errors.New(constant.TOKEN_GEN_ERR)
	}

	//生成refresh_token
	refreshToken = strings.Encryption(passwordSalt, strconv.FormatInt(time.Now().UnixNano(), 10))

	err = ser.Dao.UpdateColumns(map[string]interface{}{
		"id": userInfo.Id,
	}, map[string]interface{}{
		"refresh_token": refreshToken,
		"expir_time":    time.Now().AddDate(1, 0, 0),
	}, nil)

	ser.SetCaptchaIdentifying(email, constant.CAPTCHA_LOGIN_MODE, "0")
	ser.ClearLoginTimes(email, ip)
	return
}

func (ser *apiLoginService) Logout(uid uint64, ip string) error {
	userInfo, err := NewApiUserService().GetUserInfoRes(map[string]interface{}{"id": uid})
	if err != nil {
		return errors.New(constant.USER_NOT_EXISTS)
	}
	ser.SetCaptchaIdentifying(userInfo.Email, constant.CAPTCHA_LOGIN_MODE, "0")
	ser.ClearLoginTimes(userInfo.Email, ip)
	return nil
}

/**
* 使用refresh token 更换jtoken
 */
func (ser *apiLoginService) RefreshToken(uid uint64, req models.AppUserRefreshTokenReq) (token string, refreshToken string, err error) {
	var (
		user     models.AppUser
		myClaims jwt.MyClaims
	)

	user, err = ser.Dao.GetAppUser(map[string]interface{}{"id": uid})
	if err != nil {
		return
	}

	if user.RefreshToken != req.RefreshToken {
		return "", "", errors.New(constant.TOKEN_CHECK_ERR)
	}

	//校验过期时间
	expirTime, err := time.ParseInLocation("2006-01-02 15:04:05", user.ExpirTime.Format("2006-01-02 15:04:05"), time.Local)
	if err != nil || expirTime.IsZero() {
		return "", "", errors.New(constant.TOKEN_EXPIRE_FORMAT_ERR)
	}
	if time.Until(expirTime).Hours() < 0 {
		return "", "", errors.New(constant.TOKEN_EXPIRE)
	}

	//生成jtoken
	myClaims.Name = user.Nickname
	myClaims.ID = cast.ToString(user.Id)
	myClaims.ExpiresAt = jwtLib.NewNumericDate(time.Now().Local().Add(5 * time.Minute))
	token, err = jwt.Generate(myClaims, configs.App.Login.JwtSecret)
	if err != nil {
		return "", "", err
	}

	// 重新生成refreshToken
	refreshToken = strings.Encryption(user.Password, strconv.FormatInt(time.Now().UnixNano(), 10))
	// 更新update_time, refresh_token
	err = ser.Dao.UpdateColumns(map[string]interface{}{
		"id": uid,
	}, map[string]interface{}{
		"update_time":   time.Now(),
		"refresh_token": refreshToken,
	}, nil)
	if err != nil {
		return "", "", err
	}
	return
}

/**
 * 发送修改密码的验证码
 */
func (ser *apiLoginService) SendCode(sendEmailer string, codeType string) (bool, error) {

	ctx := context.Background()
	rc := redisClient.GetRedisClient()
	// 保存到redis的code的key
	key := fmt.Sprintf(constant.SEND_EMAIL_CODE, sendEmailer, codeType)
	code := rc.Get(ctx, key).Val()
	if code != "" {
		return false, errors.New(fmt.Sprintf(constant.SEND_EMAIL_CODE_EXPIR_ERR, strconv.Itoa(configs.App.Email.SendExpirDuration)))
	}

	// 判断频繁操作
	limitKey := fmt.Sprintf(constant.SEND_EMAIL_COUNT_4_MODIFY_PWD, sendEmailer)
	count, err := rc.Get(ctx, limitKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			count = "0"
		} else {
			return false, err
		}
	}
	c, err := strconv.Atoi(count)
	if err != nil {
		return false, err
	}

	if c >= 3 {
		flag := ser.GetCaptchaIdentifying(sendEmailer, fmt.Sprintf(constant.CAPTCHA_SEND_EMAIL_CODE_MODE, codeType))
		var isVerify = true
		if flag == "1" {
			isVerify = false
		}
		return isVerify, nil
	}

	if configs.App.Email.SendEmailLimitCount < c {
		return false, errors.New(constant.SEND_EMAIL_CODE_LIMIT_COUNT_ERR)
	}

	code = random.RandomNumber(6)
	title := constant.CHANGE_PWD_TITLE
	content := fmt.Sprintf(constant.CHANGE_PWD_CONTENT, code)
	err = emailx.SendEmail(title, content, sendEmailer)
	if err != nil {
		ser.SetCaptchaIdentifying(sendEmailer, fmt.Sprintf(constant.CAPTCHA_SEND_EMAIL_CODE_MODE, codeType), "0")
		return false, err
	}

	// 保存验证码到redis
	dur := time.Duration(configs.App.Email.SendExpirDuration)
	err = rc.Set(ctx, key, code, time.Minute*dur).Err()
	if err != nil {
		return false, err
	}

	// redis更新已经使用的限制次数
	if c == 0 {
		rc.Set(ctx, limitKey, 1, time.Minute)
	} else {
		rc.Incr(ctx, limitKey)
	}

	return false, nil
}

/**
 * 发送邮件验证码的类型： 1，修改密码；2，注册账号
 */
func (ser *apiLoginService) SendEmailCodeType() map[int]string {
	return map[int]string{1: "changePwd", 2: "register"}
}

/**
 * 验证修改密码发送的邮箱验证码
 */
func (ser *apiLoginService) ChangeNewPwdByEmailCode(req models.AppUserChangePwdByCodeReq) (bool, error) {
	ctx := context.Background()
	rc := redisClient.GetRedisClient()
	var isVerify bool = false

	countKey := fmt.Sprintf(constant.CHANGE_PWD_COUNT, req.Email)
	count, err := rc.Get(ctx, countKey).Result()
	if err != nil {
		rc.Set(ctx, countKey, 1, time.Minute)
	} else {
		rc.Incr(ctx, countKey)
	}
	c, err := strconv.Atoi(count)
	if err == nil {
		if c >= 3 {
			flag := ser.GetCaptchaIdentifying(req.Email, constant.CAPTCHA_CHANGE_PWD_MODE)
			if flag == "1" {
				isVerify = false
			} else {
				isVerify = true
			}
		}
	}
	maps := ser.SendEmailCodeType()
	ty, _ := maps[1]
	key := fmt.Sprintf(constant.SEND_EMAIL_CODE, req.Email, ty)
	code, err := rc.Get(ctx, key).Result()
	if err != nil {
		return isVerify, errors.New(constant.SEND_EMAIL_CODE_ERR)
	}
	if code != req.Code {
		return isVerify, errors.New(constant.SEND_EMAIL_CODE_NOT_EQUAL_ERR)
	}

	newPassword, err := NewApiKeyService().DecryptPasswordPrivateKey(req.NewPassword)
	confirmNewPassword, err := NewApiKeyService().DecryptPasswordPrivateKey(req.ConfirmNewPassword)

	if !regexpx.RegPassword(newPassword) {
		return isVerify, errors.New(constant.IS_PASSWORD_ERR)
	}

	if !regexpx.RegPassword(confirmNewPassword) {
		return isVerify, errors.New(constant.IS_PASSWORD_CONFIRM_ERR)
	}

	if newPassword != confirmNewPassword {
		return isVerify, errors.New(constant.TWO_NEW_PASSWORD_NOT_EQUAL_ERR)
	}

	salt := random.RandString(6)
	newPasswordSalt := strings.Encryption(newPassword, salt)
	err = ser.Dao.UpdateColumns(map[string]interface{}{
		"email": req.Email,
	}, map[string]interface{}{
		"password":    newPasswordSalt,
		"salt":        salt,
		"update_time": time.Now(),
	}, nil)
	if err != nil {
		return isVerify, errors.New(constant.CHANGE_PWD_ERR)
	}
	ser.SetCaptchaIdentifying(req.Email, constant.CAPTCHA_CHANGE_PWD_MODE, "0")
	return isVerify, nil
}

func (ser *apiLoginService) GetUseInfo(condition map[string]interface{}) (user models.AppUser, err error) {
	return ser.Dao.GetAppUser(condition)
}

func (ser *apiLoginService) GetUserInfoRes(condition map[string]interface{}) (user models.AppUserRes, err error) {
	userInfo, err := ser.GetUseInfo(condition)
	if err == gorm.ErrRecordNotFound {
		return models.AppUserRes{}, nil
	}
	var userInfoRes models.AppUserRes = models.AppUserRes{
		Id:         userInfo.Id,
		Nickname:   userInfo.Nickname,
		Email:      userInfo.Email,
		CreateTime: userInfo.CreateTime,
	}
	return userInfoRes, err
}

func (ser *apiLoginService) GetSlideCaptcha() map[string]interface{} {
	var factory = captchaService.NewCaptchaServiceFactory(configs.App.SlideCaptcha)
	factory.RegisterCache(captchaService.RedisCacheKey, captchaService.NewConfigRedisCacheService())
	factory.RegisterService(captchaService.BlockPuzzleCaptcha, captchaService.NewBlockPuzzleCaptchaService(factory))
	data, _ := factory.GetService(captchaService.BlockPuzzleCaptcha).Get()
	return data
}

func (ser *apiLoginService) CheckSlideCaptcha(captchaType, token, pointJson string) error {
	var factory = captchaService.NewCaptchaServiceFactory(configs.App.SlideCaptcha)
	factory.RegisterCache(captchaService.RedisCacheKey, captchaService.NewConfigRedisCacheService())
	factory.RegisterService(captchaService.BlockPuzzleCaptcha, captchaService.NewBlockPuzzleCaptchaService(factory))
	err := factory.GetService(captchaType).Check(token, pointJson)
	return err
}
