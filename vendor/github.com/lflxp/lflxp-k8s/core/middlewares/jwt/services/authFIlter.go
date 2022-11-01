package services

import (
	"github.com/gin-gonic/gin"
)

// 平台级别白名单列表
var whilteList []string = []string{"TENANTALL_GET"}
var allWhilteList []string = []string{"USER_GET"}

// username is test can access
func PaasAuthorizator(data interface{}, c *gin.Context) bool {
	// if v, ok := data.(*jwtmodel.User); ok && v.Username == "test" {
	// 	return true
	// }

	// return false
	return true
}

// license和角色权限验证
// 传递错误用Header ErrorCode|ErrorMessage
func AllUserAuthorizator(data interface{}, c *gin.Context) bool {
	// TODO: Add Filter

	return true
}

func AdminAuthorizator(data interface{}, c *gin.Context) bool {
	// TODO: Add Filter

	return true
}
