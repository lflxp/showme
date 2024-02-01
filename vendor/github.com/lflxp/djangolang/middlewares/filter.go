package middlewares

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/lflxp/djangolang/consted"

	// js "github.com/lflxp/djangolang/middlewares/jwt/services"

	"github.com/lflxp/djangolang/utils"

	"github.com/gin-gonic/gin"
)

const (
	LabelValueTenantPlatform = "platform"
	cacheTime                = 15 * time.Second
	XForwardedUri            = "X-Forwarded-Uri"
	XForwardedMethod         = "X-Forwarded-Method"
	XForwardedHost           = "X-Forwarded-Host"
	XForwardedProto          = "X-Forwarded-Proto"
)

// func JwtTokenFilter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if !isWhilteUrl(c) {
// 			token, err := jwts.NewGinJwtMiddlewares(js.AllUserAuthorizator).ParseToken(c)
// 			if err != nil {
// 				c.JSONP(401, gin.H{
// 					"error": err.Error(),
// 					"code":  "token is invaild",
// 				})
// 				return
// 			}

// 			info, err := json.Marshal(token)
// 			if err != nil {
// 				c.JSONP(401, gin.H{
// 					"error": err.Error(),
// 					"code":  "json.Marshal error",
// 				})
// 				return
// 			}

// 			var user admin.User
// 			err = json.Unmarshal(info, &user)
// 			if err != nil {
// 				c.JSONP(401, gin.H{
// 					"error": err.Error(),
// 					"code":  "json.Unmarshal error",
// 					"info":  string(info),
// 				})
// 				return
// 			}

// 			c.Request.Header.Set("username", user.Username)
// 			c.Request.Header.Set("token", user.Token)
// 			c.Next()
// 		} else {
// 			c.Next()
// 		}
// 	}
// }

// func TokenFilter() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// token := jwt.ExtractClaims(c)
// 		// slog.Debug("ExtractClaims token ", token)

// 		if !isWhilteUrl(c) {
// 			user, err := js.ParseJWTToken(c)
// 			if err != nil {
// 				if strings.Contains(err.Error(), "named cookie not present") {
// 					httpclient.SendErrorMessage(c, http.StatusUnauthorized, "token invalid", "/admin/index/")
// 					c.Abort()
// 					return
// 				}
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, httpclient.Result{
// 					Success:      false,
// 					Data:         "Token is not vaild!",
// 					ErrorCode:    "tokenError",
// 					ErrorMessage: err.Error(),
// 					Host:         c.Request.Host,
// 					TraceId:      "",
// 					ShowType:     "filter",
// 				})
// 				return
// 			}

// 			c.Request.Header.Set("username", user.Username)
// 			c.Request.Header.Set("name", user.Name)
// 			c.Request.Header.Set("userid", user.UserId)
// 			c.Request.Header.Set("email", user.Email)
// 			c.Request.Header.Set("token", user.Token)
// 			c.Request.Header.Set("refreshtoken", user.RefreshToken)
// 		} else {
// 			c.Next()
// 		}
// 	}
// }

// `^/api/*`, // for debug only
// `^/ws/*`,
func isWhilteUrl(c *gin.Context) bool {
	var rs bool
	url := []string{
		`^/$`,
		`^/swagger/*`,
		`^/login`,
		`^/adminfs/*`,
		`^/favicon.ico`,
		`^/auth/*`,
		`^/v3/*`,
		`^/metrics`,
		`^/cloud/xterm`,
		`^/cloud/log`,
		`^/dashboard/*`,
		`^/script/*`,
		`^/admin/auth/login`,
		`^/api/login`,
		`^/apis/*`,
		`^/node_modules/*`,
		// `^/monitor/*`,
	}

	for _, x := range url {
		rs, _ = regexp.MatchString(x, c.Request.URL.Path)
		if rs {
			break
		}
	}

	slog.Debug("method [%s] isWhite %v path %s Url.Path %s ", c.Request.Method, rs, c.Request.RequestURI, c.Request.URL.Path)
	return rs
}

// 平台级别白名单列表
var WhilteList []string = []string{
	"URLCACHE_GET", "URLCACHEALL_GET", "URLCACHE_FORWARDAUTH",
}

// var allWhilteList []string = []string{"USER_GET"}

// 定时获取角色列表
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			<-ticker.C
			// TODO: 获取角色信息
			ROLE := "TargetRole"
			// 缓存角色信息
			utils.NewCacheCliWithTTL().Set("ROLE", ROLE, -1)
		}
	}()
}

// 校验角色及角色权限
// 适配Oauth2-Proxy权限
func isRolePermission(c *gin.Context) (bool, error) {
	var (
		group, fullpath, method, urlCode string
		// isPart3                          bool // 是否第三方接口
	)
	for _, v := range utils.CacheUrlAll() {
		if c.FullPath() == URLCACHE_FORWARDAUTH {
			// isPart3 = true
			fullpath = c.Request.Header.Get(XForwardedUri)
			method = c.Request.Header.Get(XForwardedMethod)
			// slog.Debugf("探测鉴权接口: URLCACHE_FORWARDAUTH %s fullpath %s method %s Request %v", c.FullPath(), fullpath, method, c.Request.Header)
			// TODO: 跨平台接口需要支持正则表达式校验，否则接口太多
			found, err := regexp.MatchString(v.Path, fullpath)
			if err != nil {
				slog.Error(err.Error())
				return false, err
			}
			if found && strings.EqualFold(v.Method, method) {
				group = v.Group
				urlCode = v.Name
				slog.Debug("==鉴权接口==》", "URL", fmt.Sprintf("%s://%s%s method %s Name %s", c.GetHeader(XForwardedProto), c.GetHeader(XForwardedHost), fullpath, method, urlCode))
				break
			}
		} else {
			fullpath = c.FullPath()
			method = c.Request.Method
			if v.Path == fullpath && strings.EqualFold(v.Method, method) {
				group = v.Group
				urlCode = v.Name
				break
			}
			// slog.Debugf("AAA => %v", c.Request.Header)
		}
	}

	if group == "" {
		slog.Error(fmt.Sprintf("获取CacheUrl数据失败：[%s %s] url not found", method, fullpath))
		c.Request.Header.Set("ErrorCode", utils.LicenseError)
		c.Request.Header.Set("ErrorMessage", fmt.Sprintf("获取CacheUrl数据失败：[%s %s] url not found", method, fullpath))
		return false, fmt.Errorf(fmt.Sprintf("获取CacheUrl数据失败：[%s %s] url not found", method, fullpath))
	}

	// 白名单放行
	if utils.ContainsString(WhilteList, urlCode) {
		return true, nil
	}

	// TODO: 获取角色和角色的APIS CODE
	// TODO: 校验权限
	currentProject := c.Request.Header.Get("ProjectName")
	// 优先判断平台级角色 没有再判断项目级角色
	var platformRole, projectRole string
	// 缓存用户角色10秒钟
	username := c.Request.Header.Get(consted.Cookie_Username)
	roles, isExist := utils.NewCacheCliWithTTL().Get(fmt.Sprintf("%s-Roles", username))
	if !isExist {
		// TODO: 查询角色绑定信息

		utils.NewCacheCliWithTTL().Set(fmt.Sprintf("%s-Roles", username), roles, cacheTime)
	} else {
		// TODO: 获取角色
		roles := []map[string]string{}
		for _, v := range roles {
			if v["Level"] == "platform" {
				platformRole = v["Role"]
			}

			if v["Project"] == currentProject {
				projectRole = v["Role"]
			}
		}
	}

	isAdmin := c.GetHeader("X-Forwarded-Groups")

	// 如果是超级管理员，则直接放行
	if isAdmin != "admin" {
		isValid := false
		// 验证平台级非超级管理员角色
		if platformRole != "" {
			if PlatformAPIS, ok := utils.NewCacheCliWithTTL().Get(fmt.Sprintf("%s-RRR", platformRole)); !ok {
				slog.Error(fmt.Sprintf("用户 %s 平台级角色 %s 不存在", username, platformRole))
				return false, fmt.Errorf("用户 %s 平台级角色 %s 不存在", username, platformRole)
			} else {
				if utils.ContainsString(PlatformAPIS.([]string), urlCode) {
					isValid = true
				}
			}
		}

		if !isValid {
			if projectRole != "" {
				// 验证项目级角色
				if ProjectAPIS, ok := utils.NewCacheCliWithTTL().Get(fmt.Sprintf("%s-RRR", projectRole)); !ok {
					slog.Error(fmt.Sprintf("用户 %s 项目级角色 %s 不存在", username, projectRole))
					return false, fmt.Errorf("用户 %s 项目级角色 %s 不存在", username, projectRole)
				} else {
					if utils.ContainsString(ProjectAPIS.([]string), urlCode) {
						isValid = true
					} else {
						isValid = false
					}
				}
			}
		}

		if !isValid {
			// c.Request.Header.Set("ErrorCode", utils.AthorizationError)
			// c.Request.Header.Set("ErrorMessage", fmt.Sprintf("当前角色无API访问权限：[%s] %s", c.Request.Method, c.FullPath()))
			return false, fmt.Errorf("当前角色无API访问权限：[%s] %s", method, fullpath)
		}
	}

	return true, nil
}
