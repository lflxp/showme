//go:build !gopacket
// +build !gopacket

package executors

import (
	"fmt"
	"strings"

	// kubectl "github.com/lflxp/lflxp-kubectl/pkg"
	monitor "github.com/lflxp/lflxp-monitor/pkg"
	mysql "github.com/lflxp/lflxp-orzdba/pkg"
	scan "github.com/lflxp/lflxp-scan/pkg"
	"github.com/lflxp/showme/pkg/prompt/completers"
	"github.com/lflxp/showme/pkg/prompt/executors/dashboard"
	"github.com/lflxp/showme/pkg/prompt/executors/helloworld"
	"github.com/lflxp/showme/pkg/prompt/executors/layout"
	"github.com/lflxp/showme/utils"
)

/** 解析执行命令函数
else if in == "kubectl" {
		result = func() {
			kubectl.ManualInit()
		}
		status = true
	}
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
	} else if strings.Contains(in, "mysql") {
		result = func() {
			err := mysql.BeforeRun(in)
			if err != nil {
				fmt.Println(err)
			}
		}
		status = true
	} else if strings.Contains(in, "tty") {
		result = func() {
			err := utils.CommandPty(strings.Replace(in, "tty", "showme tty", -1))
			if err != nil {
				fmt.Println(err)
			}
		}
		status = true
	} else if strings.Contains(in, "sw") {
		result = func() {
			err := utils.CommandPty(strings.Replace(in, "sw", "showme watch", -1))
			if err != nil {
				fmt.Println(err)
			}
		}
		status = true
	} else {
		// fmt.Println(utils.Colorize(in, "red", "black", true, true), "未发现该命令")
		result = func() {
			err := utils.CommandPty(in)
			if err != nil {
				fmt.Println(err)
			}
		}
		status = true
	}
	return result, status
}
