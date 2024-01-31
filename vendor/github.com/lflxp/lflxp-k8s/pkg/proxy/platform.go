package proxy

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RegisterPlatform(router *gin.Engine) {
	springboot := viper.GetString("proxy.springboot")
	// prometheus := viper.GetString("proxy.prometheus")
	proxyGroup := router.Group("")
	// proxyGroup.Use(jwts.NewGinJwtMiddlewares(jwtservice.PaasAuthorizator).MiddlewareFunc())
	//proxyGroup.Use(jwts.GetMiddleware().MiddlewareFunc())
	// 代理
	{
		proxyGroup.Any("/greeting/*action", NewHttpProxyByGinCustom(fmt.Sprintf("%s/greeting", springboot), nil))
		// proxyGroup.Any("/prometheus/*action", NewHttpProxyByGinCustom(prometheus, nil))
		// proxyGroup.Any("/public/*action", NewHttpProxyByGinCustom(fmt.Sprintf("%s/public", prometheus), nil))
		// proxyGroup.Any("/oauth/*action", NewHttpProxyByGinCustom(gateway, nil))
	}
}
