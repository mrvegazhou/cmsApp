package siteConfig

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type siteConfigController struct {
	api.BaseController
}

func NewSiteConfigController() siteConfigController {
	return siteConfigController{}
}

func (con siteConfigController) Routes(rg *gin.RouterGroup) {
	rg.POST("/config/info", con.siteConfigInfo)
}

// @Summary 获取网站配置信息
// @Id 1
// @Tags 示例
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param authorization header string true "token"
// @Success 200 {object} api.SuccessResponse{data=models.User}
// @response default {object} api.DefaultResponse
// @Router /site/info [post]
func (apicon siteConfigController) siteConfigInfo(c *gin.Context) {
	var (
		err error
		req models.SiteConfigReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	siteConfigInfo, err := apiservice.NewApiSiteConfigService().GetSiteConfigInfo(map[string]interface{}{"type": 1})
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	userInfo := models.AppUserInfo{}
	if req.Uid != 0 {
		userInfo, err = apiservice.NewApiUserService().GetUserInfoRes(map[string]interface{}{"id": req.Uid})
	}

	sr := models.SiteConfigRes{
		SiteConfig: siteConfigInfo,
		UserInfo:   userInfo,
	}
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}

	apicon.Success(c, sr)
}
