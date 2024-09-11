package router

import (
	"cmsApp/configs"
	"cmsApp/internal/controllers"
	"cmsApp/pkg/cors"
	"cmsApp/pkg/utils/stringx"
	"cmsApp/pkg/validator"
	"cmsApp/web"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// 跨域
func Cors() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	//corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowHeaders = []string{"Origin", "Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	return cors.New(corsConfig)
}

func Init() (*Router, error) {

	router := NewRouter(gin.Default())

	validator.InitCustomValidator("zh")

	//设置404错误处理
	router.SetRouteError(controllers.NewHandleController().Handle)

	//设置全局中间件
	router.SetGlobalMiddleware(Cors())
	//router.SetGlobalMiddleware(middleware.Trace(), medium.GinLog(facade.NewLogger("logs"), time.RFC3339, false), medium.RecoveryWithLog(facade.NewLogger("logs"), true))

	// 开发模式设置接口文档路由
	if gin.Mode() == "debug" {
		router.SetSwaagerHandle("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 设置模板解析函数
	render, err := web.LoadTemplates()
	if err != nil {

		return nil, err
	}
	router.SetHtmlRenderer(render)

	//设置静态资源
	router.SetStaticFile("/statics", web.StaticsFs)

	//设置上传附件
	uploadPath := stringx.JoinStr(configs.RootPath, string(filepath.Separator), "uploadfile")
	err = router.SetUploadDir(uploadPath)
	if err != nil {
		return nil, err
	}

	// 设置后台全局中间件
	//store := cookie.NewStore([]byte("1GdFRMs4fcWBvLXT"))
	//router.SetAdminRoute(NewAdminRouter(), gzip.Gzip(gzip.DefaultCompression), sessions.Sessions("mysession", store))

	router.SetApiRoute(NewApiRouter())
	return &router, nil
}
