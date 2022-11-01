package services

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/lflxp/lflxp-k8s/core/middlewares/jwt/model"

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
