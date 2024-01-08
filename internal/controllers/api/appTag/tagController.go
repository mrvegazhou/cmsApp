package appTag

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type tagController struct {
	api.BaseController
}

func NewAppTagController() tagController {
	return tagController{}
}

func (apicon tagController) Routes(rg *gin.RouterGroup) {
	rg.POST("/list", apicon.getTagList)
}

func (apicon tagController) getTagList(c *gin.Context) {
	var (
		err error
		req models.AppTagReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	typeList, err := apiservice.NewApiTagService().GetTagList(req.Name)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"tagList": typeList})
}
