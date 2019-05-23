package api

// package main

import (
	"fmt"

	_ "github.com/lflxp/showme/api/routers"
	"github.com/lflxp/showme/utils"

	"github.com/astaxie/beego"
)

func Api(isSwagger bool, ip, port string) {
	utils.InitSqlite()
	if isSwagger {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run(fmt.Sprintf("%s:%s", ip, port))
}

// func main() {
// 	if beego.BConfig.RunMode == "dev" {
// 		beego.BConfig.WebConfig.DirectoryIndex = true
// 		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
// 	}
// 	beego.Run()
// }
