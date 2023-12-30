package siteInfo

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

type siteInfoController struct {
	api.BaseController
}

func NewSiteInfoController() siteInfoController {
	return siteInfoController{}
}

func (con siteInfoController) Routes(rg *gin.RouterGroup) {
	rg.POST("/page/info", con.getSiteInfo)
}

// @Summary 获取网站内容信息
// @Id 1
// @Tags 示例
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} api.SuccessResponse{data=models.SiteInfoRes}
// @response default {object} api.DefaultResponse
// @Router /page/info [post]
func (apicon siteInfoController) getSiteInfo(c *gin.Context) {
	var (
		err error
		req models.SiteInfoReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	maps := apiservice.NewApiSiteInfoService().GetSiteInfoType()
	ty, ok := maps[req.Type]
	if !ok {
		apicon.Error(c, errors.New(constant.PARAM_ERR), nil)
		return
	}
	siteInfoRes, err := apiservice.NewApiSiteInfoService().GetSiteInfo(map[string]interface{}{"type": ty})
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	fmt.Println(siteInfoRes, "----siteInfoRes---")
	apicon.Success(c, siteInfoRes)
}
