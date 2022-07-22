package pkg

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strings"
)

// var suff []string = []string{".avi", ".wma", ".rmvb", ".rm", ".mp4", ".mov", ".3gp", ".mpeg", ".mpg", ".mpe", ".m4v", ".mkv", ".flv", ".vob", ".wmv", ".asf", ".asx"}

// 获取指定目录下的所有文件,包含子目录下的文件
// strings.Replace(dirPth+PthSep+fi.Name(), ".", "/static", 1)
// 只能查找当前目录下的所有视频文件
func GetAllFiles(dirPth string, suff []string) (files []string, err error) {
	// var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			// dirs = append(dirs, dirPth+PthSep+fi.Name())
			// GetAllFiles(dirPth+PthSep+fi.Name(), suff)
			fmt.Println("不递归")
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
	// for _, table := range dirs {
	// 	temp, _ := GetAllFiles(table, suff)
	// 	for _, temp1 := range temp {
	// 		files = append(files, temp1)
	// 	}
	// }

	return files, nil
}

func GetIPs() (ips []string) {

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			// fmt.Println(ipNet.IP.String(), ipNet.Mask.String())
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

// 快速判断字符是否在字符数组中
func In(target string, source []string) bool {
	sort.Strings(source)
	index := sort.SearchStrings(source, target)
	if index < len(source) && source[index] == target {
		return true
	}
	return false
}
