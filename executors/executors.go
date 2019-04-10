package executors

import (
	"fmt"

	"github.com/lflxp/showme/executors/dashboard"
	"github.com/lflxp/showme/executors/helloworld"
	"github.com/lflxp/showme/executors/layout"
	"github.com/lflxp/showme/executors/monitor"
)

/** 解析执行命令函数
@param in   // command from
@result func() // function
@result bool // 状态 是否执行
*/
func ParseExecutors(in string) (func(), bool) {
	var result func()
	status := false
	if in == "dashboard show" {
		result = func() {
			dashboard.Run()
		}
		status = true
	} else if in == "dashboard helloworld" {
		result = func() {
			helloworld.Run()
		}
		status = true
	} else if in == "gocui active" {
		result = func() {
			layout.Run()
		}
		status = true
	} else if in == "dashboard" {
		result = func() {
			dashboard.Dashboard()
		}
		status = true
	} else if in == "monitor -lazy" {
		result = func() {
			monitor.Run()
		}
		status = true
	} else {
		fmt.Println(monitor.Colorize(in, "red", "black", true, true), " not found executors")
	}
	return result, status
}
