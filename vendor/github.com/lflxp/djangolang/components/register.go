package components

import (
	"fmt"

	"github.com/lflxp/djangolang/components/addcolumns"
	"github.com/lflxp/djangolang/components/addform"
	"github.com/lflxp/djangolang/components/editform"
	"github.com/lflxp/djangolang/components/li"
	"github.com/lflxp/djangolang/components/str2html"
	"github.com/lflxp/djangolang/components/uploads"
	"github.com/lflxp/djangolang/consted"
)

func init() {
	HtmlFactories[consted.FormComponent] = &addform.AddForm{}
	HtmlFactories[consted.EditFormComponent] = &editform.EditForm{}
	HtmlFactories[consted.StrToHtml] = &str2html.Str2html{}
	HtmlFactories[consted.AdminColumns] = &addcolumns.AddColumns{}
	HtmlFactories[consted.BeegoLi] = &li.Li{}
	HtmlFactories[consted.Upload] = &uploads.Uploads{}
}

// 定义字段工厂
// 保存用于创建Filed的工厂实例，其中map的key为插件类型，value为抽象工厂接口
var HtmlFactories = make(map[consted.HtmlComponentType]FormFactory)

// 根据FieldType类型返回对应Field类型的工厂实例
func factoryOf(f consted.HtmlComponentType) FormFactory {
	factory, ok := HtmlFactories[f]
	if !ok {
		panic(fmt.Sprintf("orm: unknown field type: %v", f))
	}
	return factory
}

// 组件适配器
type ComponentAdaptor struct {
	Type consted.HtmlComponentType
}

func (c *ComponentAdaptor) Transfer() interface{} {
	switch c.Type {
	case consted.FormComponent, consted.EditFormComponent, consted.StrToHtml, consted.AdminColumns, consted.BeegoLi, consted.Upload:
		return factoryOf(c.Type).Transfer()
	default:
		return fmt.Errorf("unsupport components type: %s", c.Type)
	}
}

func NewComponentAdaptor(t consted.HtmlComponentType) *ComponentAdaptor {
	return &ComponentAdaptor{Type: t}
}
