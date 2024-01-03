package imgs

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/AES"
	"cmsApp/pkg/utils/arrayx"
	"cmsApp/pkg/utils/imagex"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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
	rg.GET("/static/p/:name", con.personalShow)
	rg.POST("/delete", con.delImage)
	rg.POST("/personalImageList", con.personalImageList)
}

func (apicon imgsController) personalShow(c *gin.Context) {
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
	token := c.Query("t")
	if token == "" {
		apicon.Error(c, errors.New(constant.DECODE_IMG_ERR), nil)
		return
	}
	// 解密t t是时间戳
	desTime := AES.Decrypt(token, configs.App.Upload.Key)
	timeTemplate := "2006-01-02 15:04:05"
	lastStamp, err := time.ParseInLocation(timeTemplate, desTime, time.Local)
	t1 := lastStamp.Unix()
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

func (apicon imgsController) delImage(c *gin.Context) {
	var (
		err error
		req models.ImgReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	err = apiservice.NewApiImgsService().DeleteImage(req.Name)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon imgsController) personalImageList(c *gin.Context) {
	var (
		err error
		req models.AppArticleImgsReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	userId := uint64(3)
	page := req.Page
	pageSize := 8
	imgList, page, totalPage, err := apiservice.NewApiImgsService().GetImagesByUserId(userId, page, pageSize)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"imgList": imgList, "page": page, "totalPage": totalPage})
}
