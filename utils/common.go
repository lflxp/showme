package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
#include <stdio.h>
#include <termios.h>
struct termios disable_echo() {
  struct termios of, nf;
  tcgetattr(fileno(stdin), &of);
  nf = of;
  nf.c_lflag &= ~ECHO;
  nf.c_lflag |= ECHONL;
  if (tcsetattr(fileno(stdin), TCSANOW, &nf) != 0) {
    perror("tcsetattr");
  }
  return of;
}
void restore_echo(struct termios f) {
  if (tcsetattr(fileno(stdin), TCSANOW, &f) != 0) {
    perror("tcsetattr");
  }
}
*/
import "C"

func Prompt(msg string) string {
	fmt.Printf("%s: ", msg)
	oldFlags := C.disable_echo()
	passwd, err := bufio.NewReader(os.Stdin).ReadString('\n')
	C.restore_echo(oldFlags)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(passwd)
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

// 解密base64
func DecodeBase64Bytes(in string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(in)
	return decoded, err
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
