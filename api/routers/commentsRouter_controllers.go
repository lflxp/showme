package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "GetDataByTablename",
            Router: `/:dbname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:dbname/:tablename`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:tablename/:key`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/kv`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "PostInfo",
            Router: `/kv/info`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "AddTables",
            Router: `/tables`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "GetAllTables",
            Router: `/tables/:dbname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:BboltController"],
        beego.ControllerComments{
            Method: "GetAllTables2",
            Router: `/tables2/:dbname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:ObjectController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:objectId`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Post",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Put",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: `/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Login",
            Router: `/login`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:UserController"],
        beego.ControllerComments{
            Method: "Logout",
            Router: `/logout`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
