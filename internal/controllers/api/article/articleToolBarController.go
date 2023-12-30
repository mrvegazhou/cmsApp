package article

import (
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type articleToolBarController struct {
	api.BaseController
}

func NewArticleToolBarController() articleToolBarController {
	return articleToolBarController{}
}

func (con articleToolBarController) Routes(rg *gin.RouterGroup) {
	//rg.POST("/info", middleware.JwtAuth(), con.info)
	rg.POST("/like", con.doArticleLike)
	rg.POST("/unlike", con.doArticleUnlike)
	rg.POST("/toolBarData", con.getArticleToolBar)
}

func (apicon articleToolBarController) doArticleLike(c *gin.Context) {
	var (
		err error
		req models.AppArticleLikeReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := 3
	uid := uint64(userId)
	err = apiservice.NewApiArticleToolBarService().DoArticleLike(req.ArticleId, uid)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon articleToolBarController) doArticleUnlike(c *gin.Context) {
	var (
		err error
		req models.AppArticleLikeReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := 3
	uid := uint64(userId)
	err = apiservice.NewApiArticleToolBarService().DoArticleUnlike(req.ArticleId, uid)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon articleToolBarController) getArticleToolBar(c *gin.Context) {
	var (
		err  error
		req  models.AppArticleLikeReq
		resp models.AppArticleToolBarDataResp
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	userId := 3
	uid := uint64(userId)

	// 是否点赞
	resp = apiservice.NewApiArticleToolBarService().GetArticleToolBarData(req.ArticleId, uid)
	apicon.Success(c, resp)
}
