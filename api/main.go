package api

import (
	"fmt"

	_ "github.com/lflxp/showme/api/routers"

	"github.com/astaxie/beego"
	"github.com/lflxp/showme/utils"
)

func Api(isSwagger bool, ip, port, dbName, defaultDb string) {
	utils.InitDB(dbName, defaultDb)
	if isSwagger {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run(fmt.Sprintf("%s:%s", ip, port))
}

// func main() {
// 	utils.InitDB("dbName", "defaultDb")
// 	if beego.BConfig.RunMode == "dev" {
// 		beego.BConfig.WebConfig.DirectoryIndex = true
// 		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
// 	}
// 	beego.Run()
// }
