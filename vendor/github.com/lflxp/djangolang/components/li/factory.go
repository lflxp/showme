package li

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

// string转html
type Li struct {
	Raw    []map[string]string
	Result string
}

func (s *Li) Transfer() interface{} {
	return func(info []map[string]string) string {
		var result string
		if err := s.Check(info); err != nil {
			slog.Error(err.Error())
			return err.Error()
		}

		// 按Service分组
		tmp := map[string][]map[string]string{}

		for _, in := range s.Raw {
			// log.Debugf("Service %s", in["Services"])
			if d, ok := tmp[in["Services"]]; ok {
				d = append(d, in)
				tmp[in["Services"]] = d
				// log.Debugf("Service ADD %v", tmp)
			} else {
				t := []map[string]string{in}
				tmp[in["Services"]] = t
				// log.Debugf("New Service %v", tmp)
			}
		}

		// 去重
		var sortsKey []string
		for key := range tmp {
			sortsKey = append(sortsKey, key)
		}

		// 排序
		sort.Strings(sortsKey)

		for _, key := range sortsKey {
			value := tmp[key]
			// result := fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">%s</li>", strings.ToUpper(beego.AppConfig.String("appname")))
			result += fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">Project %s</li>", strings.ToUpper(key))
			for _, data := range value {
				// result += fmt.Sprintf("<li class=\"list-group-item list-group-item-info\">%s</li>", data["Struct"])
				if data["Struct"] != "User" && data["Struct"] != "Group" && data["Struct"] != "Userauth" && data["Struct"] != "History" {
					// result += fmt.Sprintf("<li class=\"list-group-item\"><a class=\"badge\" href=\"/admin/table?name=%s\"><span class=\"glyphicon glyphicon-edit\"></span>Change</a><a class=\"badge\" href=\"/admin/add?name=%s\"><span class=\"glyphicon glyphicon-plus\"></span>Add</a><a href=\"/admin/table?name=%s\" target=\"_self\">%s</a></li>", data["Struct"], data["Struct"], data["Struct"], data["Struct"])
					result += fmt.Sprintf("<li class=\"list-group-item\"><button type=\"button\" style=\"float: right;\" class=\"btn btn-link btn-xs\" data-toggle=\"modal\" data-target=\"#addModelinit\" onclick=\"Add('%s')\"> <span class=\"glyphicon glyphicon-plus\"></span>添加</button><a href=\"/admin/table?name=%s\" target=\"_self\">%s</a></li>", data["Struct"], data["Struct"], data["Struct"])
				}
				// for _, x := range strings.Split(data["Name"], " ") {
				// 	result += fmt.Sprintf("<li class=\"list-group-item\"><a class=\"badge\" href=\"#\">Change</a><a class=\"badge\" href=\"#\">Add</a>%s</li>", x)
				// }
			}
		}

		result += "<div class=\"row\">&nbsp;</div>"
		return result
	}
}

func (s *Li) Check(data []map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("data is none")
	}
	s.Raw = data
	return nil
}
