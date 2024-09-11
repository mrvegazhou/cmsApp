package report

import (
	"cmsApp/internal/controllers/api"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type reportController struct {
	api.BaseController
}

func NewReportController() reportController {
	return reportController{}
}

func (con reportController) Routes(rg *gin.RouterGroup) {
	rg.POST("/reason/list", con.getReportList)

}

func (apicon reportController) getReportList(c *gin.Context) {
	list, err := apiservice.NewApiReportReasonService().GetReportReasons()
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, list)
}
