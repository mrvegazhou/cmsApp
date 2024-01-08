package appType

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type typeController struct {
	api.BaseController
}

func NewTypeController() typeController {
	return typeController{}
}

func (apicon typeController) Routes(rg *gin.RouterGroup) {
	rg.POST("/list", apicon.getTypeList)
	rg.POST("/pid", apicon.getTypeListByPid)
	rg.POST("/id", apicon.getTypeInfoById)
}

func (apicon typeController) getTypeList(c *gin.Context) {
	var (
		err error
		req models.AppTypeReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	typeList, err := apiservice.NewApiTypeService().GetTypeList(req.Name)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"typeList": typeList})
}

func (apicon typeController) getTypeListByPid(c *gin.Context) {
	var (
		err error
		req models.AppTypeByPidReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	typeList, err := apiservice.NewApiTypeService().GetTypeListByPid(req.Pid)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"typeList": typeList})
}

func (apicon typeController) getTypeInfoById(c *gin.Context) {
	var (
		err error
		req models.AppTypeByIdReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	typeInfo, err := apiservice.NewApiTypeService().GetTypeInfoById(req.Id)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"typeInfo": typeInfo})
}
