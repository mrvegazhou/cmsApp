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

type commentController struct {
	api.BaseController
}

func NewCommentController() commentController {
	return commentController{}
}

func (con commentController) Routes(rg *gin.RouterGroup) {
	rg.POST("/save/comment", middleware.JwtAuth(), con.saveComment)
}

// 保存文章评论
func (apicon commentController) saveComment(c *gin.Context) {
	var (
		err error
		req models.ArticleCommentPost
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
	info, err := apiservice.NewApiArticleCommentService().SaveArticleComment(cast.ToUint64(userId), req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"commentInfo": info})
}

func (apicon commentController) saveCommentReply(c *gin.Context) {
	var (
		err error
		req models.ArticleReplyPost
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
	info, err := apiservice.NewApiArticleCommentService().SaveArticleReply(cast.ToUint64(userId), req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"replyInfo": info})
}

func (apicon commentController) commentList(c *gin.Context) {
	var (
		err error
		req models.ArticleCommentListPost
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	commentList, err := apiservice.NewApiArticleCommentService().GetCommentList(req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, map[string]interface{}{"commentList": commentList})
}

func (apicon commentController) replyList(c *gin.Context) {
	var (
		err error
		req models.ArticleCommentListPost
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	commentList, err := apiservice.NewApiArticleCommentService().GetCommentList(req)
}
