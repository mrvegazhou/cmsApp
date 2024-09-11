package imgs

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/middleware"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/AES"
	"cmsApp/pkg/DES"
	"cmsApp/pkg/jwt"
	"cmsApp/pkg/utils/arrayx"
	"cmsApp/pkg/utils/imagex"
	"cmsApp/pkg/utils/stringx"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"net/http"
	"strings"
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
	rg.POST("/upload", middleware.JwtAuth(), con.upload)
	rg.POST("/delete", con.delImage)
	rg.POST("/personalImageList", con.personalImageList)
}

func (apicon imgsController) personalShow(c *gin.Context) {
	name := c.Param("name")
	token, err := c.Cookie("__t")
	token = strings.Replace(token, " ", "+", -1)
	if err != nil {
		apicon.Error(c, errors.New(constant.IMAGE_CHECK_PERMISSION_ERR), nil)
		return
	}
	token2, err := AES.DecryptJsStr(token, configs.App.Upload.ImgCookieSecret, configs.App.Upload.ImgCookieSecret)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	_, err = jwt.Check(token2, configs.App.Login.JwtSecret, false)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId := payload.ID
	//// 检查是否为自己上传的图片
	//flag, err := apiservice.NewApiImgsService().CheckImageIsYours(cast.ToUint64(userId), name)
	//if err != nil {
	//	apicon.Error(c, err, nil)
	//	return
	//}
	//if !flag {
	//	apicon.Error(c, errors.New(constant.IMAGE_CHECK_PERMISSION_ERR), nil)
	//	return
	//}
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
	token = strings.Replace(token, " ", "+", -1)
	// 解密t t是时间戳
	des, err := DES.DesCbcDecryptByBase64(token, stringx.String2bytes(configs.App.Upload.Key), nil)
	parts := strings.SplitN(stringx.Bytes2string(des), " ", 2)
	timeTemplate := "2006-01-02"
	lastStamp, err := time.ParseInLocation(timeTemplate, parts[1], time.Local)
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

func (apicon imgsController) upload(c *gin.Context) {
	var (
		resourceId uint64
		err        error
		req        models.AppImgTempUploadReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	uid, _ := c.Get("uid")
	userId := cast.ToUint64(uid)
	imgName := ""
	fileName := ""
	if req.Type == "2" {
		resourceId, _, imgName, fileName, err = apiservice.NewApiArticleService().UploadCoverImage(req, userId)
	} else if req.Type == "1" {
		// 检查文章图片的上传限制次数50次
		_, err := apiservice.NewApiArticleService().CheckUploadLimitNum(userId)
		if err != nil {
			apicon.Error(c, err, nil)
			return
		}
		resourceId, _, imgName, fileName, err = apiservice.NewApiArticleService().UploadImage(req, userId)
	} else if req.Type == "3" || req.Type == "4" || req.Type == "5" {
		_, _, imgName, fileName, err = apiservice.NewApiImgsService().UploadImage(req, userId)
	}

	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"imageName": imgName, "fileName": fileName, "resourceId": resourceId})
}

func (apicon imgsController) delete(c *gin.Context) {
	var (
		err error
		req models.AppImgTempDeleteReq
	)
	uid, _ := c.Get("uid")
	userId := cast.ToUint64(uid)

}
