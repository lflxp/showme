package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PostCommon",
            Router: `/common`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "GetCommonAll",
            Router: `/common/:nname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PutCommon",
            Router: `/common/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "DeleteCommon",
            Router: `/common/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "DistinctNote",
            Router: `/distinct/:gname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "GetGroupAll",
            Router: `/group`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PostGroup",
            Router: `/group`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "GetGroupAllSecond",
            Router: `/group/:gname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PutGroup",
            Router: `/group/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "DeleteGroup",
            Router: `/group/:uid`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PostNote",
            Router: `/note`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "GetNoteAll",
            Router: `/note/:gname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "PutNote",
            Router: `/note/:uid`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"] = append(beego.GlobalControllerRouter["github.com/lflxp/showme/api/controllers:NoteController"],
        beego.ControllerComments{
            Method: "DeleteNote",
            Router: `/note/:uid`,
            AllowHTTPMethods: []string{"delete"},
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
