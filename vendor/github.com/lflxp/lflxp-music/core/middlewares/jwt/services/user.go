package services

import (
	"github.com/lflxp/lflxp-music/core/middlewares/jwt/model"

	"github.com/lflxp/tools/httpclient"
	"github.com/lflxp/tools/rsa"
	"github.com/lflxp/tools/utils"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
)

func DecodeUser(c *gin.Context, data map[string]interface{}) (*model.Auth, error) {
	auth := &model.Auth{}

	if data != nil {
		if username, ok := data["username"]; ok {
			auth.Username = username.(string)
		}
		if password, ok := data["password"]; ok {
			auth.Password = password.(string)
		}
		if description, ok := data["description"]; ok {
			auth.Description = description.(string)
		}
		if responseType, ok := data["responseType"]; ok {
			auth.ResponseType = responseType.(string)
		}
	} else {
		err := c.BindJSON(auth)
		if err != nil {
			log.Error(err)
			httpclient.SendErrorMessage(c, 500, "500", err.Error())
			return nil, err
		}
	}

	// 解密
	userBytes, _ := utils.DecodeBase64(auth.Username)
	user, err := rsa.RsaDecrypt([]byte(userBytes))
	if err != nil {
		httpclient.SendErrorMessage(c, 500, "500", err.Error())
		return nil, err
	}
	auth.Username = string(user)

	pwdBytes, _ := utils.DecodeBase64(auth.Password)
	pwd, err := rsa.RsaDecrypt([]byte(pwdBytes))
	if err != nil {
		httpclient.SendErrorMessage(c, 500, "500", err.Error())
		return nil, err
	}
	auth.Password = string(pwd)
	return auth, nil
}
