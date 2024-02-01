package middlewares

import (
	"fmt"
	"net/http"

	"github.com/lflxp/djangolang/utils"

	urlcache "github.com/lflxp/djangolang/model"

	"github.com/gin-gonic/gin"
)

var (
	URLCACHE_GET         string = "/apis/v1/urlcache/:name"
	URLCACHEALL_GET      string = "/apis/v1/urlcache"
	URLCACHE_FORWARDAUTH string = "/forwardauth"
)

func init() {
	// 注册URL接口信息
	utils.CacheUrlSet(urlcache.UrlCache{
		Name:        "URLCACHE_GET",
		Method:      http.MethodGet,
		Description: "查询指定URLCache",
		Path:        URLCACHE_GET,
		Group:       "other",
	})
	utils.CacheUrlSet(urlcache.UrlCache{
		Name:        "URLCACHEALL_GET",
		Method:      http.MethodGet,
		Description: "查询所有URLCache",
		Path:        URLCACHEALL_GET,
		Group:       "other",
	})
	utils.CacheUrlSet(urlcache.UrlCache{
		Name:        "URLCACHE_FORWARDAUTH",
		Method:      http.MethodGet,
		Description: "第三方权限校验",
		Path:        URLCACHE_FORWARDAUTH,
		Group:       "other",
	})
}

// 注册UrlCache接口
func RegisterUrlCache(router *gin.Engine) {
	proxyGroup := router.Group("")
	// TODO: 后期生产直接取消接口
	// proxyGroup.Use(jwts.GetMiddleware().MiddlewareFunc())
	// 代理
	{
		proxyGroup.GET(URLCACHE_GET, getByName)
		proxyGroup.GET(URLCACHEALL_GET, getall)
		proxyGroup.GET(URLCACHE_FORWARDAUTH, func(c *gin.Context) {
			// log.Info("===============FORWARDAUTH=============")
			utils.SendSuccessMessage(c, 200, "success")
		})
	}
}

// @Summary  获取指定的URL
// @Description 根据key获取指定的URL
// @Tags UrlCache
// @Param name path string true "URL Key名称"
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /apis/v1/urlcache/{name} [get]
func getByName(c *gin.Context) {
	name := c.Param("name")
	url, ok := utils.CacheUrlGet(name)
	if !ok {
		utils.SendErrorMessage(c, 404, "not found", fmt.Sprintf("%s key not found", name))
	} else {
		utils.SendSuccessMessage(c, 200, url)
	}
}

// @Summary  获取所有URL
// @Description 所有所有URL信息
// @Tags UrlCache
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /apis/v1/urlcache [get]
func getall(c *gin.Context) {
	utils.SendSuccessMessage(c, 200, utils.CacheUrlAll())
}
