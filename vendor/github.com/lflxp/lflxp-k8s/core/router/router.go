package router

import (
	"github.com/lflxp/lflxp-k8s/core/controller"
	"github.com/lflxp/lflxp-k8s/core/middlewares"
	"github.com/lflxp/lflxp-k8s/core/pages"
	"github.com/lflxp/lflxp-k8s/pkg/apiserver"

	"github.com/lflxp/lflxp-k8s/asset"
	"github.com/lflxp/lflxp-k8s/pkg/auth"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
)

// 注册插件和路由
func PreGinServe(r *gin.Engine) {
	log.Info("注册Gin路由")

	r.Use(middlewares.TokenFilter())
	// r.Use(gin.Logger())
	// TODO: 自定义Panic 处理
	// r.Use(middlewares.PanicAdviceMiddleware())
	//r.Use(gin.RecoveryWithWriter(os.Stdout))
	r.Use(middlewares.Cors())

	// 添加prometheus监控
	middlewares.RegisterPrometheusMiddleware(r, false)

	// gzip
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{".*"})))

	// 404
	r.NoRoute(middlewares.NoRouteHandler)

	r.GET("/", func(c *gin.Context) {
		// c.Redirect(301, "/login")
		// c.Redirect(301, "/dashboard")
		c.Redirect(301, "/d2admin")
		// c.Redirect(301, "/swaggers/index.html")
	})

	// 健康检查
	r.GET("/health", middlewares.RegisterHealthMiddleware)

	// Swagger
	middlewares.RegisterSwaggerMiddleware(r)

	// 自动注册，请勿修改

	// demo.Registerdemo(r)

	// url接口
	// urlcache.RegisterUrlCache(r)

	// 登陆
	auth.RegisterAuth(r)

	asset.RegisterAsset(r)
	pages.RegisterTemplate(r)
	controller.RegisterAdmin(r)
	controller.Registertest(r)
	apiserver.RegisterApiserver(r)
	apiserver.RegisterApiserverWS(r)
}
