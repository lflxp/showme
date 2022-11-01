package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary  健康检查
// @Description 接口健康检查接口
// @Tags Health
// @Success 200 {string} string "success"
// @Router /health [get]
func RegisterHealthMiddleware(c *gin.Context) {
	c.String(http.StatusOK, "success")
}
