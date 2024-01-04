package article

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type articleTypeController struct {
	api.BaseController
}

func NewArticleTypeController() articleTypeController {
	return articleTypeController{}
}

func (apicon articleTypeController) Routes(rg *gin.RouterGroup) {
	rg.POST("/typeList", apicon.getArticleTypeList)
}

func (apicon articleTypeController) getArticleTypeList(c *gin.Context) {
	var (
		err error
		req models.AppArticleTypeReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	typeList, err := apiservice.NewApiArticleTypeService().GetArticleTypeList(req.Name)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"typeList": typeList})
}
