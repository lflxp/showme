package executors

import (
	"github.com/lflxp/showme/executors/dashboard"
	"github.com/lflxp/showme/executors/dashboard2"
)

func ParseExecutors(in string) func() {
	var result func()
	if in == "dashboard show" {
		result = func() {
			dashboard.Run()
		}
	} else if in == "dashboard tcell" {
		result = func() {
			dashboard2.Run()
		}
	}
	return result
}
