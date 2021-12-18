package roles

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

type RoleType string

const (
	CopyType     RoleType = "copy"
	FetchType    RoleType = "fetch"
	DockerType   RoleType = "image"
	K8sType      RoleType = "k8s"
	ShellType    RoleType = "shell"
	SwarmType    RoleType = "swarm"
	TemplateType RoleType = "template"
	IncludeType  RoleType = "include"
	SystemdType  RoleType = "systemd"
	UserType     RoleType = "user"
	ScriptType   RoleType = "script"
)

// roles 插件模块名称与类型的映射关系，主要用于通过反射创建roles对象
var roleNames = make(map[RoleType]reflect.Type)

func addRoles(key RoleType, value reflect.Type) {
	log.Debugf("添加Roles插件 %s", key)
	roleNames[key] = value
}
