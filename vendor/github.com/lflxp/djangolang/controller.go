package djangolang

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/lflxp/djangolang/components"
	"github.com/lflxp/djangolang/consted"
	"github.com/lflxp/djangolang/middlewares"
	raw "github.com/lflxp/djangolang/model"
	model "github.com/lflxp/djangolang/model/admin"
	"github.com/lflxp/djangolang/utils"

	"github.com/lflxp/djangolang/utils/orm/sqlite"

	"github.com/gin-gonic/gin"
)

func init() {
	// vpn := Vpn{}s
	RegisterAdmin(new(model.Vpn), new(model.Machine), new(model.Cdn), new(model.More), new(model.User), new(model.Claims), new(model.Groups), new(model.Userauth), new(model.History))

	user := model.User{Username: "admin"}
	has, err := sqlite.NewOrm().Get(&user)
	if err != nil {
		slog.Error(err.Error())
	}

	if !has {
		claims := model.Claims{
			Auth:  "admin",
			Type:  "nav",
			Value: "dashboard",
		}

		sqlite.NewOrm().Insert(&claims)

		sql := "insert into user('username','password') values ('admin','admin');"
		n, err := sqlite.NewOrm().Query(sql)
		if err != nil {
			slog.Error("init admin user", "err", err.Error(), "sql", sql)
		}
		slog.Debug("insert admin user", "count", len(n), "sql", sql)
	}
}

// @param router => *gin.Engine
// @param isJwt => bool 是否开启jwt认证模式
func RegisterControllerAdmin(router *gin.Engine, isJwt bool) {
	// 注册模版
	RegisterTemplate(router)
	// 默认跳转
	router.NoRoute(middlewares.NoRouteHandler)
	// form upload serve
	router.StaticFS("/upload", http.Dir("upload"))
	// admin路由
	adminGroup := router.Group("/admin")
	if isJwt {
		// 默认开启jwt认证
		middlewares.RegisterJWT(router)
		var authMiddleware = middlewares.NewGinJwtMiddlewares(middlewares.AllUserAuthorizator)
		adminGroup.Use(authMiddleware.MiddlewareFunc())
	}
	// adminGroup.Use(authMiddleware.MiddlewareFunc())
	{
		adminGroup.Any("/:type", process)
	}
}

// @Summary admin动态渲染接口
// @Description 获取type，逐级数据查询
// @Tags admin
// @Param type path string true "动态类型，有index/table/add/edit/data"
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /admin/:type [get|post]
func process(c *gin.Context) {
	typed := c.Params.ByName("type")
	if c.Request.Method == "GET" {
		if typed == "index" {
			Hs := make([]model.History, 0)
			err := sqlite.NewOrm().Desc("id").Limit(10, 0).Find(&Hs)
			if err != nil {
				c.JSON(400, err)
			} else {
				c.HTML(200, "admin/home.html", gin.H{
					"Data":    GetRegistered(),
					"History": Hs,
					"User":    "Boss",
					"Nav":     true,
				})
			}
		} else if typed == "userinfo" {
			c.JSONP(200, gin.H{
				"username": "admin",
				"userid":   "001",
				"email":    "admin@example.com",
				"data": map[string]string{
					"username": "admin",
					"userid":   "001",
					"email":    "admin@example.com",
					"avatar":   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
					"name":     "Super Admin",
				},
			})
		} else if typed == "tablelist" {
			tmp := []map[string]interface{}{}

			status := []string{"published", "draft", "deleted"}
			for i := 0; i < 30; i++ {
				current := i % 3
				tmp = append(tmp, map[string]interface{}{
					"id":           i,
					"title":        utils.GetRandomString(i + 1),
					"status":       status[current],
					"author":       utils.GetRandomString(10),
					"display_time": time.Now().String(),
					"pageviews":    i,
				})
			}
			c.JSONP(200, gin.H{
				"data": map[string]interface{}{
					"total": 30,
					"items": tmp,
				},
			})
		} else if typed == "history" {
			c.JSON(200, gin.H{
				"json": "ok",
			})
		} else if typed == "test" {
			c.HTML(200, "test.html", gin.H{
				"User": "Boss",
				"Nav":  true,
			})
		} else if typed == "add" {
			data := map[string]interface{}{}
			name := c.Query("name")
			slog.Debug("add name ", name)
			if name != "" {
				data["Col"] = GetRegisterByName(name)
				// data["File"] = false
				columns := []raw.Upload{}
				for _, value := range strings.Split(strings.TrimSpace(GetRegisterByName(name)["Col"]), " ") {
					// slog.Debugf("index %d value %s", index, value)
					tmp := strings.Split(strings.TrimSpace(value), ":")
					// if tmp[1] == "file" {
					// 	data["File"] = true
					// }
					tmp_col := raw.Upload{
						VerboseName: tmp[0],
						ColType:     tmp[1],
						Name:        tmp[2],
						Detail:      tmp[3],
					}
					columns = append(columns, tmp_col)
				}
				data["Columns"] = columns
				// slog.Debugf("Colums is %v", columns)
			}
			data["Up"] = components.NewComponentAdaptor(consted.Upload).Transfer().(func(map[string]string) string)(GetRegisterByName(name))
			data["Data"] = GetRegistered()
			data["Name"] = name
			data["User"] = "Boss"
			data["Nav"] = true
			c.HTML(200, "admin/addmenu.html", data)
		} else if typed == "test1" {
			// fmt.Println(models.Registered)
			Hs := make([]model.History, 0)
			err := sqlite.NewOrm().Desc("id").Limit(10, 0).Find(&Hs)
			if err != nil {
				panic(err)
			}
			Data := map[string]interface{}{}
			Data["Data"] = GetRegistered()
			Data["History"] = Hs
			Data["User"] = "Boss"
			Data["Nav"] = true
			c.HTML(200, "admin/home.html", Data)
		} else if typed == "table" {
			// data, err := utils.DirectJson("First", "Second", "Three", "Four", "Op", "Datetime")
			// if err == nil {
			// 	this.Data["Col"] = data
			// }
			Data := map[string]interface{}{}
			Name := c.Query("name")
			if Name != "" {
				Data["Col"] = GetRegisterByName(Name)
			}
			Data["Data"] = GetRegistered()
			Data["Name"] = Name
			Data["User"] = "Boss"
			Data["Nav"] = true
			c.HTML(200, "admin/tablemenu.html", Data)
		} else if typed == "edit" {
			name := c.Query("name")
			id := c.Query("id")
			Data := map[string]interface{}{}
			if name != "" && id != "" {
				//查询sql
				check_sql := fmt.Sprintf("select * from %s where id=%s", name, id)
				result, err := sqlite.NewOrm().Query(check_sql)
				if err != nil {
					c.String(400, err.Error())
					return
				}
				if len(result) == 1 {
					// TODO: 返回 map[string]string{"table": "tableName", "id": "123"} | editform
					// Data["Col"] = template.EditFormColumns(template.GetRegisterByName(name), result[0])
					Data["Col"] = components.NewComponentAdaptor(consted.EditFormComponent).Transfer().(func(data map[string]string, id, name string) string)(GetRegisterByName(name), id, name)
				} else {
					c.String(400, fmt.Sprintf("Id %s 返回数据超过1条 实际为 %d", id, len(result)))
					return
				}
				Data["Name"] = name

				columns := []raw.Upload{}
				for index, value := range strings.Split(strings.TrimSpace(GetRegisterByName(name)["Col"]), " ") {
					slog.Debug(fmt.Sprintf("index %d value %s", index, value))
					tmp := strings.Split(strings.TrimSpace(value), ":")
					// if tmp[1] == "file" {
					// 	data["File"] = true
					// }
					tmp_col := raw.Upload{
						VerboseName: tmp[0],
						ColType:     tmp[1],
						Name:        tmp[2],
						Detail:      tmp[3],
					}
					columns = append(columns, tmp_col)
				}
				Data["Columns"] = columns
				slog.Debug(fmt.Sprintf("Colums is %v", columns))
			}
			Data["Data"] = GetRegistered()
			Data["User"] = "Boss"
			Data["Nav"] = true
			Data["Id"] = id
			c.HTML(200, "admin/editmenu.html", Data)
		} else if typed == "data" {
			var sql string
			name := c.Query("name")
			order := c.Query("order")
			search := c.Query("search")
			offset := c.Query("offset")
			limit := c.Query("limit")

			searchs := GetRegisterByName(name)
			// 获取字段类型
			columnsType := map[string]string{}
			for _, v := range strings.Split(strings.TrimSpace(searchs["Col"]), " ") {
				// slog.Infof("vvvv is %s", v)
				tmp := strings.Split(strings.TrimSpace(v), ":")
				columnsType[strings.ToLower(tmp[2])] = tmp[1]
			}

			if search == "" {
				sql = fmt.Sprintf("select * from %s order by id %s", strings.ToLower(name), order)
			} else {
				slog.Debug("searchs", searchs)
				if searchs != nil {
					if strings.Contains(searchs["Search"], ",") {
						tmp := []string{}
						for _, v := range strings.Split(searchs["Search"], ",") {
							tmp = append(tmp, fmt.Sprintf("%s like '%%%s%%'", v, search))
						}
						sql = fmt.Sprintf("select * from %s where %s order by id %s", strings.ToLower(name), strings.Join(tmp, " or "), order)
					} else {
						sql = fmt.Sprintf("select * from %s where %s order by id %s", strings.ToLower(name), fmt.Sprintf("id = '%s'", search), order)
					}
				}
			}

			if offset != "" && limit != "" {
				sql = fmt.Sprintf("%s limit %s offset %s", sql, limit, offset)
			}
			slog.Debug(sql)
			result, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				slog.Error(err.Error())
			}
			total, err := sqlite.NewOrm().Table(strings.ToLower(name)).Count()
			if err != nil {
				slog.Error(err.Error())
			}
			ttt := map[string]interface{}{}
			t2 := []map[string]string{}
			for _, x := range result {
				tmp := map[string]string{}
				for key, value := range x {
					// result[n][key] = string(value)
					// 判断字段类型
					// slog.Debugf("=== 字段类型 %v key %s value %s", columnsType, key, value)
					if types, ok := columnsType[strings.ToLower(key)]; ok {
						if types == "file" {
							// slog.Debugf("====== FILE === data key %s value %s", key, value)
							if len(value) > 5 {
								tmp[strings.ToLower(key)] = fmt.Sprintf(`<a href="/%s" target="_self" download>下载</a>`, string(value))
							} else {
								tmp[strings.ToLower(key)] = "-"
							}
						} else {
							// slog.Debugf("====== NOT FILE === data key %s value %s", key, value)
							tmp[strings.ToLower(key)] = string(value)
						}
					}
				}
				op := `<button type="button" id="btnEdit" data-toggle="modal" data-target="#editModel" class="btn btn-warning btn-xs" aria-label="Left Align" onclick="Edit('$NAME','$ID')">
				<span class="glyphicon glyphicon-edit" aria-hidden="true"></span>
			  </button>
				<button type="button" class="btn btn-danger btn-xs" aria-label="Left Align" onclick="Delete('$NAME','$ID')">
				<span class="glyphicon glyphicon-remove" aria-hidden="true"></span>
				</button>`
				tmp["操作"] = strings.Replace(strings.Replace(op, "$NAME", name, -1), "$ID", string(x["id"]), -1)

				t2 = append(t2, tmp)
			}
			ttt["rows"] = t2
			ttt["total"] = total
			c.JSONP(200, ttt)
		}
	} else if c.Request.Method == "POST" {
		if typed == "check" {
			c.JSONP(200, gin.H{
				"json": "hello",
			})
		} else if typed == "delete" {
			ids := c.Query("ids")
			name := c.Query("name")
			history := model.History{
				Name: name,
				Op:   fmt.Sprintf("Delete Id %s", ids),
			}
			slog.Debug(ids, name)
			sql := fmt.Sprintf("delete from %s where id in (%s)", name, ids)
			_, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				c.String(400, err.Error())
				history.Common = err.Error()
				return
			}
			history.Common = "success"
			defer model.AddHistory(&history)
			c.String(200, fmt.Sprintf("delete %s %s success", name, ids))
		} else if typed == "add" {
			// col := []string{}
			// value := []string{}

			name := c.Query("table")
			coldata := GetRegisterByName(name)
			// slog.Debugf("================= coldata is %v", coldata)
			colList := []string{}
			valueList := []string{}
			for _, x := range strings.Split(strings.TrimSpace(coldata["Col"]), " ") {
				tmp := strings.Split(strings.TrimSpace(x), ":")
				if tmp[2] == "id" {
					continue
				} else if tmp[1] == "file" {
					file, err := c.FormFile(tmp[2])
					if err != nil {
						slog.Error(err.Error())
						c.String(400, "file upload error %v", err.Error())
						return
					}
					fmt.Println("接收的数据", name, file.Filename)
					//获取文件名称
					fmt.Println(file.Filename)
					//文件大小
					fmt.Println(file.Size)
					// //获取文件的后缀名
					// extstring := path.Ext(file.Filename)
					// fmt.Println(extstring)
					//根据当前时间鹾生成一个新的文件名
					fileNameInt := time.Now().Unix()
					fileNameStr := strconv.FormatInt(fileNameInt, 10)
					//新的文件名
					fileName := fmt.Sprintf("%s_%s", file.Filename, fileNameStr)
					//保存上传文件
					filePath := filepath.Join(utils.Mkdir("upload"), "/", fileName)
					slog.Debug(fmt.Sprintf("=================上传文件 %s size %d name %s", filePath, file.Size, file.Filename, fileName))
					c.SaveUploadedFile(file, filePath)
					colList = append(colList, tmp[2])
					valueList = append(valueList, fmt.Sprintf("'%s'", filePath))
				} else {
					tt := c.PostForm(tmp[2])
					// slog.Debugf("postForm %s is %s", tmp[2], tt)
					colList = append(colList, tmp[2])
					valueList = append(valueList, fmt.Sprintf("'%s'", tt))
				}
			}

			// slog.Error(string(this.Ctx.Input.RequestBody))
			//获取字段和所有值
			// bs, err := io.ReadAll(c.Request.Body)
			// if err != nil {
			// 	c.String(400, err.Error())
			// 	return
			// }
			// slog.Debugf("================= bs is %s", string(bs))
			// body := strings.Replace(string(bs), "&_save=%E4%BF%9D%E5%AD%98", "", -1)
			// body = strings.Replace(body, "%5B%5D", "", -1)
			// parseBody, err := url.Parse(body)
			// if err != nil {
			// 	c.String(400, err.Error())
			// 	return
			// }
			// slog.Error(parseBody.Path)
			// l, err := url.ParseQuery(parseBody.Path)
			// if err != nil {
			// 	slog.Error(err.Error())
			// 	c.String(400, err.Error())
			// 	return
			// }

			// for keyed, valueed := range l {
			// 	col = append(col, fmt.Sprintf("'%s'", keyed))
			// 	value = append(value, fmt.Sprintf("'%s'", strings.Join(valueed, ",")))
			// }

			history := model.History{
				Name: name,
				Op:   "Add",
			}

			sql := fmt.Sprintf("insert into %s(%s) values (%s)", name, strings.Join(colList, ","), strings.Join(valueList, ","))
			slog.Debug(sql)
			_, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				history.Common = err.Error()
				c.String(400, err.Error())
				return
			}

			history.Common = "success"
			defer model.AddHistory(&history)

			// search
			// issearch := viper.GetBool("meili.enable")

			// if issearch {
			// 	tmpData := map[string]interface{}{
			// 		"uid": uuid.New(),
			// 	}
			// 	for keyed, valueed := range l {
			// 		tmpData[keyed] = strings.Join(valueed, ",")
			// 	}

			// 	index := utils.SearchCli.Index(name)
			// 	update, err := index.AddDocuments([]map[string]interface{}{tmpData})
			// 	if err != nil {
			// 		c.String(400, err.Error())
			// 		return
			// 	}
			// 	core.Debugf("add meili docs %d %v", update.UpdateID, tmpData)
			// }

			// c.String("insert ok")
			c.Redirect(301, fmt.Sprintf("/admin/add?name=%s", name))
		} else if typed == "edit" {
			name := c.Query("table")
			id := c.Query("id")
			if name != "" && id != "" {
				history := model.History{
					Name: name,
					Op:   fmt.Sprintf("Edit Id %s", id),
				}
				// bs, err := io.ReadAll(c.Request.Body)
				// if err != nil {
				// 	history.Common = err.Error()
				// 	c.String(400, err.Error())
				// 	return
				// }
				// // core.Debug(string(bs))
				// //获取字段和所有值
				// body := strings.Replace(string(bs), "&_save=%E4%BF%9D%E5%AD%98", "", -1)
				// // 过滤claims_id[] => claims_id
				// body = strings.Replace(body, "%5B%5D", "", -1)
				// // core.Infof("body is %s", body)
				// parseBody, err := url.Parse(body)
				// if err != nil {
				// 	history.Common = err.Error()
				// 	c.String(400, err.Error())
				// 	return
				// }
				// // slog.Error(parseBody.Path)
				// l, err := url.ParseQuery(parseBody.Path)
				// if err != nil {
				// 	slog.Error(err.Error())
				// 	history.Common = err.Error()
				// 	c.String(400, err.Error())
				// 	return
				// }

				coldata := GetRegisterByName(name)
				slog.Debug(fmt.Sprintf("================= coldata is %v", coldata))
				colList := []string{}
				valueList := []string{}
				for _, x := range strings.Split(strings.TrimSpace(coldata["Col"]), " ") {
					tmp := strings.Split(strings.TrimSpace(x), ":")
					if tmp[2] == "id" {
						continue
					} else if tmp[1] == "file" {
						file, err := c.FormFile(tmp[2])
						if err != nil {
							slog.Error(err.Error())
							c.String(400, "file upload error %v", err.Error())
							return
						}
						fmt.Println("接收的数据", name, file.Filename)
						//获取文件名称
						fmt.Println(file.Filename)
						//文件大小
						fmt.Println(file.Size)
						// //获取文件的后缀名
						// extstring := path.Ext(file.Filename)
						// fmt.Println(extstring)
						//根据当前时间鹾生成一个新的文件名
						fileNameInt := time.Now().Unix()
						fileNameStr := strconv.FormatInt(fileNameInt, 10)
						//新的文件名
						fileName := fmt.Sprintf("%s_%s", file.Filename, fileNameStr)
						//保存上传文件
						filePath := filepath.Join(utils.Mkdir("upload"), "/", fileName)
						slog.Debug(fmt.Sprintf("=================上传文件 %s size %d name %s", filePath, file.Size, file.Filename, fileName))
						c.SaveUploadedFile(file, filePath)
						colList = append(colList, tmp[2])
						valueList = append(valueList, filePath)
					} else {
						tt := c.PostForm(tmp[2])
						slog.Debug(fmt.Sprintf("postForm %s is %s", tmp[2], tt))
						colList = append(colList, tmp[2])
						valueList = append(valueList, tt)
					}
				}

				set_string := []string{}
				for index, value := range colList {
					if value != "id" {
						set_string = append(set_string, fmt.Sprintf("'%s'='%s'", value, valueList[index]))
					}
				}
				sql := fmt.Sprintf("update %s set %s where id=%s", name, strings.Join(set_string, ","), id)
				// core.Debug(sql)
				_, err := sqlite.NewOrm().Query(sql)
				if err != nil {
					history.Common = err.Error()
					slog.Error(fmt.Sprintf("sql %s error %s", sql, err.Error()))
					c.String(400, err.Error())
					return
				}

				history.Common = "success"
				defer model.AddHistory(&history)

				// search
				// issearch := viper.GetBool("meili.enable")

				// if issearch {
				// 	tmpData := map[string]interface{}{
				// 		"uid": uuid.New(),
				// 	}
				// 	for keyed, valueed := range l {
				// 		tmpData[keyed] = strings.Join(valueed, ",")
				// 	}

				// 	index := utils.SearchCli.Index(name)
				// 	update, err := index.AddDocuments([]map[string]interface{}{tmpData})
				// 	if err != nil {
				// 		c.String(400, err.Error())
				// 		return
				// 	}
				// 	core.Debugf("add meili docs %d %v", update.UpdateID, tmpData)
				// }
			}
			// c.String("insert ok")
			c.Redirect(301, fmt.Sprintf("/admin/add?name=%s", name))
		}
	}

}

// @Summary  查询指定key的值
// @Description 获取value，逐级数据查询
// @Tags admin
// @Param token query string false "token"
// @Param key query string true "查询key"
// @Success 200 {string} string "success"
// @Security ApiKeyAuth
// @Router /api/v1/admin/get [get]
func getadmin(c *gin.Context) {
	key, isok := c.GetQuery("key")
	if !isok {
		c.JSONP(200, gin.H{
			"status": false,
			"msg":    "key is ",
		})
	}

	admin := new(model.Vpn)

	c.JSONP(200, gin.H{
		"status": true,
		"data":   admin,
		"source": key,
	})
}
