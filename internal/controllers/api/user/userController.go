package user

import (
	"cmsApp/internal/controllers/api"
	apiservice "cmsApp/internal/services/api"
	"github.com/gin-gonic/gin"
)

type userController struct {
	api.BaseController
}

func NewUserController() userController {
	return userController{}
}

func (con userController) Routes(rg *gin.RouterGroup) {
	//rg.POST("/info", middleware.JwtAuth(), con.info)
	rg.POST("/info", con.info)
	rg.POST("/search/trendingToday", con.search)
}

// @Summary 展示用户信息
// @Id 1
// @Tags 示例
// @version 1.0
// @Accept multipart/form-data
// @Produce json
// @Param authorization header string true "token"
// @Success 200 {object} api.SuccessResponse{data=models.User}
// @response default {object} api.DefaultResponse
// @Router /user/info [post]
func (apicon userController) info(c *gin.Context) {
	id, _ := c.Get("uid")
	id = 3
	userInfo, err := apiservice.NewApiUserService().GetUserInfoRes(map[string]interface{}{"id": id})
	if err != nil {
		apicon.Error(c, err, nil)
		return
	}
	apicon.Success(c, userInfo)
}

func (apicon userController) search(c *gin.Context) {
	apicon.Success(c, true)
}
