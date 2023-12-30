package imgs

import (
	"cmsApp/internal/controllers/api"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/utils/arrayx"
	"cmsApp/pkg/utils/imagex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type imgsController struct {
	api.BaseController
}

func NewImgsController() imgsController {
	return imgsController{}
}

func (con imgsController) Routes(rg *gin.RouterGroup) {
	rg.GET("/static/:name", con.show)
}

func (apicon imgsController) show(c *gin.Context) {
	referer := c.Request.Header.Get("Referer")
	if referer != "" {
		whiteList := []string{"http://localhost:3000/", "https://localhost:3000/", "http://localhost:3015/", "https://localhost:3015/"}
		if !arrayx.IsContain(whiteList, referer) {
			c.String(http.StatusForbidden, "Access Denied")
			return
		}
	}
	name := c.Param("name")
	url, err := apiservice.NewApiImgsService().GetImageDirs(name)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	err = imagex.CheckImage(url)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	c.File(url)
}
