package framework

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lflxp/lflxp-music/core/middlewares/jwt/model"
	"github.com/lflxp/lflxp-music/core/middlewares/jwt/services"

	"github.com/lflxp/tools/httpclient"
	"github.com/lflxp/tools/utils"

	"github.com/guonaihong/gout"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/tools/orm/sqlite"
	"github.com/spf13/viper"
)

var (
	identityKey    string
	jwtMiddleware  *jwt.GinJWTMiddleware
	IsRancherLogin bool
)

func init() {
	identityKey = viper.GetString("token.jwtIdentityKey")
	// log.Debugf("初始化jwtIdentityKey %s", identityKey)
}

func GetMiddleware() *jwt.GinJWTMiddleware {
	if jwtMiddleware == nil {
		jwtMiddleware = NewGinJwtMiddlewares(services.AllUserAuthorizator)
	}
	return jwtMiddleware
}

// 验证用户密码
func VerifyAuth(username, password string) (bool, error) {
	// 优先查询数据库
	var user = new(model.User)
	// 忽略[]claims与string 解析
	has, err := sqlite.NewOrm().Where("username=?", username).Get(user)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if has {
		log.Debugf("Found User %s Login", username)
		if user.Password == password {
			return true, nil
		} else {
			return false, errors.New("用户名或密码错误")
		}
	} else {
		// TODO: verify username pwd
		if account := viper.GetStringMap("account"); account != nil {
			if user, ok := account[username]; !ok {
				return false, errors.New("查无此人")
			} else {
				if pwd, ok := user.(map[string]interface{})["password"]; ok {
					// TODO: password md5 jiami
					if pwd.(string) == password {
						return true, nil
					} else {
						return false, errors.New("用户名或密码错误！")
					}
				}
			}
		} else {
			return false, errors.New("用户名或密码错误")
			// 如果又没数据库又没yaml配置文件
			// 默认密码： admin
			// if user.Password == "admin" {
			// 	return true, nil
			// } else {
			// 	return false, errors.New("用户名或密码错误")
			// }
		}
	}

	return false, errors.New("account is not exists")
}

// 验证kubevela VelaUX OpenAPI
// https://kubevela.net/zh/docs/platform-engineers/openapi/overview
func VerifyKubeVelaLogin(username, password string) (*model.KubeVelaToken, error) {
	var result model.KubeVelaToken
	host := viper.GetString("auth.url")
	url := fmt.Sprintf("%s/api/v1/auth/login", host)
	log.Debug("kubevela login url: %s", url)
	code := 0
	body := ""
	err := httpclient.NewGoutClient().
		POST(url).
		Debug(true).
		SetHeader(gout.H{
			"Content-Type": "application/json",
		}).
		SetJSON(gout.H{
			"username": username,
			"password": password,
		}).
		BindBody(&body).
		Code(&code).
		Do()
	if err != nil {
		return &result, err
	}

	if code != 200 {
		log.Error(body)
		return &result, errors.New(body)
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

func VerifyAuthByRancher(username, password string) (string, string, bool) {
	rancherHost := viper.GetString("proxy.host")
	log.Infof("初始化RancherHost地址 %s", rancherHost)
	url := fmt.Sprintf("%s/v3-public/localProviders/local?action=login", rancherHost)
	log.Debugf("url is %s", url)
	code := 0
	body := ""
	header := make(http.Header)
	err := httpclient.NewGoutClient().
		POST(url).
		Debug(true).
		SetHeader(gout.H{
			"Content-Type": "application/json;charset=UTF-8",
			"Accept":       "application/json",
		}).
		SetJSON(gout.H{
			"description":  "APILogin",
			"password":     password,
			"username":     username,
			"responseType": "cookie",
		}).
		BindBody(&body).
		BindHeader(&header).
		Code(&code).
		Do()

	if err != nil {
		log.Error(err.Error())
		return "", "", false
	}

	if code != 200 {
		log.Error(body)
		return "", "", false
	}

	ress := header.Get("Set-Cookie")
	if ress == "" {
		log.Error("no token found")
		return "", "", false
	}

	log.Debugf("Ress is %s", ress)
	var rs string
	for _, t := range strings.Split(ress, ";") {
		if strings.Contains(t, "R_SESS") {
			tmp := strings.Split(t, "=")
			rs = tmp[1]
		}
	}
	if rs == "" {
		log.Error("no R_SESS found")
		return "", "", false
	}
	return rs, ress, true
}

var (
	keyOnce sync.Once
	keys    string
)

func newKey() string {
	keyOnce.Do(func() {
		keys = utils.GetRandomString(64)
	})
	return keys
}

type IdentityKey string

// 接口权限
type JwtAuthorizator func(data interface{}, c *gin.Context) bool

// 根据不同接口的权限规则生成不同权限的jwt中间件
func NewGinJwtMiddlewares(jwta JwtAuthorizator) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: "apizone",
		// Key:         []byte(newKey()), // 每次随机不利于k8s无状态部署
		Key:         []byte("9f5b3737-ad22-4426-a517-6877618813a4"), // 每次随机不利于k8s无状态部署
		Timeout:     10 * time.Hour,
		MaxRefresh:  10 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				// claimJson, err := GetUserClaims(v.Username)
				// if err != nil {
				// 	log.Error(err.Error())
				// 	return jwt.MapClaims{}
				// }
				// maps the claims in the JWT
				return jwt.MapClaims{
					"username":     v.Username,
					"token":        v.Token,
					"email":        v.Email,
					"tenant":       v.Tenant,
					"authProvider": v.AuthProvider,
					"userId":       v.UserId,
					"role":         v.Role,
					"roleLevel":    v.RoleLevel,
					"roleReal":     v.RoleReal,
					"isGlobal":     v.IsGlobal,
				}
			} else if vv, ok := data.(map[string]string); ok {
				tmp := jwt.MapClaims{}
				for k, v := range vv {
					tmp[k] = v
				}
				return tmp
			} else if vela, ok := data.(*model.KubeVelaToken); ok {
				return jwt.MapClaims{
					"username":           vela.User.Name,
					"token":              vela.AccessToken,
					"refreshToken":       vela.RefreshToken,
					"createTimestamp":    vela.User.CreateTime,
					"lastLoginTimestamp": vela.User.LastLoginTime,
					"email":              vela.User.Email,
					"disabled":           vela.User.Disabled,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Debugf("I Claims: %v", claims)
			return nil
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals model.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			c.Request.Header.Set("user", loginVals.Username)

			log.Debugf("JWT Authenticator %v username %s password %s", loginVals, userID, password)
			if IsRancherLogin {
				if token, ress, ok := VerifyAuthByRancher(userID, password); ok {
					c.Header("Set-Cookie", strings.Replace(ress, "Secure", "", -1))
					return &model.User{
						Username: userID,
						Token:    token,
					}, nil
				}
			} else {
				isDev := viper.GetBool("auth.dev")
				if isDev {
					if token, err := VerifyKubeVelaLogin(userID, password); err == nil {
						return token, nil
					} else {
						log.Error(err.Error())
					}
				} else {
					if ok, err := VerifyAuth(userID, password); ok {
						return &model.User{
							Username: userID,
							Token:    "verifyAuth",
						}, nil
					} else {
						log.Error(err.Error())
					}
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},
		//receives identity and handles authorization logic
		Authorizator: jwta,
		//handles unauthorized logic
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.Redirect(http.StatusFound, "/login?url="+c.Request.RequestURI)
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:       time.Now,
		SendCookie:     true,
		SecureCookie:   false, // non HTTPS dev environments
		CookieHTTPOnly: false, // js can't modify
		// CookieDomain: "localhost:8080",
		CookieName:     "token",
		CookieSameSite: http.SameSiteLaxMode,
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			c.JSONP(code, gin.H{
				"token":  message,
				"time":   time.String(),
				"code":   code,
				"name":   c.Request.Header.Get("user"),
				"uuid":   newKey(),
				"expire": time,
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.SetCookie("token", "", -1, "/", "localhost", false, true)
			c.Redirect(302, "/login")
		},
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return authMiddleware
}
