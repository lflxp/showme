// 工厂模式生产具体的类
package plugin

import "reflect"

// ============================根据接口实现struct==============================
// 接着，我们定义input、filter、output三类插件接口的具体实现：

// 插件抽象工厂接口
type Factory interface {
	Create(conf Config) Plugin
}

// check插件工厂，实现Factory接口
type CheckFactory struct{}

// 读取配置，通过反射机制进行对象实例化
func (i *CheckFactory) Create(conf Config) Plugin {
	t, _ := checkNames[conf.Name]
	p := reflect.New(t).Interface().(Plugin)
	return p
}

// input插件工厂对象，实现Factory接口
type InputFactory struct{}

// 读取配置，通过反射机制进行对象实例化
func (i *InputFactory) Create(conf Config) Plugin {
	t, _ := inputNames[conf.Name]
	p := reflect.New(t).Interface().(Plugin)
	return p
}

// filter插件工厂对象，实现Factory接口
type FilterFactory struct{}

// 读取配置，通过反射机制进行对象实例化
func (i *FilterFactory) Create(conf Config) Plugin {
	t, _ := filterNames[conf.Name]
	p := reflect.New(t).Interface().(Plugin)
	return p
}

// output插件工厂对象，实现Factory接口
type OutputFactory struct{}

// 读取配置，通过反射机制进行对象实例化
func (i *OutputFactory) Create(conf Config) Plugin {
	t, _ := outputNames[conf.Name]
	p := reflect.New(t).Interface().(Plugin)
	return p
}
