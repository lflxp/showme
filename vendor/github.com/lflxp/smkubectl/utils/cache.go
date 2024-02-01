/*
1. ~/.smkubectl.yaml文件的查询、创建、修改功能
2. 获取缓存数据，计算缓存时间
3. 每24小时更新一次
4. 记录每个集群的api-resources
*/
package utils

import (
	"io"
	"log/slog"
	"os"
	"time"

	"path/filepath"

	yaml "gopkg.in/yaml.v3"
	"k8s.io/client-go/util/homedir"
)

const kubectlYaml = ".kubectl-smart.yaml"

var targetFile string

func init() {
	if home := homedir.HomeDir(); home != "" {
		targetFile = filepath.Join(home, kubectlYaml)
	} else {
		targetFile = kubectlYaml
	}
}

type ResourceCache struct {
	Cluster   string   `yaml:"cluster"`
	Resources []string `yaml:"resources"`
}

type ResourceYaml struct {
	UpdateTime time.Time       `yaml:"updatetime"` // 最近一次更新时间
	Data       []ResourceCache `yaml:"data"`
}

// 判断是否在24小时以内
func (r *ResourceYaml) IsExpire() bool {
	slog.Debug("判断缓存是否过期", "UpdateTime", r.UpdateTime.Add(24*time.Hour), "Now", time.Now(), "IsExpire", time.Now().After(r.UpdateTime.Add(24*time.Hour)))
	if len(r.Data) == 0 {
		return true
	}
	return time.Now().After(r.UpdateTime.Add(24 * time.Hour))
}

// Yaml文件转ResourceYaml
// 不存在则创建
func ReadYamlToStruct() (*ResourceYaml, error) {
	var result ResourceYaml
	slog.Debug("打开缓存配置文件", "File", targetFile, "操作", "READ")
	file, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return &result, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return &result, err
	}

	if len(bytes) == 0 {
		result = ResourceYaml{
			UpdateTime: time.Now(),
			Data:       []ResourceCache{},
		}
		return &result, nil
	}

	err = yaml.Unmarshal(bytes, &result)
	return &result, err
}

// ResourceYaml转Yaml
// 1. 获取本地文件
// 2. 判断本地文件中是否含有这个集群数据
// 3. 有则更新 没有就新增
func WriteResourcesToYaml(data *ResourceCache) error {
	list, err := ReadYamlToStruct()
	if err != nil {
		return err
	}

	isExist := false
	for index, info := range list.Data {
		if info.Cluster == data.Cluster {
			isExist = true
			// 更新时间
			list.UpdateTime = time.Now()
			list.Data[index] = *data
		}
	}

	if !isExist {
		list.Data = append(list.Data, *data)
	}

	yml, err := yaml.Marshal(list)
	if err != nil {
		return err
	}

	return writeFile(yml)
}

func writeFile(data []byte) error {
	slog.Debug("打开缓存配置文件", "File", targetFile, "操作", "Write")
	file, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return os.WriteFile(targetFile, data, 0644)
}
