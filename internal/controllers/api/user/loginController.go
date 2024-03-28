package user

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/errorx"
	"cmsApp/internal/middleware"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
)

type loginController struct {
	api.BaseController
}

func NewLoginController() loginController {
	return loginController{}
}

func (con loginController) Routes(rg *gin.RouterGroup) {
	rg.POST("/register/email", con.registerByEmail)
	rg.POST("/login/email", con.loginByEmail)
	rg.POST("/logout/email", con.logout)
	// 此请求不进行token的过期判断
	rg.POST("/token/refresh", con.refreshToken)
	rg.POST("/publicKey/password", con.getPasswordPublicKey)
	rg.POST("/email/code", con.sendEmailCode)
	rg.POST("/change/newPassword", middleware.JwtAuth(), con.changeNewPwdByEmailCode)
	rg.POST("/captcha/get", con.slideCaptchaGet)
	rg.POST("/captcha/check", con.slideCaptchaCheck)
}

// @Summary 用户注册
// @Id 2
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.UserRegisterReq true "User info"
// @Success 200 {object} api.SuccessResponse
// @response default {object} api.DefaultResponse
// @Router /user/register [post]
func (apicon loginController) registerByEmail(c *gin.Context) {
	var (
		err error
		req models.AppUserRegisterReq
	)

	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, errorx.NewCustomError(errorx.HTTP_BIND_PARAMS_ERR, err.Error()), nil)
		return
	}

	ip := c.ClientIP()
	err = apiservice.NewApiLoginService().RegisterByEmail(req, ip)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	apicon.Success(c, true)
}

// @Summary 用户登录
// @Id 3
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.AppUserLoginReq true "User info"
// @Success 200 {object} api.SuccessResponse
// @response default {object} models.AppUserLoginRes
// @Router /login/email [post]
func (apicon loginController) loginByEmail(c *gin.Context) {

	var (
		err          error
		req          models.AppUserLoginReq
		token        string
		refreshToken string
		times        uint
		resp         models.AppUserLoginRes
	)

	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	ip := c.ClientIP()
	flag := apiservice.NewApiLoginService().GetCaptchaIdentifying(req.Email, constant.CAPTCHA_LOGIN_MODE)
	times = apiservice.NewApiLoginService().GetLoginTimes(req.Email, ip)
	// flag!=1 表示滑动验证没有通过
	if flag != "1" && times > 3 {
		data := map[string]bool{
			"captchaVerify": true,
		}
		apicon.Error(c, errors.New(constant.CAPTCHA_NEED), data)
		return
	}

	userInfo, token, refreshToken, times, err := apiservice.NewApiLoginService().LoginByEmail(req, ip, times)
	resp.Token = token
	resp.RefreshToken = refreshToken
	resp.CaptchaVerify = false
	resp.UserInfo = userInfo

	if err != nil {
		if times > 3 {
			// 显示滑动验证
			resp.CaptchaVerify = true
		}
		apicon.Error(c, err, resp)
		return
	}

	// 将 Cookie 添加到 HTTP 响应中
	//c.SetSameSite(http.SameSiteNoneMode)
	//c.SetCookie("refresh_token", refreshToken, int(30*24*60*60), "/", "localhost", false, false)

	apicon.Success(c, resp)
}

func (apicon loginController) logout(c *gin.Context) {
	var (
		err error
		req models.AppUserLogoutReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	ip := c.ClientIP()
	id := c.GetUint64("uid")
	if id <= 0 {
		apicon.Error(c, errors.New(constant.USER_NOT_EXISTS), nil)
		return
	}
	err = apiservice.NewApiLoginService().Logout(id, ip)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

// @Summary 刷新jtoken
// @Id 4
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.UserRefreshTokenReq true "info"
// @Success 200 {json} {"status":1,"message":"success","data":{"jtoken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHAiOiIyMDIxLTEyLTI2VDE5OjI1OjI4Ljg0OTIzNzUrMDg6MDAiLCJOYW1lIjoiZ3BocGVyIiwiVWlkIjo0fQ==.ab81bb7134978afe976df55b45789aefd10f6c3edb969bae283c32c080083b89"}}
// @response default {object} api.DefaultResponse
// @Router /user/refresh [post]
func (apicon loginController) refreshToken(c *gin.Context) {
	var (
		err      error
		req      models.AppUserRefreshTokenReq
		newToken string
	)
	auth := c.Request.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		apicon.Error(c, errors.New(constant.TOKEN_NIL), nil)
		return
	}
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	// 解析token
	payload, err := jwt.Check(token, configs.App.Login.JwtSecret, true)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	uid := cast.ToUint64(payload.ID)
	newToken, newRefreshToken, err := apiservice.NewApiLoginService().RefreshToken(uid, req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	apicon.Success(c, gin.H{
		"token":        newToken,
		"refreshToken": newRefreshToken,
	})
}

func (apicon loginController) getPasswordPublicKey(c *gin.Context) {
	apicon.Success(c, configs.App.Rsa.PublicStr)
}

// @Summary 发送忘记密码的验证码到邮箱
// @Id 6
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.AppUserEmailReq true "info"
// @Success 200 {json} {"status":1,"message":"success","data":true}
// @response default {object} api.DefaultResponse
// @Router /email/code [post]
func (apicon loginController) sendEmailCode(c *gin.Context) {
	var (
		err error
		req models.AppUserEmailReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	maps := apiservice.NewApiLoginService().SendEmailCodeType()
	ty, ok := maps[req.CodeType]
	if !ok {
		apicon.Error(c, errors.New(constant.PARAM_ERR), nil)
		return
	}
	// codeType=1 表示： 忘记密码发送邮件；
	// codeType=2 表示： 注册账号
	if req.CodeType == 1 {
		// 检查邮箱是否存在
		_, err := apiservice.NewApiLoginService().GetUseInfo(map[string]interface{}{"email": req.Email})
		if err == gorm.ErrRecordNotFound {
			apicon.Error(c, errors.New(constant.LOGIN_EMAIL_NO_EXIST), nil)
			return
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			apicon.Error(c, err, nil)
			return
		}
	}

	isCaptchaVerify, err := apiservice.NewApiLoginService().SendCode(req.Email, ty)
	if isCaptchaVerify {
		data := map[string]bool{
			"captchaVerify": true,
		}
		apicon.Error(c, errors.New(constant.CAPTCHA_NEED), data)
		return
	}

	if err != nil {
		apicon.Error(c, err, err)
		return
	}
	apicon.Success(c, true)
}

// @Summary 通过邮箱验证码修改密码
// @Id 7
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.AppUserChangePwdByCodeReq true "info"
// @Success 200 {json} {"code":1,"message":"success","data":true}
// @response default {object} api.DefaultResponse
// @Router /verify/emailCode [post]
func (apicon loginController) changeNewPwdByEmailCode(c *gin.Context) {
	var (
		err error
		req models.AppUserChangePwdByCodeReq
	)

	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	// 检查邮箱是否存在
	_, err = apiservice.NewApiLoginService().GetUseInfo(map[string]interface{}{"email": req.Email})
	if err == gorm.ErrRecordNotFound {
		apicon.Error(c, errors.New(constant.LOGIN_EMAIL_NO_EXIST), nil)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		apicon.Error(c, err, nil)
		return
	}

	flag, err := apiservice.NewApiLoginService().ChangeNewPwdByEmailCode(req)
	data := map[string]bool{
		"captchaVerify": flag,
	}
	if err != nil {
		apicon.Error(c, err, data)
		return
	}
	apicon.Success(c, data)
}

// @Summary 返回验证码
// @Id 8
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.AppUserChangePwdByCodeReq true "info"
// @Success 200 {json} {"status":200,"message":"success","data":""}
// @response default {object} api.DefaultResponse
// @Router /captcha/get [post]
func (apicon loginController) slideCaptchaGet(c *gin.Context) {
	data := apiservice.NewApiLoginService().GetSlideCaptcha()
	apicon.Success(c, data)
}

// @Summary 验证验证码
// @Id 9
// @Tags 用户
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param info formData models.AppUserChangePwdByCodeReq true "info"
// @Success 200 {json} {"status":200,"message":"success","data":""}
// @response default {object} api.DefaultResponse
// @Router /captcha/check [post]
func (apicon loginController) slideCaptchaCheck(c *gin.Context) {
	var (
		err error
		req models.AppUserSlideCaptchaReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	err = apiservice.NewApiLoginService().CheckSlideCaptcha(req.CaptchaType, req.Token, req.PointJson)
	if err == nil {
		apiservice.NewApiLoginService().SetCaptchaIdentifying(req.Email, req.Mode, "1")
	} else {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, nil)
}
