package addform

import (
	"fmt"
	"log/slog"
	"strings"

	orm "github.com/lflxp/djangolang/consted"
)

type AddForm struct{}

func (a *AddForm) Transfer() interface{} {
	return func(data map[string]string) string {
		var result string
		if err := a.Check(data); err != nil {
			slog.Error(err.Error())
			return ""
		}

		slog.Debug("Col formcolumns %v", data)
		args := data["Col"]

		for _, info := range strings.Split(strings.TrimSpace(args), " ") {
			tmp := strings.Split(info, ":")
			// 移除id字段在form表单添加和修改的时候
			if tmp[2] != "id" {
				switch tmp[1] {
				case string(orm.Password):
					result += a.GetFunc(orm.Password)(a.GetTemplate(orm.Text), tmp)
				case string(orm.StringField):
					result += a.GetFunc(orm.StringField)(a.GetTemplate(orm.Text), tmp)
				case string(orm.IntField):
					result += a.GetFunc(orm.IntField)(a.GetTemplate(orm.Text), tmp)
				case string(orm.Int16Field):
					result += a.GetFunc(orm.Int16Field)(a.GetTemplate(orm.Text), tmp)
				case string(orm.Int64Field):
					result += a.GetFunc(orm.Int64Field)(a.GetTemplate(orm.Text), tmp)
				case string(orm.Textarea):
					result += a.GetFunc(orm.Textarea)(a.GetTemplate(orm.Textarea), tmp)
				case string(orm.Radio):
					result += a.GetFunc(orm.Radio)(a.GetTemplate(orm.Radio), tmp)
				case string(orm.Select):
					result += a.GetFunc(orm.Select)(a.GetTemplate(orm.Select), tmp)
				case string(orm.MultiSelect):
					result += a.GetFunc(orm.MultiSelect)(a.GetTemplate(orm.MultiSelect), tmp)
				case string(orm.File):
					result += a.GetFunc(orm.File)(a.GetTemplate(orm.Text), tmp)
				case string(orm.Time):
					result += a.GetFunc(orm.Time)(a.GetTemplate(orm.Time), tmp)
				case string(orm.OneToOne):
					result += a.GetFunc(orm.OneToOne)(a.GetTemplate(orm.Select), tmp)
				case string(orm.OneToMany):
					result += a.GetFunc(orm.OneToMany)(a.GetTemplate(orm.MultiSelect), tmp)
				default:
					slog.Error("unsupport Columns Type", "Type", tmp[1])
				}
			}
		}
		return result
	}
}

// 检查并赋值参数
func (a *AddForm) Check(data map[string]string) error {
	if value, ok := data["Col"]; ok {
		slog.Debug("Col ", "Value", value)
	} else {
		return fmt.Errorf("Col is empty")
	}

	if value, ok := data["List"]; ok {
		slog.Debug("List ", value)
	}
	if value, ok := data["Search"]; ok {
		slog.Debug("Search ", value)
	}

	return nil
}

func (a *AddForm) GetTemplate(t orm.FieldType) string {
	column, ok := addColumns[t]
	if !ok {
		slog.Error("模板字段 %v 不存在", t)
		return fmt.Sprintf("模板字段 %v 不存在", t)
	}
	return column
}

func (a *AddForm) GetFunc(t orm.FieldType) func(text string, tmp []string) string {
	function, ok := addFuncMap[t]
	if !ok {
		slog.Error("模板字段 %v 不存在", a)
	}
	return function
}
