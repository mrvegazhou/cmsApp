package article

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/controllers/api"
	"cmsApp/internal/models"
	apiservice "cmsApp/internal/services/api"
	"cmsApp/pkg/utils/arrayx"
	"errors"
	"github.com/gin-gonic/gin"
)

type collabController struct {
	api.BaseController
}

func NewCollabController() collabController {
	return collabController{}
}

func (con collabController) Routes(rg *gin.RouterGroup) {
	//rg.POST("/info", middleware.JwtAuth(), con.info)
	rg.POST("/collab/invite", con.invite)
	rg.POST("/collab/kickout", con.kickOut)
	rg.POST("/collab/exit", con.exitCollab)
	rg.POST("/collab/view", con.listCollab)
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
		req models.AppArticleReq
	)
	err = apicon.FormBind(c, &req)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	err = apiservice.NewApiCollabService().ExitCollab(userId, req.ArticleId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, true)
}

func (apicon collabController) listCollab(c *gin.Context) {
	var (
		err error
	)
	//userId, _ := c.Get("uid")
	userId := uint64(3)
	collabList, err := apiservice.NewApiCollabService().ShowKeysCollab(userId)
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, collabList)
}
