package editform

import (
	"fmt"
	"log/slog"
	"strings"

	orm "github.com/lflxp/djangolang/consted"
	"github.com/lflxp/djangolang/utils/orm/sqlite"
)

// 编辑表单
type EditForm struct{}

// 按照表字段循环生成所有表单字段
func (f *EditForm) Transfer() interface{} {
	return func(data map[string]string, id, name string) string {
		slog.Debug("Raw EditForm id %s data %v\n", id, data)
		if err := f.Check(data, id, name); err != nil {
			slog.Error(err.Error())
			return ""
		}

		args := data["Col"]
		result := fmt.Sprintf("<input name=\"id\" value=\"%s\" style=\"display:none\" />", id)
		check_sql := fmt.Sprintf("select * from %s where id=%s", name, id)
		db_result, err := sqlite.NewOrm().Query(check_sql)
		if err != nil {
			slog.Error(err.Error())
			return err.Error()
		}

		for _, info := range strings.Split(strings.TrimSpace(args), " ") {
			tmp := strings.Split(info, ":")
			// 移除id字段在form表单添加和修改的时候
			if tmp[2] != "id" {
				switch tmp[1] {
				case string(orm.Password):
					result += f.GetFunc(orm.Password)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.StringField):
					result += f.GetFunc(orm.StringField)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.IntField):
					result += f.GetFunc(orm.IntField)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.Int16Field):
					result += f.GetFunc(orm.Int16Field)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.Int64Field):
					result += f.GetFunc(orm.Int64Field)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.Textarea):
					result += f.GetFunc(orm.Textarea)(f.GetTemplate(orm.Textarea), tmp, db_result[0])
				case string(orm.Radio):
					result += f.GetFunc(orm.Radio)(f.GetTemplate(orm.Radio), tmp, db_result[0])
				case string(orm.Select):
					result += f.GetFunc(orm.Select)(f.GetTemplate(orm.Select), tmp, db_result[0])
				case string(orm.MultiSelect):
					result += f.GetFunc(orm.MultiSelect)(f.GetTemplate(orm.MultiSelect), tmp, db_result[0])
				case string(orm.File):
					result += f.GetFunc(orm.File)(f.GetTemplate(orm.Text), tmp, db_result[0])
				case string(orm.Time):
					result += f.GetFunc(orm.Time)(f.GetTemplate(orm.Time), tmp, db_result[0])
				case string(orm.OneToOne):
					result += f.GetFunc(orm.OneToOne)(f.GetTemplate(orm.Select), tmp, db_result[0])
				case string(orm.OneToMany):
					result += f.GetFunc(orm.OneToMany)(f.GetTemplate(orm.MultiSelect), tmp, db_result[0])
				default:
					slog.Error("not found field type: %s", tmp[1])
				}
			}
		}
		return result
	}
}

// 检查并赋值参数
func (f *EditForm) Check(data map[string]string, id, name string) error {
	slog.Debug("EditForm Check id %s name %s data: %v", id, name, data)
	if value, ok := data["Col"]; ok {
		slog.Debug("edit Col ", value)
	} else {
		return fmt.Errorf("edit Col is empty")
	}

	// if value, ok := data["List"]; ok {
	// 	slog.Debug("List ", value)
	// }
	// if value, ok := data["Search"]; ok {
	// 	slog.Debug("Search ", value)
	// }
	if name == "" {
		slog.Error("Table name is empty")
		return fmt.Errorf("Table name is empty")
	}

	if id == "" {
		slog.Error("id is empty")
		return fmt.Errorf("id is empty")
	}

	return nil
}

func (f *EditForm) GetTemplate(t orm.FieldType) string {
	column, ok := editColumns[t]
	if !ok {
		slog.Error("模板字段 %v 不存在", t)
		return fmt.Sprintf("模板字段 %v 不存在", t)
	}
	return column
}

func (f *EditForm) GetFunc(t orm.FieldType) func(text string, tmp []string, data map[string][]byte) string {
	function, ok := editFuncMap[t]
	if !ok {
		slog.Error("模板字段 %v 不存在", f)
	}
	return function
}
