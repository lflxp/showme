// 读取指定yaml文件
package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// 指定读取yaml文件
func ReadYamlConfig(path string) (interface{}, error) {
	var data interface{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
