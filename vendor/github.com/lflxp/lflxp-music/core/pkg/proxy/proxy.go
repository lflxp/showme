package proxy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lflxp/tools/proxy"
)

func ProxyRegister(router *gin.Engine) {
	proxyGroup := router.Group("")
	// proxyGroup.Use(jwts.NewGinJwtMiddlewares(jwtservice.PaasAuthorizator).MiddlewareFunc())
	//proxyGroup.Use(jwts.GetMiddleware().MiddlewareFunc())
	// 代理
	{
		proxyGroup.GET("/playlist/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/playlist", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/lyric/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/lyric", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/toplist/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/toplist", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/personalized/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/personalized", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/search/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/search", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/song/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/song", "https://mu-api.yuk0.com"), nil))
		proxyGroup.GET("/comment/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/comment", "https://mu-api.yuk0.com"), nil))
	}
}
