package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
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
