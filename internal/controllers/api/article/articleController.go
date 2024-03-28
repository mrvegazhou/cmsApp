package article

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/middleware"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type articleController struct {
	api.BaseController
}

func NewArticleController() articleController {
	return articleController{}
}

func (con articleController) Routes(rg *gin.RouterGroup) {
	rg.POST("/info", middleware.JwtAuth(), con.info)
	rg.POST("/uploadImage", middleware.JwtAuth(), con.uploadImage)
	rg.POST("/save/draft", middleware.JwtAuth(), con.saveArticleDraft)
	rg.POST("/draft/history", middleware.JwtAuth(), con.getDraftHistoryList)
	rg.POST("/draft/info", middleware.JwtAuth(), con.getDraftHistoryInfo)
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

// 保存文章草稿
func (apicon articleController) saveArticleDraft(c *gin.Context) {
	var (
		err error
		req models.ArticleDraft
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	userId, ok := c.Get("uid")
	if !ok {
		apicon.Error(c, errors.New(constant.TOKEN_CHECK_ERR), nil)
		return
	}
	info, err := apiservice.NewApiArticleService().SaveArticleDraft(cast.ToUint64(userId), req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
}

func (apicon articleController) uploadImage(c *gin.Context) {
	var (
		articleId uint64
		err       error
		req       models.AppArticleUploadImage
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	imgName := ""
	fileName := ""
	if req.Type == 2 {
		articleId, _, imgName, fileName, err = apiservice.NewApiArticleService().UploadCoverImage(req, userId)
	} else if req.Type == 1 {
		// 检查文章图片的上传限制次数50次
		_, err := apiservice.NewApiArticleService().CheckUploadLimitNum(userId)
		if err != nil {
			apicon.Error(c, err, nil)
			return
		}
		articleId, _, imgName, fileName, err = apiservice.NewApiArticleService().UploadImage(req, userId)
	}

	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"imageName": imgName, "fileName": fileName, "articleId": articleId})
}

// 文章草稿保存记录列表
func (apicon articleController) getDraftHistoryList(c *gin.Context) {
	var (
		err error
		req models.AppArticleReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	list, err := apiservice.NewApiArticleService().GetDraftHistoryList(req.ArticleId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, list)
}

func (apicon articleController) getDraftHistoryInfo(c *gin.Context) {
	var (
		err error
		req models.AppDraftInfoReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	info, err := apiservice.NewApiArticleService().GetDraftHistoryInfo(req.Id)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
}

// 自动保存草稿到历史记录
func (apicon articleController) autoSaveDraftInfo(c *gin.Context) {
	//var (
	//	err error
	//	req models.ArticleDraft
	//)

}
