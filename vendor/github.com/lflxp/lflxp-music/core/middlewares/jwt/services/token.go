package services

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	log "github.com/go-eden/slf4go"
	"github.com/lflxp/lflxp-music/core/middlewares/jwt/model"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// 获取用户详细信息
func ParseToken(info string, sc *corev1.Secret) (*model.User, error) {
	// TODO: 填补用户信息
	jwtUser := &model.User{}

	return jwtUser, nil
}

// 解析JWT Token
func ParseJWTToken(c *gin.Context) (*model.User, error) {
	jwtoken, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}

	info := strings.Split(jwtoken, ".")[1]

	payload, err := base64.RawURLEncoding.DecodeString(info)
	if err != nil {
		return nil, err
	}

	var user *model.User
	err = json.Unmarshal(payload, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type JwtData struct {
	Exp       int64 `json:"exp,omitempty"`
	Iat       int64 `json:"iat,omitempty"`
	Auth_time int64 `json:"auth_time,omitempty"`
}

// 判断jwt是否过期
func IsExpiresHeader(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return false
	}

	info := strings.Split(token, ".")[1]

	payload, err := base64.RawURLEncoding.DecodeString(info)
	if err != nil {
		log.Errorf("error decoding payload %s: %v", info, err)
		return false
	}

	var data JwtData
	err = json.Unmarshal(payload, &data)
	if err != nil {
		log.Errorf("error Unmarshal payload %s", err.Error())
		return false
	}

	now := time.Now()

	t := time.Unix(data.Exp, 0)

	if now.Before(t) {
		log.Debugf("now %s before exp %s", now.String(), t.String())
		return true
	}
	log.Debugf("ERROR: now %s after exp %s", now.String(), t.String())
	return false
}
