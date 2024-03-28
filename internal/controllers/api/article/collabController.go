package article

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/middleware"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/utils/arrayx"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type collabController struct {
	api.BaseController
}

func NewCollabController() collabController {
	return collabController{}
}

func (con collabController) Routes(rg *gin.RouterGroup) {
	rg.POST("/collab/invite", middleware.JwtAuth(), con.invite)
	rg.POST("/collab/kickout", middleware.JwtAuth(), con.kickOut)
	rg.POST("/collab/exit", middleware.JwtAuth(), con.exitCollab)
	rg.POST("/collab/view", middleware.JwtAuth(), con.listCollab)
	rg.POST("/collab/check", middleware.JwtAuth(), con.checkCollab)
}

// @Summary 通过名称搜索用户列表后加入协作
func (apicon collabController) invite(c *gin.Context) {
	var (
		err error
		req models.AppArticleCollabInviteReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	if len(req.UserIds) == 0 {
		apicon.Error(c, errors.New(constant.COLLAB_INVITE_USER_NONE), nil)
		return
	}
	ttls := make(map[string]int64, 5)
	ttls["ttl0"] = -1
	ttls["ttl1w"] = 604800
	ttls["ttl1d"] = 86400
	ttls["ttl5h"] = 18000
	ttls["ttl30s"] = 1800
	keys := make([]string, 0, 5)
	for k := range ttls {
		keys = append(keys, k)
	}
	if !arrayx.IsContain(keys, req.ExpireName) {
		apicon.Error(c, errors.New(constant.COLLAB_INVITE_TTL_ERR), nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	token, err := apiservice.NewApiCollabService().JoinCollab(userId, req.ArticleId, req.UserIds, ttls[req.ExpireName])
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, token)
}

func (apicon collabController) kickOut(c *gin.Context) {
	var (
		err error
		req models.AppArticleCollabInviteReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	err = apiservice.NewApiCollabService().KickOutCollab(userId, req.ArticleId, req.UserIds)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon collabController) exitCollab(c *gin.Context) {
	var (
		err error
		req models.CollabToken
	)

	//userId, _ := c.Get("uid")
	userId := uint64(3)
	err = apiservice.NewApiCollabService().ExitCollab(userId, req.Token)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon collabController) listCollab(c *gin.Context) {
	var err error
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	collabList, err := apiservice.NewApiCollabService().ShowKeysCollab(userId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, collabList)
}

func (apicon collabController) checkCollab(c *gin.Context) {
	//var (
	//	err error
	//	req models.CollabToken
	//)
	//
	//err = apicon.FormBind(c, &req)
	//if err != nil {
	//	apicon.Error(c, err, nil)
	//	return
	//}
	//tokenInfo := models.CollabTokenInfo{}
	////userId, _ := c.Get("uid")
	//userId := uint64(4)
	//if req.Token == "" {
	//	tokenInfo.IsCollab = false
	//} else {
	//	tokenInfo = apiservice.NewApiCollabService().CheckCollabToken(userId, req.Token)
	//}
	//apicon.Success(c, tokenInfo)

	jsonInfo := make(map[string]string) //注意该结构接受的内容
	c.BindJSON(&jsonInfo)
	tokenInfo := models.CollabTokenInfo{}
	userId, ok := c.Get("uid")
	if !ok {
		apicon.Error(c, errors.New(constant.TOKEN_CHECK_ERR), nil)
		return
	}
	if jsonInfo["token"] == "" {
		tokenInfo.IsCollab = false
	} else {
		tokenInfo = apiservice.NewApiCollabService().CheckCollabToken(cast.ToUint64(userId), jsonInfo["token"])
	}
	apicon.Success(c, tokenInfo)
}
