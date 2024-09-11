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
	"time"
)

type commentController struct {
	api.BaseController
}

func NewCommentController() commentController {
	return commentController{}
}

func (con commentController) Routes(rg *gin.RouterGroup) {
	rg.POST("/save/comment", middleware.JwtAuth(), middleware.RateLimit(), con.saveComment)
	rg.POST("/save/reply", middleware.JwtAuth(), con.saveReply)
	rg.POST("/comment/list", con.commentList)
	rg.POST("/reply/list", con.replyList)
	rg.POST("/comment/report", middleware.JwtAuth(), con.reportComment)
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
	ipStr := c.ClientIP()
	info, err := apiservice.NewApiArticleCommentService().SaveArticleComment(cast.ToUint64(userId), req, ipStr)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
}

func (apicon commentController) saveReply(c *gin.Context) {
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
	ipStr := c.ClientIP()
	info, err := apiservice.NewApiArticleCommentService().SaveArticleReply(cast.ToUint64(userId), req, ipStr)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, info)
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
	if &req.CurrentTime == nil {
		req.CurrentTime = time.Now().Unix()
	} else {
		if len(cast.ToString(req.CurrentTime)) == 10 {

		} else if len(cast.ToString(req.CurrentTime)) == 13 {
			req.CurrentTime = req.CurrentTime / 1000
		} else {
			apicon.Error(c, errors.New(constant.ARTICLE_COMMENT_CURRENT_TIME_ERR), nil)
			return
		}
	}
	if &req.OrderBy == nil {
		req.OrderBy = "score"
	}
	articleCommentListResp, err := apiservice.NewApiArticleCommentService().GetCommentList(req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	if articleCommentListResp.CommentList == nil {
		articleCommentListResp.CommentList = []models.ArticleCommentReplies{}
	}
	articleCommentListResp.Page = req.Page
	articleCommentListResp.CurrentTime = req.CurrentTime
	apicon.Success(c, articleCommentListResp)
}

func (apicon commentController) replyList(c *gin.Context) {
	var (
		err error
		req models.ArticleReplyListPost
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	replyResp, err := apiservice.NewApiArticleCommentService().GetReplyList(req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, replyResp)
}

func (apicon commentController) reportComment(c *gin.Context) {
	var (
		err error
		req models.ArticleCommentReportReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apiservice.NewApiArticleCommentService().HandleReport(req)
}
