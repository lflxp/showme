package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

var (
	types      []string
	pageSize   int
	isvideo    bool
	path       string
	port       string
	closeChan  chan os.Signal
	uri        string
	initnum    int
	raw        bool
	staticPort string
)

const maxUploadSize = 2000 * 1024 * 2014 // 2 MB
const uploadPath = "/tmp"

func init() {
	initnum = 0
	// path = utils.GetCurrentDirectory()
	// port = "9090"
	closeChan = make(chan os.Signal)
}

func HttpStaticServeForCorba(data *Apis) {
	// httpstatic -port 9090 -path ./
	port = data.Port
	path = data.Path
	isvideo = data.IsVideo
	pageSize = data.PageSize
	types = strings.Split(data.Types, ",")
	raw = data.Raw
	staticPort = data.StaticPort

	ips := GetIPs()
	var openUrl string
	for index, ip := range ips {
		if index == 0 {
			openUrl = fmt.Sprintf("http://%s:%s", ip, port)
		}
		log.Printf("前端访问地址: http://%s:%s", ip, port)
		log.Printf("文件访问地址: http://%s:%s", ip, staticPort)
	}
	dir, _ := os.Getwd()
	log.Printf("当前目录: %s", dir)
	log.Println("curl -X POST http://127.0.0.1:9090/upload -F \"file=@/Users/lxp/123.mp4\" -H \"Content-Type:multipart/form-data\"")

	open.Start(openUrl)
	serverGin()
}
