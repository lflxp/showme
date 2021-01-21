package plugin

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

// filter插件名称与类型的映射关系，主要用于通过反射创建filter对象
var checkNames = make(map[string]reflect.Type)

func AddCheck(key string, value reflect.Type) {
	log.Debugf("添加Check插件 %s", key)
	checkNames[key] = value
}

// input插件名称与类型的映射关系，主要用于通过反射创建input对象
var inputNames = make(map[string]reflect.Type)

func AddInput(key string, value reflect.Type) {
	log.Debugf("添加Input插件 %s", key)
	inputNames[key] = value
}

// filter插件名称与类型的映射关系，主要用于通过反射创建filter对象
var filterNames = make(map[string]reflect.Type)

func AddFilter(key string, value reflect.Type) {
	log.Debugf("添加Filter插件 %s", key)
	filterNames[key] = value
}

// output插件名称与类型的映射关系，主要用于通过反射创建output对象
var outputNames = make(map[string]reflect.Type)

func AddOutput(key string, value reflect.Type) {
	log.Debugf("添加Output插件 %s", key)
	outputNames[key] = value
}
