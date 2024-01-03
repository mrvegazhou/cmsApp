package article

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type articleController struct {
	api.BaseController
}

func NewArticleController() articleController {
	return articleController{}
}

func (con articleController) Routes(rg *gin.RouterGroup) {
	//rg.POST("/info", middleware.JwtAuth(), con.info)
	rg.POST("/info", con.info)
	rg.POST("/uploadImage", con.uploadImage)
	rg.POST("/saveArticleDraft", con.saveArticleDraft)
}

func (apicon articleController) info(c *gin.Context) {
	var (
		err error
		req models.AppArticleReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	info, err := apiservice.NewApiArticleService().GetArticleInfo(req.ArticleId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
}

func (apicon articleController) saveArticleDraft(c *gin.Context) {
	var (
		err error
		req models.AppArticle
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	info, err := apiservice.NewApiArticleService().SaveArticleDraft(req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
}

func (apicon articleController) uploadImage(c *gin.Context) {
	var (
		err error
		req models.AppArticleUploadImage
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	_, imgName, fileName, err := apiservice.NewApiArticleService().UploadImage(req, userId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"imageName": imgName, "fileName": fileName})
}
