package router

import (
	"cmsApp/internal/controllers/api/article"
	"cmsApp/internal/controllers/api/imgs"
	"cmsApp/internal/controllers/api/siteConfig"
	"cmsApp/internal/controllers/api/siteInfo"
	"cmsApp/internal/controllers/api/user"
	"github.com/gin-gonic/gin"
)

// API控制器接口
type IApiController interface {
	Success(*gin.Context, interface{})
	Error(*gin.Context, error, interface{})
	FormBind(*gin.Context, interface{}) error
	Routes(*gin.RouterGroup)
}

type ApiRouter struct {
	root *gin.RouterGroup
}

func NewApiRouter() *ApiRouter {
	return &ApiRouter{}
}

func (ar ApiRouter) addRouter(con IApiController, router *gin.RouterGroup) {
	con.Routes(router)
}

func (ar ApiRouter) AddRouters() {
	{
		apiUserRouter := ar.root.Group("/user")
		{
			ar.addRouter(user.NewUserController(), apiUserRouter)
			ar.addRouter(user.NewLoginController(), apiUserRouter)
		}
	}
	{
		apiSiteConfigRouter := ar.root.Group("/site")
		{
			ar.addRouter(siteConfig.NewSiteConfigController(), apiSiteConfigRouter)
			ar.addRouter(siteInfo.NewSiteInfoController(), apiSiteConfigRouter)
		}
	}
	{
		apiArticleRouter := ar.root.Group("/article")
		{
			ar.addRouter(article.NewArticleController(), apiArticleRouter)
			ar.addRouter(article.NewArticleToolBarController(), apiArticleRouter)
			ar.addRouter(article.NewArticleTypeController(), apiArticleRouter)
		}
	}
	{
		apiImgRouter := ar.root.Group("/image")
		{
			ar.addRouter(imgs.NewImgsController(), apiImgRouter)
		}
	}
}
