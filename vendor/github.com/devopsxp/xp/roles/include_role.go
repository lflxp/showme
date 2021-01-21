package roles

import (
	"fmt"
	"reflect"
	"time"
)

func init() {
	// 初始化include注册模块
	addRoles(IncludeType, reflect.TypeOf(IncludeRole{}))
}

type IncludeRole struct {
	RoleLC
	include string // include 路径
}

func (i *IncludeRole) After() {
	stoptime := time.Now()
	i.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(i.starttime))
	i.msg.CallBack[fmt.Sprintf("%s-%s-%s", i.host, i.stage, i.name)] = i.logs
}
