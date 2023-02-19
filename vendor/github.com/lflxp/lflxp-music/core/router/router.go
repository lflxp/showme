package router

import (
	"github.com/lflxp/lflxp-music/core/controller"
	"github.com/lflxp/lflxp-music/core/middlewares"
	"github.com/lflxp/lflxp-music/core/pages"
	"github.com/lflxp/lflxp-music/core/pkg/auth"
	"github.com/lflxp/lflxp-music/core/pkg/music"
	"github.com/lflxp/lflxp-music/core/pkg/proxy"

	"github.com/lflxp/lflxp-music/core/asset"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
)

// 注册插件和路由
func PreGinServe(r *gin.Engine) {
	log.Info("注册Gin路由")

	r.Use(middlewares.TokenFilter())
	r.Use(gin.Logger())

	r.Use(middlewares.Cors())

	// 添加prometheus监控
	middlewares.RegisterPrometheusMiddleware(r, false)

	// gzip
	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{".*"})))

	// 404
	r.NoRoute(middlewares.NoRouteHandler)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/music")
	})

	// 健康检查
	r.GET("/health", middlewares.RegisterHealthMiddleware)

	// Swagger
	middlewares.RegisterSwaggerMiddleware(r)

	// 登陆
	auth.RegisterAuth(r)
	// auth.RegisterShop(r)

	pages.RegisterTemplate(r)
	controller.RegisterAdmin(r)
	asset.RegisterAsset(r)
	proxy.ProxyRegister(r)
	music.RegisterMusic(r)
}
