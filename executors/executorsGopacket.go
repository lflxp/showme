// +build gopacket

package executors

import (
	"fmt"
	"strings"

	monitor "github.com/lflxp/lflxp-monitor/pkg"
	scan "github.com/lflxp/lflxp-scan/pkg"
	"github.com/lflxp/showme/completers"
	"github.com/lflxp/showme/executors/dashboard"
	"github.com/lflxp/showme/executors/gopacket"
	"github.com/lflxp/showme/executors/helloworld"
	"github.com/lflxp/showme/executors/layout"
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
	} else if strings.Contains(in, "gopacket") {
		if strings.Contains(in, "gopacket in") {
			result = func() {
				gopacket.Run(in)
			}
		}
		if strings.Contains(in, "gopacket screen") {
			result = func() {
				// gopacket.Gopacket(in)
				gopacket.Screen(strings.Split(in, " ")[2])
			}
		}
		status = true
	} else if strings.Contains(in, "scan") {
		result = func() {
			scan.Scan(in)
		}
		status = true
	} else {
		fmt.Println(utils.Colorize(in, "red", "black", true, true), " not found executors")
	}
	return result, status
}
