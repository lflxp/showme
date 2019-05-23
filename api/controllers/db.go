package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"

	"github.com/lflxp/showme/api/models"
	"github.com/lflxp/showme/utils"
)

// Operations about Users
type DbController struct {
	BaseController
}

// @Title Select
// @Description create users
// @Param	body		body 	models.Db	true		"body for user content"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /select [post]
func (u *DbController) Select() {
	var db models.Db
	json.Unmarshal(u.Ctx.Input.RequestBody, &db)
	rs, err := utils.Engine.QueryString(db.Sql)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = rs
	}

	u.ServeJSON()
}

// @Title Insert
// @Description create users
// @Param	body		body 	models.Db	true		"body for user content"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /exec [post]
func (u *DbController) Exec() {
	var db models.Db
	json.Unmarshal(u.Ctx.Input.RequestBody, &db)
	rs, err := utils.Engine.Exec(db.Sql)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = rs
	}
	u.ServeJSON()
}

// @Title Shell
// @Description create users
// @Param	body		body 	models.Db	true		"body for user content"
// @Success 200 {string} success
// @Failure 403 body is empty
// @router /shell [post]
func (u *DbController) Shell() {
	var db models.Db
	json.Unmarshal(u.Ctx.Input.RequestBody, &db)
	rs, err := utils.ExecCommand(db.Sql)
	beego.Critical("111", string(rs), db.Sql, string(u.Ctx.Input.RequestBody))
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = string(rs)
	}
	u.ServeJSON()
}
