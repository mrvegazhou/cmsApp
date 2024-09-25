package models

import (
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"time"
)

type AppUser struct {
	postgresqlx.BaseModle
	UserFields
}

func (user *AppUser) TableName() string {
	return "cms_app.app_user"
}

func (user *AppUser) FillData(db *gorm.DB) {

}

func (user *AppUser) GetConnName() string {
	return "default"
}

type UserFields struct {
	Id           uint64    `gorm:"column:id;primary_key" json:"id" form:"id" name:"id"`
	Nickname     string    `gorm:"column:nickname;not null" json:"nickname" form:"nickanme" name:"nickname"`
	Phone        string    `gorm:"column:phone" json:"phone" form:"phone" name:"phone"`
	About        string    `gorm:"column:about" json:"about" form:"about" name:"about"`
	AvatarUrl    string    `gorm:"column:avatar_url" json:"avatarUrl" form:"avatarUrl" name:"avatar_url"`
	Email        string    `gorm:"column:email;not null" json:"email" form:"email" name:"email"`
	Password     string    `gorm:"column:password;not null" json:"-" form:"password" name:"password"`
	Salt         string    `gorm:"column:salt;not null" json:"-" form:"salt" name:"salt"`
	RefreshToken string    `gorm:"column:refresh_token" json:"-" name:""`
	CreateTime   time.Time `gorm:"autoCreateTime;column:create_time;<-:create" json:"createTime" name:"create_time"`
	ExpirTime    time.Time `gorm:"column:expir_time" json:"expirTime,omitempty" name:"expir_time"`
	UpdateTime   time.Time `gorm:"column:update_time" json:"updateTime,omitempty" name:"update_time"`
	DeleteTime   time.Time `gorm:"column:delete_time" json:"-" name:"delete_time"`
}

type AppUserInfo struct {
	Id       uint64 `json:"id" form:"id"`
	Nickname string `json:"nickname" form:"nickanme"`
	//Phone      string    `json:"phone" form:"phone" label:"手机号"`
	About      string    `json:"about" form:"about"`
	AvatarUrl  string    `json:"avatarUrl" form:"avatarUrl" label:"头像"`
	Email      string    `json:"email" form:"email" label:"邮箱"`
	CreateTime time.Time `json:"createTime,omitempty" form:"createTime" label:"注册时间"`
}

// 注册请求
type AppUserRegisterReq struct {
	Email           string `json:"email" form:"email" binding:"required,email" label:"邮箱"`
	Code            string `json:"Code" form:"Code" binding:"required" label:"验证码"`
	Password        string `json:"password" form:"password" binding:"required" label:"密码"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" binding:"required" label:"确认密码"`
}

type AppUserLoginReq struct {
	Email    string `json:"email" form:"email" binding:"required,email" label:"邮箱"`
	Password string `json:"password" form:"password" binding:"required" label:"密码"`
}

// 验证码
type AppUserCaptchaRes struct {
	B64s          string `json:"b64s" form:"b64s" label:"验证码图片"`
	CaptchaVerify bool   `json:"captchaVerify" form:"captchaVerify" label:"是否需要验证码"`
}

type AppUserReq struct {
	Username string `json:"username" form:"username" binding:"required" label:"用户名"` //用户名
	Sex      uint   `json:"sex" form:"sex" binding:"required" label:"性别"`            //性别
	Age      uint   `json:"age" form:"age" binding:"required" label:"年龄"`            //年龄
}

type AppUserLoginRes struct {
	UserInfo      AppUser `json:"userInfo"`
	Token         string  `json:"token"`        //Jtoken 验证字符串
	RefreshToken  string  `json:"refreshToken"` //retoken 刷新token
	CaptchaVerify bool    `json:"captchaVerify"`
}

type AppUserRefreshTokenReq struct {
	RefreshToken string `form:"refreshToken" binding:"required" label:"refreshToken" json:"refreshToken"`
}

type AppUserLogoutReq struct {
	Token string `form:"token" binding:"required" label:"token" json:"token"`
}

// 发送验证码到邮箱
type AppUserEmailReq struct {
	Email    string `json:"email" form:"email" binding:"email" label:"邮箱"`
	CodeType int    `json:"codeType" form:"codeType" label:"验证码类型"`
}

// 发送验证码到邮箱后进行验证的请求
type AppUserChangePwdByCodeReq struct {
	Email              string `json:"email" form:"email" binding:"email" label:"邮箱"`
	Code               string `json:"code" form:"code" binding:"required" label:"验证码"`
	NewPassword        string `json:"newPassword" form:"newPassword" binding:"required" label:"新密码"`
	ConfirmNewPassword string `json:"confirmNewPassword" form:"confirmNewPassword" binding:"required,eqfield=NewPassword" label:"确认密码"`
}

// 验证码check判断
type AppUserSlideCaptchaReq struct {
	Email       string `json:"email" form:"email" binding:"email" label:"邮箱"`
	Mode        string `json:"mode" form:"mode" binding:"required"`
	CaptchaType string `json:"captchaType" form:"captchaType" binding:"required" label:"验证码类型"`
	PointJson   string `json:"pointJson" form:"pointJson" binding:"required"`
	Token       string `json:"token" form:"token" binding:"required"`
}

type AppUserPwdPubKeyReq struct {
}

type AppUserSearchNameRes struct {
	Name string `json:"name"`
	Page int    `json:"page"`
}
