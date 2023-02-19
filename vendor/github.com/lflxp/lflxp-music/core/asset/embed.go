package asset

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed music
var music embed.FS

func RegisterAsset(router *gin.Engine) {
	router.GET("/music/*any", func(c *gin.Context) {
		// staticServer := wrapHandler(http.FS(dashboard))
		staticServer := http.FileServer(http.FS(music))
		// TODO: 遇到404 就返回前端根路径下index.html的资源 路径不变
		staticServer.ServeHTTP(c.Writer, c.Request)
	})

	// router.StaticFS("/favicon.ico", http.FS(favicon))
}
