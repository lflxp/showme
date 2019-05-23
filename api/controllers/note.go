package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/lflxp/showme/api/models"
)

// Operations about Users
type NoteController struct {
	BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.Group	true		"body for user content"
// @Success 200 {int} models.Group.Id
// @Failure 403 body is empty
// @router /group [post]
func (u *NoteController) PostGroup() {
	var user models.Group
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	num, err := user.Insert()
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("success insert %d ", num)
	}
	u.ServeJSON()
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.Note	true		"body for user content"
// @Success 200 {int} models.Note.Id
// @Failure 403 body is empty
// @router /note [post]
func (u *NoteController) PostNote() {
	var user models.Note
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	num, err := user.Insert()
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("success insert %d ", num)
	}
	u.ServeJSON()
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.Common	true		"body for user content"
// @Success 200 {int} models.Common.Id
// @Failure 403 body is empty
// @router /common [post]
func (u *NoteController) PostCommon() {
	var user models.Common
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	num, err := user.Insert()
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("success insert %d ", num)
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.Group
// @router /group [get]
func (u *NoteController) GetGroupAll() {
	data, err := models.GetAllGroups()
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = data
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	gname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Group
// @router /group/:gname [get]
func (u *NoteController) GetGroupAllSecond() {
	gname := u.GetString(":gname")
	if gname != "" {
		data, err := models.GetAllGroupsSecond(gname)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = data
		}
	}
	u.ServeJSON()
}

// @Title Distinct
// @Description distinct all Users
// @Param	gname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Note
// @router /distinct/:gname [get]
func (u *NoteController) DistinctNote() {
	gname := u.GetString(":gname")
	if gname != "" {
		data, err := models.DistinctNote(gname)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = data
		}
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	gname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Pay
// @router /note/:gname [get]
func (u *NoteController) GetNoteAll() {
	gname := u.GetString(":gname")
	if gname != "" {
		data, err := models.GetAllNotes(gname)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = data
		}
	}

	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Param	nname		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Pay
// @router /common/:nname [get]
func (u *NoteController) GetCommonAll() {
	nname := u.GetString(":nname")
	if nname != "" {
		data, err := models.GetAllCommon(nname)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = data
		}
	}

	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.Group	true		"body for user content"
// @Success 200 {object} models.Group
// @Failure 403 :uid is not int
// @router /group/:uid [put]
func (u *NoteController) PutGroup() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.Group
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := user.Update(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = fmt.Sprintf("update success %d", uu)
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.Note	true		"body for user content"
// @Success 200 {object} models.Note
// @Failure 403 :uid is not int
// @router /note/:uid [put]
func (u *NoteController) PutNote() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.Note
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := user.Update(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = fmt.Sprintf("update success %d", uu)
		}
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.Common	true		"body for user content"
// @Success 200 {object} models.Common
// @Failure 403 :uid is not int
// @router /common/:uid [put]
func (u *NoteController) PutCommon() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.Common
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := user.Update(uid)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = fmt.Sprintf("update success %d", uu)
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /group/:uid [delete]
func (u *NoteController) DeleteGroup() {
	uid := u.GetString(":uid")
	n, err := models.DeleteGroup(uid)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("delete success %d", n)
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /note/:uid [delete]
func (u *NoteController) DeleteNote() {
	uid := u.GetString(":uid")
	n, err := models.DeleteNote(uid)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("delete success %d", n)
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /common/:uid [delete]
func (u *NoteController) DeleteCommon() {
	uid := u.GetString(":uid")
	n, err := models.DeleteCommon(uid)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("delete success %d", n)
	}
	u.ServeJSON()
}
