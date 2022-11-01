package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// 写文件
func WriteFile(path string, data []byte) (int, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return 0, err
	}

	return f.Write(data)
}

// 渲染模板
func ApplyTemplate(temp string, data map[string]interface{}) (string, error) {
	var out bytes.Buffer
	// TODO: More FuncMap => https://blog.csdn.net/shida_csdn/article/details/113760434
	// ctx := map[string]interface{}{}
	// template.Must(template.New("demo").Funcs(funcMap).Parse(demoTmpl)).Execute(os.Stdout, ctx)
	t := template.Must(template.New("now").Parse(temp))
	err := t.Execute(&out, data)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// 加密base64
func EncodeBase64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

// 解密base64
func DecodeBase64(in string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(in)
	return string(decoded), err
}

// 加密
func Jiami(code string) string {
	w := md5.New()
	io.WriteString(w, code)
	md5str2 := fmt.Sprintf("%x", w.Sum(nil))
	return md5str2
}

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(strings.TrimSpace(text)))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(32)
}

// 生成随机字符串
func GetRandomString(len int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len; i++ {
		result = append(result, bytes[r.Intn(62)])
	}
	return string(result)
}

func ExecCommand(cmd string) ([]byte, error) {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Stdout = &out
	pipeline.Stderr = &stderr
	err := pipeline.Run()
	if err != nil {
		return stderr.Bytes(), err
	}
	// fmt.Println(stderr.String())
	return out.Bytes(), nil
}

func ExecCommandString(cmd string) (string, error) {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Stdout = &out
	pipeline.Stderr = &stderr
	err := pipeline.Run()
	if err != nil {
		return stderr.String(), err
	}
	return out.String(), nil
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

// 0:cf:e0:44:dd:be,enp1s0
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, fmt.Sprintf("%s,%s", macAddr, netInterface.Name))
	}
	return macAddrs
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

// If 实现三元表达式的功能
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func ParseJson(data string, header http.Header) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}

	var rs map[string]interface{}
	// 判断数据是否是gzip压缩
	if header != nil && strings.Contains(header.Get("Content-Encoding"), "gzip") {
		var rs map[string]interface{}
		reader, err := gzip.NewReader(strings.NewReader(data))
		if err != nil {
			return rs, err
		}
		decodeRs, err := ioutil.ReadAll(reader)
		if err != nil {
			return rs, err
		}
		data = string(decodeRs)
	}

	if err := json.Unmarshal([]byte(data), &rs); err != nil {
		return rs, err
	}
	return rs, nil
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

func IsDir(path string) bool {
	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("%s is not exist", path)
		return false
	} else if err != nil {
		return false
	}

	return file.IsDir()
}
