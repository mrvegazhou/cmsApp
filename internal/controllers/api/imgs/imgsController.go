package imgs

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/AES"
	"cmsApp/pkg/utils/arrayx"
	"cmsApp/pkg/utils/imagex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
	id := c.Query("id")
	if id == "" {
		apicon.Error(c, errors.New(constant.DECODE_IMG_ERR), nil)
		return
	}
	// 解密id id是时间戳
	desTime := AES.AesDecrypt(configs.App.Upload.Key, id)
	t1, err := strconv.ParseInt(desTime, 10, 64)
	if err != nil {
		apicon.Error(c, errors.New(constant.DECODE_IMG_ERR), nil)
		return
	}
	t2 := time.Now().Unix()
	// 大于一天不显示图片
	if t2-t1 > 86400 {
		apicon.Error(c, errors.New(constant.DECODE_IMG_ERR), nil)
		return
	}
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
