package auth

import (
	"time"

	jwts "github.com/lflxp/lflxp-k8s/core/middlewares/jwt/framework"
	"github.com/lflxp/lflxp-k8s/utils"

	"github.com/gin-gonic/gin"
)

var (
	AUTH_LOGOUT_POST      string = "/admin/auth/logout"
	AUTH_LOGIN_LOCAL_POST string = "/admin/auth/login"
)

func RegisterAuth(router *gin.Engine) {
	keycloakGroup := router.Group("")
	{
		keycloakGroup.GET("/api/tablelist", tablelist)
		keycloakGroup.POST("/api/login", loginlocal)
		keycloakGroup.POST(AUTH_LOGIN_LOCAL_POST, loginlocal)
		keycloakGroup.POST(AUTH_LOGOUT_POST, logout)
	}
}

func tablelist(c *gin.Context) {
	tmp := []map[string]interface{}{}

	status := []string{"published", "draft", "deleted"}
	for i := 0; i < 30; i++ {
		current := i % 3
		tmp = append(tmp, map[string]interface{}{
			"id":           i,
			"title":        utils.GetRandomString(i + 1),
			"status":       status[current],
			"author":       utils.GetRandomString(10),
			"display_time": time.Now().String(),
			"pageviews":    i,
		})
	}
	c.JSONP(200, gin.H{
		"data": map[string]interface{}{
			"total": 30,
			"items": tmp,
		},
	})
}

// @Summary  本地登录接口
// @Description 后端服务登录接口
// @Tags Auth
// @Success 200 {object} model.Resp{}
// @Security ApiKeyAuth
// @Router /apis/auth/login/local [post]
func loginlocal(c *gin.Context) {
	jwts.GetMiddleware().LoginHandler(c)
}

// @Summary  注销接口
// @Description 后端服务注销接口
// @Tags Auth
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /apis/auth/logout [post]
func logout(c *gin.Context) {
	jwts.GetMiddleware().LogoutHandler(c)
}
