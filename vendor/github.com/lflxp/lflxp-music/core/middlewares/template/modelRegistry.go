package template

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"

	"github.com/lflxp/tools/orm/sqlite"
)

var registered []map[string]string

func GetRegistered() []map[string]string {
	if registered == nil {
		registered = []map[string]string{}
	}
	return registered
}

func GetRegisterByName(name string) map[string]string {
	for _, maps := range GetRegistered() {
		if strings.EqualFold(strings.ToUpper(maps["Struct"]), strings.ToUpper(name)) {
			return maps
		}
	}
	return nil
}

// http://code.lflxp.cn/life/services/src/master/util/db/xormTest.go
// 收集struct信息
func RigsterStruct(data interface{}) map[string]string {
	tmp := map[string]string{}
	// ttt := reflect.TypeOf(data)
	// beego.Critical(ttt.Elem())
	// beego.Critical(ttt.NumMethod())
	// beego.Critical(ttt.MethodByName().Func.)

	vv := reflect.ValueOf(data)
	v := reflect.Indirect(vv)

	// beego.Critical("API", v.Type().Name(), v.Type().PkgPath())
	services := strings.Split(v.Type().PkgPath(), "/")
	tmp["Services"] = services[len(services)-1]
	// vv.MethodByName().FieldByName()
	// vv.MethodByName().IsValid

	tmp["Struct"] = v.Type().Name()
	for i := 0; i < v.NumField(); i++ {
		if strings.ToLower(v.Type().Field(i).Name) != "id" && v.Type().Field(i).Tag.Get("name") != "id" && v.Type().Field(i).Tag.Get("xorm") != "-" && strings.ToLower(v.Type().Field(i).Name) != "create" && strings.ToLower(v.Type().Field(i).Name) != "update" {
			//利用反射获取structTag
			tmp["Tag"] = fmt.Sprintf("%s", v.Type().Field(i).Tag)
			// html 字段类型
			if v.Type().Field(i).Tag.Get("colType") != "" {
				tmp["ColType"] = v.Type().Field(i).Tag.Get("colType")
			} else {
				tmp["ColType"] = v.Type().Field(i).Type.String()
			}
			// beego.Critical("ColType", tmp["ColType"])
			switch tmp["ColType"] {
			case "radio":
				tmp["detail"] = v.Type().Field(i).Tag.Get("radio")
			case "manytomany":
				tmp["detail"] = v.Type().Field(i).Tag.Get("manytomany")
			case "select":
				tmp["detail"] = v.Type().Field(i).Tag.Get("select")
			case "multiselect":
				tmp["detail"] = v.Type().Field(i).Tag.Get("multiselect")
			case "o2o":
				tmp["detail"] = v.Type().Field(i).Tag.Get("o2o")
			case "o2m":
				tmp["detail"] = v.Type().Field(i).Tag.Get("o2m")
			}

			//字段名 供前端form表单使用 还有table字段
			if v.Type().Field(i).Tag.Get("name") != "" {
				tmp["Name"] = v.Type().Field(i).Tag.Get("name")
			} else {
				tmp["Name"] = v.Type().Field(i).Name
			}
			tmp["Names"] += fmt.Sprintf("%s ", tmp["Name"])
			//收集昵称
			if v.Type().Field(i).Tag.Get("verbose_name") != "" {
				tmp["Col"] = fmt.Sprintf("%s %s:%s:%s:%s", tmp["Col"], v.Type().Field(i).Tag.Get("verbose_name"), tmp["ColType"], tmp["Name"], tmp["detail"])
			} else {
				tmp["Col"] = fmt.Sprintf("%s %s:%s:%s:%s", tmp["Col"], v.Type().Field(i).Name, tmp["ColType"], tmp["Name"], tmp["detail"])
			}
			//collect show list on table columns
			if v.Type().Field(i).Tag.Get("list") == "true" || v.Type().Field(i).Tag.Get("list") == "" {
				//name
				tmp["List"] += fmt.Sprintf("%s ", tmp["Name"])
			}
			//collect show search on table columns
			if v.Type().Field(i).Tag.Get("search") == "true" {
				//name
				if tmp["Search"] == "" {
					tmp["Search"] = tmp["Name"]
				} else {
					tmp["Search"] = fmt.Sprintf("%s,%s", tmp["Search"], tmp["Name"])
				}
			}
		}
	}
	return tmp
}

func Register(data ...interface{}) error {
	tmp := GetRegistered()
	for _, model := range data {
		// log.Debugf("Register %v", RigsterStruct(model))
		tmp = append(tmp, RigsterStruct(model))
	}

	// 注册Model
	err := sqlite.NewOrm().Sync2(data...)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// log.Debug("完成模型注册")
	registered = tmp
	return nil
}

// func ReadStruct(data ...interface{}) []map[string]string {
// 	result := []map[string]string{}
// 	for _, x := range data {
// 		result = append(result, RigsterStruct(x))
// 	}
// 	return result
// }
