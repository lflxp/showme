package executors

import (
	"github.com/lflxp/showme/executors/dashboard"
	"github.com/lflxp/showme/executors/helloworld"
	"github.com/lflxp/showme/executors/layout"
)

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
	}
	return result, status
}
