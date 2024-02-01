package middlewares

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/guonaihong/gout"

	. "github.com/lflxp/djangolang/model/admin"
	"github.com/lflxp/djangolang/utils"

	"github.com/lflxp/djangolang/utils/httpclient"
	"github.com/lflxp/djangolang/utils/orm/sqlite"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	identityKey    = "thisiskeyword"
	IsRancherLogin bool
)

// 用户表
// type Auth struct {
// 	Id           int64  `json:"id"`
// 	Username     string `form:"username" json:"username" binding:"required" xorm:"varchar(255) notnull index unique"`
// 	Password     string `form:"password" json:"password" binding:"required" xorm:"varchar(255) not null"`
// 	Name         string `json:"name" xorm:"varchar(255)"`
// 	Avatar       string `json:"avatar" xorm:"varchar(255)"`
// 	Status       bool   `json:"status" xorm:"bool"`
// 	Telephone    string `json:"telephone" xorm:"varchar(255)"`
// 	LastLoginIp  string `json:"lastLoginIp" xorm:"varchar(255)"`
// 	CreateTime   string `json:"createTime" xorm:"varchar(255)"`
// 	CreatorId    string `json:"creatorId" xorm:"varchar(255)"`
// 	MerchantCode string `json:"merchantCode" xorm:"varchar(255)"`
// 	Deleted      bool   `json:"deleted" xorm:"bool"`
// 	RoleId       string `json:"roleId" xorm:"varchar(255)"`
// }
// type User struct {
// 	Username string   `json:"username"`
// 	Password string   `json:"password"`
// 	Claims   []Claims `json:"claims"`
// }

// // 用户权限表
// type Claims struct {
// 	Id    int64  `json:"id"`
// 	Auth  string `json:"auth" xorm:"varchar(255) unique(only)"`  // 对应Auth => Username  eg: admin
// 	Type  string `json:"type" xorm:"varchar(255) unique(only)"`  // 权限类型 eg: nav
// 	Value string `json:"value" xorm:"varchar(255) unique(only)"` // 权限指 eg: dashboard
// }

// 获取用户权限列表
func GetUserClaims(username string) ([]Claims, error) {
	var result []Claims

	sql := fmt.Sprintf("select * from user where username='%s'", username)
	rs, err := sqlite.NewOrm().Query(sql)
	if err != nil {
		return nil, err
	}

	if len(rs) <= 0 {
		// TODO: finished get user claims
		// if account := viper.GetStringMap("account"); account != nil {
		// 	if user, ok := account[username]; !ok {
		// 		return nil, errors.New("查无此人")
		// 	} else {
		// 		if claim, ok := user.(map[string]interface{})["claim"]; ok {
		// 			// TODO: password md5 jiami
		// 			err := json.Unmarshal([]byte(claim.(string)), &result)
		// 			if err != nil {
		// 				return result, err
		// 			}
		// 			return result, nil
		// 		}
		// 	}
		// }
		return nil, errors.New("用户不存在")
	} else {
		ttt := strings.Split(string(rs[0]["claims_id"]), ",")
		err = sqlite.NewOrm().In("id", ttt).Find(&result)
		return result, err
	}

	return result, errors.New("nothing found")
}

// 验证用户密码
func VerifyAuth(username, password string) (bool, error) {
	// 优先查询数据库
	var user = User{Username: username}
	// 忽略[]claims与string 解析
	has, _ := sqlite.NewOrm().Get(&user)

	if has {
		slog.Debug("Found User %s Login", username)
		slog.Info("password %s %s %s", user.Password, utils.MD5(password), password)
		// if user.Password == utils.MD5(password) {
		if user.Password == password {
			return true, nil
		} else {
			return false, errors.New("用户名或密码错误")
		}
	} else if username == "admin" {
		// TODO: verify username pwd
		// if account := viper.GetStringMap("account"); account != nil {
		// 	if user, ok := account[username]; !ok {
		// 		return false, errors.New("查无此人")
		// 	} else {
		// 		if pwd, ok := user.(map[string]interface{})["password"]; ok {
		// 			// TODO: password md5 jiami
		// 			if pwd.(string) == password {
		// 				return true, nil
		// 			} else {
		// 				return false, errors.New("用户名或密码错误！")
		// 			}
		// 		}
		// 	}
		// }
		if password == "admin" {
			return true, nil
		} else {
			return false, errors.New("用户名或密码错误！")
		}
	}

	return false, errors.New("account is not exists")
}

func VerifyAuthByRancher(username, password string) (string, string, bool) {
	// rancherHost := viper.GetString("rancher.host")
	rancherHost := "127.0.0.1"
	slog.Info("初始化RancherHost地址 %s", rancherHost)
	url := fmt.Sprintf("%s/v3-public/localProviders/local?action=login", rancherHost)
	slog.Debug("url is %s", url)
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
		slog.Error(err.Error())
		return "", "", false
	}

	if code != 200 {
		slog.Error(body)
		return "", "", false
	}

	ress := header.Get("Set-Cookie")
	if ress == "" {
		slog.Error("no token found")
		return "", "", false
	}

	slog.Debug("Ress is %s", ress)
	var rs string
	for _, t := range strings.Split(ress, ";") {
		if strings.Contains(t, "R_SESS") {
			tmp := strings.Split(t, "=")
			rs = tmp[1]
		}
	}
	if rs == "" {
		slog.Error("no R_SESS found")
		return "", "", false
	}
	return rs, ress, true
}

type IdentityKey string

// 接口权限
type JwtAuthorizator func(data interface{}, c *gin.Context) bool

// 根据不同接口的权限规则生成不同权限的jwt中间件
func NewGinJwtMiddlewares(jwta JwtAuthorizator) *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "apizone",
		Key:         []byte("demo"),
		Timeout:     7 * 24 * time.Hour,
		MaxRefresh:  7 * 24 * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				// claimJson, err := GetUserClaims(v.Username)
				// if err != nil {
				// 	slog.Error(err.Error())
				// }
				// maps the claims in the JWT
				return jwt.MapClaims{
					"username": v.Username,
					"token":    v.Token,
				}
			} else if vData, ok := data.(map[string]string); ok {
				tmp := jwt.MapClaims{}
				for k, v := range vData {
					tmp[k] = v
				}
				return tmp
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			// slog.Error("claims ", claims)
			//extracts identity from claims
			// if claims["userName"] != nil {
			// 	return &User{
			// 		Username: claims["userName"].(string),
			// 	}
			// }
			// jsonClaim := claims["userClaims"].([]interface{})
			// data, err := json.Marshal(jsonClaim)
			// if err != nil {
			// 	slog.Error(err)
			// 	return &User{
			// 		Username: claims["userName"].(string),
			// 	}
			// }

			// var userClaims []Claims
			// json.Unmarshal(data, &userClaims)
			// //Set the identity
			// return &User{
			// 	Username: claims["userName"].(string),
			// 	Claims:   userClaims,
			// }
			slog.Debug("-----------------CLAIMS %v", claims)
			return claims
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			//handles the login logic. On success LoginResponse is called, on failure Unauthorized is called
			var loginVals User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			slog.Debug("JWT Authenticator %v username %s password %s", loginVals, userID, password)
			// if IsRancherLogin {
			// 	if token, ress, ok := VerifyAuthByRancher(userID, password); ok {
			// 		c.Header("Set-Cookie", strings.Replace(ress, "Secure", "", -1))
			// 		return &User{
			// 			Username: userID,
			// 			Token:    token,
			// 		}, nil
			// 	}
			// } else {
			// 	if ok, err := VerifyAuth(userID, password); ok {
			// 		return &User{
			// 			Username: userID,
			// 			Token:    "verifyAuth",
			// 		}, nil
			// 	} else {
			// 		slog.Error(err.Error())
			// 	}
			// }
			if ok, err := VerifyAuth(userID, password); ok {
				return &User{
					Username: userID,
					Token:    "verifyAuth",
				}, nil
			} else {
				slog.Error(err.Error())
			}

			return nil, jwt.ErrFailedAuthentication
		},
		//receives identity and handles authorization logic
		Authorizator: jwta,
		//handles unauthorized logic
		Unauthorized: func(c *gin.Context, code int, message string) {
			// c.JSON(code, gin.H{
			//	"code":    code,
			//	"message": message,
			// })
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
		CookieHTTPOnly: true,  // js can't modify
		// CookieDomain: "localhost:8080",
		CookieName:     "token",
		CookieSameSite: http.SameSiteDefaultMode,
		LogoutResponse: func(c *gin.Context, code int) {
			c.Redirect(http.StatusPermanentRedirect, "/login")
		},
	})
	if err != nil {
		slog.Error("JWT Error:" + err.Error())
	}
	return authMiddleware
}

// role is admin can access
func AdminAuthorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok {
		// for _, itemClaim := range v.Claims {
		// 	if itemClaim.Type == "role" && itemClaim.Value == "admin" {
		// 		return true
		// 	}
		// }

		if v.Username == "admin" {
			return true
		}
	}

	return false
}

// username is test can access
func TestAuthorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Username == "test" {
		return true
	}

	return false
}

// 不限制用户权限
func AllUserAuthorizator(data interface{}, c *gin.Context) bool {
	return true
}
