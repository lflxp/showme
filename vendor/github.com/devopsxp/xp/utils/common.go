package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

// 渲染模板
func ApplyTemplate(temp string, data map[string]interface{}) (string, error) {
	var out bytes.Buffer
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
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// return len=8  salt
func GetRandomSalt() string {
	return GetRandomString(32)
}

//生成随机字符串
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

func ExecCommandStd(cmd string) error {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	pipeline.Stdin = os.Stdin
	pipeline.Stdout = os.Stdout
	pipeline.Stderr = os.Stderr
	err := pipeline.Run()
	if err != nil {
		return err
	}
	// fmt.Println(stderr.String())
	return nil
}

func ExecCommand(cmd string) ([]byte, error) {
	pipeline := exec.Command("/bin/sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	pipeline.Stdin = os.Stdin
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
	pipeline.Stdin = os.Stdin
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

func ScanPort(host string, port int) bool {
	remote := fmt.Sprintf("%s:%d", host, port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", remote) //转换IP格式 // 根据域名查找ip
	//fmt.Printf("%s", tcpAddr)
	// conn, err := net.DialTCP("tcp", nil, tcpAddr) //查看是否连接成功
	conn, err := net.DialTimeout("tcp", tcpAddr.String(), time.Second) //建立tcp连接
	if err != nil {
		// fmt.Printf("no==%s:%s\r\n", host, port)
		return false
	}
	defer conn.Close()
	// fmt.Printf("ok==%s:%s\r\n", host, port)
	return true
}

// 解析IP地址
func ParseIps(in string) ([]string, error) {
	rs := []string{}
	if strings.Contains(in, "-") {
		tmp_a := strings.Split(in, ".")
		if len(tmp_a) != 4 {
			fmt.Println(tmp_a)
			return nil, errors.New("ip地址不正确")
		}
		A := []string{}
		B := []string{}
		C := []string{}
		D := []string{}
		for m, n := range tmp_a {
			if strings.Contains(n, "-") {
				tmp := strings.Split(n, "-")
				a, err := strconv.Atoi(tmp[0])
				if err != nil {
					return rs, err
				}
				b, err := strconv.Atoi(tmp[1])
				if err != nil {
					return rs, err
				}
				for i := a; i <= b; i++ {
					if m == 0 {
						A = append(A, fmt.Sprintf("%d", i))
					} else if m == 1 {
						B = append(B, fmt.Sprintf("%d", i))
					} else if m == 2 {
						C = append(C, fmt.Sprintf("%d", i))
					} else if m == 3 {
						D = append(D, fmt.Sprintf("%d", i))
					}
				}
			} else {
				if m == 0 {
					A = append(A, n)
				} else if m == 1 {
					B = append(B, n)
				} else if m == 2 {
					C = append(C, n)
				} else if m == 3 {
					D = append(D, n)
				}
			}
		}
		for _, a1 := range A {
			for _, b1 := range B {
				for _, c1 := range C {
					for _, d1 := range D {
						rs = append(rs, fmt.Sprintf("%s.%s.%s.%s", a1, b1, c1, d1))
					}
				}
			}
		}
	} else {
		rs = append(rs, in)
	}
	return rs, nil
}

// 转换星期中文为数字
func TransformCHN(data []string) ([]int, error) {
	rs := []int{}
	if len(data) == 0 {
		return nil, errors.New("data is 0 length")
	}

	for _, x := range data {
		switch x {
		case "星期一":
			rs = append(rs, 1)
		case "星期二":
			rs = append(rs, 2)
		case "星期三":
			rs = append(rs, 3)
		case "星期四":
			rs = append(rs, 4)
		case "星期五":
			rs = append(rs, 5)
		case "星期六":
			rs = append(rs, 6)
		case "星期天":
			rs = append(rs, 0)
		}
	}
	return rs, nil
}

// 判断是否在两个时间范围内
func IsBetweenAB(start, end string) (bool, error) {
	rs := false
	format := "2006-01-02"
	now := time.Now().Local()

	days := now.Format(format)

	format_mm := "2006-01-02 15:04"
	startTime, err := time.ParseInLocation(format_mm, fmt.Sprintf("%s %s", days, start), time.Local)
	if err != nil {
		return rs, err
	}
	endTime, err := time.ParseInLocation(format_mm, fmt.Sprintf("%s %s", days, end), time.Local)
	if err != nil {
		return rs, err
	}

	if now.After(startTime) && now.Before(endTime) {
		rs = true
	}
	log.Println("验证时间", rs, fmt.Sprintf("%s %s", days, start), fmt.Sprintf("%s %s", days, end), now.Format(format_mm))
	return rs, nil
}

//PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}
