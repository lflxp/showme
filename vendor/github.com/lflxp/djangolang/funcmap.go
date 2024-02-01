package djangolang

import (
	"sync"
	"text/template"

	"github.com/lflxp/djangolang/components"
	"github.com/lflxp/djangolang/consted"
)

var (
	once    sync.Once
	FuncMap template.FuncMap
)

func init() {
	once.Do(func() {
		if FuncMap == nil {
			FuncMap = map[string]interface{}{}
		}
	})
	// 注册自定义函数
	FuncMap["beegoli"] = components.NewComponentAdaptor(consted.BeegoLi).Transfer()
	FuncMap["admincolumns"] = components.NewComponentAdaptor(consted.AdminColumns).Transfer()
	FuncMap["formcolumns"] = components.NewComponentAdaptor(consted.FormComponent).Transfer()
	FuncMap["str2html"] = components.NewComponentAdaptor(consted.StrToHtml).Transfer()
	FuncMap["uploads"] = components.NewComponentAdaptor(consted.Upload).Transfer()
}

func AddFuncMap(funcName string, funcs interface{}) {
	if _, ok := FuncMap[funcName]; !ok {
		FuncMap[funcName] = funcs
	}
}
