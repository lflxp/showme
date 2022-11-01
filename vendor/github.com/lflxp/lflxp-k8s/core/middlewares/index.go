package middlewares

import (
	"github.com/gin-gonic/gin"
)

func RegisterIndex(c *gin.Context) {
	c.Redirect(302, "/admin/index")
}
