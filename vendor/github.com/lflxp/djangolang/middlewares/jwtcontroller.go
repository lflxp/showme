package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册auth
func RegisterJWT(router *gin.Engine) {
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "admin/login.html", gin.H{
			"User":     "admin",
			"Password": "admin",
		})
	})
	apiGroup := router.Group("/auth")
	// login
	apiGroup.POST("/login", Login)
	apiGroup.GET("/logout", Logout)

	var authMiddleware = NewGinJwtMiddlewares(AllUserAuthorizator)
	authGroup := router.Group("/auth")
	authGroup.Use(authMiddleware.MiddlewareFunc())
	{
		// Refresh time can be longer than token timeout
		authGroup.GET("/refreshtoken", RefreshToken)
	}
}

// @Summary  通用接口
// @Description 登陆、swagger、注销、404等
// @Tags Auth
// @Param token query string false "token"
// @Param data body User true "data"
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var authMiddleware = NewGinJwtMiddlewares(AllUserAuthorizator)
	authMiddleware.LoginHandler(c)
}

// @Summary  通用接口
// @Description 登陆、swagger、注销、404等
// @Tags Auth
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /auth/logout [get]
func Logout(c *gin.Context) {
	var authMiddleware = NewGinJwtMiddlewares(AllUserAuthorizator)
	authMiddleware.LogoutHandler(c)
	c.Redirect(http.StatusFound, "/login")
}

// @Summary  通用接口
// @Description 登陆、swagger、注销、404等
// @Tags Auth
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /auth/refreshtoken [get]
func RefreshToken(c *gin.Context) {
	var authMiddleware = NewGinJwtMiddlewares(AllUserAuthorizator)
	authMiddleware.RefreshHandler(c)
}
