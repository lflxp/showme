package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/lflxp/showme/api/models"
	pkg "github.com/lflxp/showme/utils"
)

type BboltController struct {
	BaseController
}

// @Title CreateBucket
// @Description 添加一个完整的数据 表、key、value
// @Param	body		body 	models.Bucket	true		"body for Bucket content"
// @Success 200 {int} models.Bucket.Key
// @Failure 403 body is empty
// @router /kv [post]
func (u *BboltController) Post() {
	var user models.Bucket
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	beego.Critical(user)
	err := pkg.AddKeyValueByBucketName(user.Dbname, user.Tablename, user.Key, user.Value)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("add %s %s success", user.Tablename, user.Key)
	}
	u.ServeJSON()
}

// @Title CreateBucketInfo
// @Description 添加一个针对Tree记录的完整的数据 表、key、value
// @Param	body		body 	models.Bucket	true		"body for Bucket content"
// @Success 200 {int} models.Bucket.Key
// @Failure 403 body is empty
// @router /kv/info [post]
func (u *BboltController) PostInfo() {
	var user models.BucketInfo
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	beego.Critical(user)
	err := pkg.AddKeyValueByBucketName(user.Dbname, user.Tablename, user.Key, string(u.Ctx.Input.RequestBody))
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = fmt.Sprintf("add %s %s success", user.Tablename, user.Key)
	}
	u.ServeJSON()
}

// @Title CreateUser
// @Description 只添加表
// @Param	body		body 	models.Bucket	true		"Tablename字段中只有tablename字段会被用到"
// @Success 200 {int} models.Tablename
// @Failure 403 body is empty
// @router /tables [post]
func (u *BboltController) AddTables() {
	var user models.Bucket
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	fmt.Println(user)
	err := pkg.CreateBucket(user.Tablename)
	if err != nil {
		u.Data["json"] = err
	} else {
		u.Data["json"] = fmt.Sprintf("add table %s success", user.Tablename)
	}
	u.ServeJSON()
}

// @Title GetAllTables
// @Description 获取所有表的数据
// @Param  dbname  path  int true  "name of db"
// @Success 200 {string} success
// @router /tables/:dbname [get]
func (u *BboltController) GetAllTables() {
	dbname := u.GetString(":dbname")
	data, err := pkg.GetAllTables(dbname)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = data
	}
	u.ServeJSON()
}

// @Title GetAllTables
// @Description 获取所有表的数据
// @Param  dbname  path  int true  "name of db"
// @Success 200 {string} success
// @router /tables2/:dbname [get]
func (u *BboltController) GetAllTables2() {
	dbname := u.GetString(":dbname")
	data, err := pkg.GetAllTables(dbname)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = data
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description 遍历指定表的所有数据
// @Param  dbname  path  string true  "name of table"
// @Success 200 {string} success
// @router /:dbname [get]
func (u *BboltController) GetDataByTablename() {
	dbname := u.GetString(":dbname")
	if dbname != "" {
		data, err := pkg.GetAllByTables(dbname)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = data
		}
	} else {
		u.Data["json"] = "tablename or key is none"
	}
	u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *BboltController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJSON()
}

// @Title Get
// @Description get user by uid
// @Param	tablename path 	string	true		"The key for staticblock"
// @Param	key		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Bucket
// @Failure 403 :tablename is empty
// @router /:tablename/:key [get]
func (u *BboltController) Get() {
	tablename := u.GetString(":tablename")
	key := u.GetString(":key")
	if tablename != "" && key != "" {
		data, err := pkg.GetValueByBucketName(tablename, key)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = string(data)
		}
	} else {
		u.Data["json"] = "tablename or key is none"
	}
	u.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *BboltController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err.Error()
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJSON()
}

// @Title Delete
// @Description delete the Table
// @Param  dbname  path  int true  "name of db"
// @Param	tablename		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 tablename is empty
// @router /:dbname/:tablename [delete]
func (u *BboltController) Delete() {
	dbname := u.GetString(":dbname")
	tablename := u.GetString(":tablename")
	err := pkg.DeleteBucket(dbname, tablename)
	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		u.Data["json"] = "delete success!"
	}
	u.ServeJSON()
}
