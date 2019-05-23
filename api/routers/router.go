// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/lflxp/showme/api/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
			beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
			beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		),
		beego.NSNamespace("/db",
			beego.NSInclude(
				&controllers.DbController{},
			),
			beego.NSRouter("/*", &controllers.BaseController{}, "options:Options"),
		),
	)
	beego.AddNamespace(ns)
}
