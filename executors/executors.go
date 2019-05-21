// +build !gopacket

package executors

import (
	"fmt"
	"strings"

	"github.com/lflxp/showme/completers"
	"github.com/lflxp/showme/executors/dashboard"
	"github.com/lflxp/showme/executors/helloworld"
	"github.com/lflxp/showme/executors/httpstatic"
	"github.com/lflxp/showme/executors/layout"
	"github.com/lflxp/showme/executors/monitor"
	"github.com/lflxp/showme/executors/mysql"
	"github.com/lflxp/showme/executors/scan"
	"github.com/lflxp/showme/utils"
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
	} else if strings.Contains(in, "monitor") {
		result = func() {
			monitor.Run(in)
		}
		status = true
	} else if in == "help" {
		result = func() {
			completers.Help()
		}
		status = true
	} else if strings.Contains(in, "scan") {
		result = func() {
			scan.Scan(in)
		}
		status = true
	} else if strings.Contains(in, "httpstatic") {
		tmp := strings.Split(in, " ")
		if len(tmp)%2 == 1 {
			result = func() {
				httpstatic.HttpStaticServe(in)
			}
			status = true
		} else {
			fmt.Println(utils.Colorize("输入错误，请输入完整的【-port】或【-path】", "red", "black", true, true))
			status = false
		}
	} else if strings.Contains(in, "mysql") {
		result = func() {
			err := mysql.BeforeRun(in)
			if err != nil {
				fmt.Println(err)
			}
		}
		status = true
	} else {
		fmt.Println(utils.Colorize(in, "red", "black", true, true), " not found executors")
	}
	return result, status
}
