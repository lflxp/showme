// 获取指定目录下的所有视频文件
package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

var suff []string = []string{".avi", ".wma", ".rmvb", ".rm", ".mp4", ".mov", ".3gp", ".mpeg", ".mpg", ".mpe", ".m4v", ".mkv", ".flv", ".vob", ".wmv", ".asf", ".asx"}

// 获取指定目录下的所有文件,包含子目录下的文件
// strings.Replace(dirPth+PthSep+fi.Name(), ".", "/static", 1)
// 只能查找当前目录下的所有视频文件
func GetAllFiles(dirPth string) (files []string, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetAllFiles(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			for _, x := range suff {
				ok := strings.HasSuffix(fi.Name(), x)
				if ok {
					files = append(files, strings.Replace(dirPth+PthSep+fi.Name(), ".", "/static", 1))
				}
			}

		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}
