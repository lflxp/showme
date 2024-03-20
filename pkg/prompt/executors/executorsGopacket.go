//go:build gopacket
// +build gopacket

package executors

import (
	"fmt"
	"strings"

	// kubectl "github.com/lflxp/lflxp-kubectl/pkg"
	monitor "github.com/lflxp/lflxp-monitor/pkg"
	// mysql "github.com/lflxp/lflxp-orzdba/pkg"
	scan "github.com/lflxp/lflxp-scan/pkg"
	"github.com/lflxp/showme/pkg/prompt/completers"
	"github.com/lflxp/showme/pkg/prompt/executors/dashboard"
	"github.com/lflxp/showme/pkg/prompt/executors/gopacket"
	"github.com/lflxp/showme/pkg/prompt/executors/helloworld"
	"github.com/lflxp/showme/pkg/prompt/executors/layout"
	"github.com/lflxp/showme/utils"
)

/*
  - 解析执行命令函数else if in == "kubectl" {
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
		// } else if strings.Contains(in, "mysql") {
		// 	result = func() {
		// 		err := mysql.BeforeRun(in)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 		}
		// 	}
		// 	status = true
	} else {
		fmt.Println(utils.Colorize(in, "red", "black", true, true), " not found executors")
	}
	return result, status
}
