// prometheus中间件
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 默认prometheus基础监控
// func RegisterPrometheusMiddleware(router *gin.Engine, isauth bool) {
// 	if isauth {
// 		group := router.Group("/metrics", gin.BasicAuth(gin.Accounts{
// 			"root": "system",
// 		}))
// 		{
// 			group.GET("/", ginprom.PromHandler(promhttp.Handler()))
// 			// group.GET("/")
// 		}
// 	} else {
// 		router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
// 	}
// }

func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// 默认prometheus监控+自定义监控
func RegisterPrometheusMiddleware(router *gin.Engine, isauth bool) {
	if isauth {
		group := router.Group("/metrics", gin.BasicAuth(gin.Accounts{
			"root": "system",
		}))
		{
			group.GET("/", PromHandler(promhttp.HandlerFor(
				prometheus.DefaultGatherer,
				promhttp.HandlerOpts{
					// Opt into OpenMetrics to support exemplars.
					EnableOpenMetrics: true,
				},
			)))
			// group.GET("/")
		}
	} else {
		router.GET("/metrics", PromHandler(promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				EnableOpenMetrics: true,
			},
		)))
	}
}
