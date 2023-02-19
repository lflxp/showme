package controller

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/lflxp/lflxp-music/core/model/admin"
	"github.com/lflxp/lflxp-music/core/utils"

	"github.com/lflxp/lflxp-music/core/middlewares/jwt/framework"
	"github.com/lflxp/lflxp-music/core/middlewares/jwt/services"
	"github.com/lflxp/lflxp-music/core/middlewares/template"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/tools/orm/sqlite"
)

func RegisterAdmin(router *gin.Engine) {
	noauthGroup := router.Group("")
	noauthGroup.GET("/login", login)
	authMiddleware := framework.NewGinJwtMiddlewares(services.AllUserAuthorizator)
	adminGroup := router.Group("/admin")
	adminGroup.Use(authMiddleware.MiddlewareFunc())
	{
		adminGroup.Any("/:type", process)
	}

}

func login(c *gin.Context) {
	c.HTML(200, "admin/login.html", gin.H{
		"User": "Boss",
	})
}

func process(c *gin.Context) {
	typed := c.Params.ByName("type")
	if c.Request.Method == "GET" {
		if typed == "index" {
			Hs := make([]admin.History, 0)
			err := sqlite.NewOrm().Desc("id").Limit(10, 0).Find(&Hs)
			if err != nil {
				c.JSON(400, err)
			} else {
				c.HTML(200, "admin/home.html", gin.H{
					"Data":    template.GetRegistered(),
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
			if name != "" {
				data["Col"] = template.GetRegisterByName(name)
			}
			data["Name"] = name
			data["User"] = "Boss"
			data["Nav"] = false
			c.HTML(200, "admin/add.html", data)
		} else if typed == "test1" {
			// fmt.Println(models.Registered)
			Hs := make([]admin.History, 0)
			err := sqlite.NewOrm().Desc("id").Limit(10, 0).Find(&Hs)
			if err != nil {
				panic(err)
			}
			Data := map[string]interface{}{}
			Data["Data"] = template.GetRegistered()
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
				Data["Col"] = template.GetRegisterByName(Name)
			}
			Data["Data"] = template.GetRegistered()
			Data["Name"] = Name
			Data["User"] = "Boss"
			Data["Nav"] = true
			c.HTML(200, "admin/table.html", Data)
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
					Data["Col"] = template.EditFormColumns(template.GetRegisterByName(name), result[0])
				} else {
					c.String(400, fmt.Sprintf("Id %s 返回数据超过1条 实际为 %d", id, len(result)))
					return
				}
				Data["Name"] = name
			}
			Data["User"] = "Boss"
			Data["Nav"] = true
			c.HTML(200, "admin/edit.html", Data)
		} else if typed == "data" {
			var sql string
			name := c.Query("name")
			order := c.Query("order")
			search := c.Query("search")
			offset := c.Query("offset")
			limit := c.Query("limit")
			if search == "" {
				sql = fmt.Sprintf("select * from %s order by id %s limit %s offset %s", strings.ToLower(name), order, limit, offset)
			} else {
				searchs := template.GetRegisterByName(name)
				if searchs != nil {
					if strings.Contains(searchs["Search"], ",") {
						sql = fmt.Sprintf("select * from %s where %s order by id %s limit %s offset %s", strings.ToLower(name), strings.Replace(searchs["Search"], ",", fmt.Sprintf("='%s' and ", search), -1), order, limit, offset)
					} else {
						sql = fmt.Sprintf("select * from %s where %s order by id %s limit %s offset %s", strings.ToLower(name), fmt.Sprintf("%s='%s'", searchs["Search"], search), order, limit, offset)
					}
				}
			}
			log.Error(sql)
			result, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				log.Error(err.Error())
			}
			total, err := sqlite.NewOrm().Table(strings.ToLower(name)).Count()
			if err != nil {
				log.Error(err.Error())
			}
			log.Error(sql)
			ttt := map[string]interface{}{}
			t2 := []map[string]string{}
			for _, x := range result {
				tmp := map[string]string{}
				for key, value := range x {
					// result[n][key] = string(value)
					tmp[strings.ToLower(key)] = string(value)
				}
				op := `<button type="button" class="btn btn-warning" aria-label="Left Align" onclick="Edit('$NAME','$ID')">
				<span class="glyphicon glyphicon-edit" aria-hidden="true"></span>
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
			log.Error(ids, name)
			sql := fmt.Sprintf("delete from %s where id in (%s)", name, ids)
			_, err := sqlite.NewOrm().Query(sql)
			if err != nil {
				c.String(400, err.Error())
				return
			}
			c.String(200, fmt.Sprintf("delete %s %s success", name, ids))
		} else if typed == "add" {
			col := []string{}
			value := []string{}

			name := c.Query("table")
			// log.Error(string(this.Ctx.Input.RequestBody))
			//获取字段和所有值
			bs, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.String(400, err.Error())
				return
			}
			body := strings.Replace(string(bs), "&_save=%E4%BF%9D%E5%AD%98", "", -1)
			body = strings.Replace(body, "%5B%5D", "", -1)
			parseBody, err := url.Parse(body)
			if err != nil {
				c.String(400, err.Error())
				return
			}
			// log.Error(parseBody.Path)
			l, err := url.ParseQuery(parseBody.Path)
			if err != nil {
				log.Error(err.Error())
				c.String(400, err.Error())
				return
			}

			for keyed, valueed := range l {
				col = append(col, fmt.Sprintf("'%s'", keyed))
				value = append(value, fmt.Sprintf("'%s'", strings.Join(valueed, ",")))
			}

			sql := fmt.Sprintf("insert into %s(%s) values (%s)", name, strings.Join(col, ","), strings.Join(value, ","))
			log.Debug(sql)
			_, err = sqlite.NewOrm().Query(sql)
			if err != nil {
				c.String(400, err.Error())
				return
			}

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
				bs, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					c.String(400, err.Error())
					return
				}
				// core.Debug(string(bs))
				//获取字段和所有值
				body := strings.Replace(string(bs), "&_save=%E4%BF%9D%E5%AD%98", "", -1)
				// 过滤claims_id[] => claims_id
				body = strings.Replace(body, "%5B%5D", "", -1)
				// core.Infof("body is %s", body)
				parseBody, err := url.Parse(body)
				if err != nil {
					c.String(400, err.Error())
					return
				}
				// log.Error(parseBody.Path)
				l, err := url.ParseQuery(parseBody.Path)
				if err != nil {
					log.Error(err.Error())
					c.String(400, err.Error())
					return
				}

				set_string := []string{}
				for keyed, valueed := range l {
					if keyed != "id" {
						set_string = append(set_string, fmt.Sprintf("'%s'='%s'", keyed, strings.Join(valueed, ",")))
					}
				}
				sql := fmt.Sprintf("update %s set %s where id=%s", name, strings.Join(set_string, ","), id)
				// core.Debug(sql)
				_, err = sqlite.NewOrm().Query(sql)
				if err != nil {
					log.Error("sql error", err.Error())
					c.String(400, err.Error())
					return
				}

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
// @Success 200 {object} admin.Vpn admin{}
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

	admin := new(admin.Vpn)

	c.JSONP(200, gin.H{
		"status": true,
		"data":   admin,
		"source": key,
	})
}
