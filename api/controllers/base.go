package controllers

/*
	1.允许Headers表头通过验证 YC-Token
	2.增加Options 功能验证
	3.允许服务器跨域访问
*/
import (
	"github.com/astaxie/beego"
)

type Other struct {
	User string
	Pwd  string
	Name string
}

type BaseController struct {
	beego.Controller
	Others Other
}

func (this *BaseController) Prepare() {
	beego.ReadFromRequest(&this.Controller)
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")                                     //允许访问源
	this.Ctx.Output.Header("Access-Control-Allow-Methods", "DELETE,POST, GET, PUT, OPTIONS")       //允许post访问
	this.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,LXP-Token") //header的类型
	this.Ctx.Output.Header("Access-Control-Max-Age", "1728000")
	this.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	this.Ctx.Output.Header("content-type", "application/json") //返回数据格式是json
}

func (this *BaseController) AllowCross() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")                                     //允许访问源
	this.Ctx.Output.Header("Access-Control-Allow-Methods", "DELETE,POST, GET, PUT, OPTIONS")       //允许post访问
	this.Ctx.Output.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,LXP-Token") //header的类型
	this.Ctx.Output.Header("Access-Control-Max-Age", "1728000")
	this.Ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	this.Ctx.Output.Header("content-type", "application/json") //返回数据格式是json
}

func (this *BaseController) Options() {
	this.AllowCross() //允许跨域
	this.Data["json"] = map[string]interface{}{"status": 200, "message": "ok", "moreinfo": ""}
	this.ServeJSON()
}
