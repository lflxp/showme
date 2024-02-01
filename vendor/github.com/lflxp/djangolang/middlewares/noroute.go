package middlewares

import (
	"github.com/gin-gonic/gin"
)

// 默认404跳转地址
var NoRoutePath string = "/admin/index"

// 404 handler
func NoRouteHandler(c *gin.Context) {
	// c.HTML(404, "layout/404.html", nil)
	// c.String(http.StatusNotFound, "%s", "Page Not Found")
	c.Redirect(301, NoRoutePath)
}
