package tty

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 标准化
func wstty(c *gin.Context) {
	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	cmd := exec.Command("bash")
	//这里得到标准输出和标准错误输出的两个管道，此处获取了错误处理
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Error(err.Error())
		return
	}

	connections := int64(0)
	connections = atomic.AddInt64(&connections, 1)
	xterm := XtermJs{
		Options: Options{
			PermitWrite:    true,
			CloseSignal:    9,
			MaxConnections: 100,
		},
		Title:       "message",
		Connections: &connections,
		Server:      NewServer(),
	}

	log.Printf("Command is running for client %s with PID %d (args=%q), connections: %d",
		c.Request.RemoteAddr, cmd.Process.Pid, "bash", connections)
	xterm.Server.StartGo()

	context := &ClientContext{
		Xtermjs:    &xterm,
		Request:    c.Request,
		WsConn:     ws,
		Cmd:        cmd,
		Pty:        ptmx,
		Cache:      bytes.NewBuffer([]byte("")),
		CacheMutex: &sync.Mutex{},
		WriteMutex: &sync.Mutex{},
	}

	context.HandleClient()
	xterm.Server.WaitGo()
}

func ServeGin(host, port, username, password string, isdebug bool) {
	if isdebug {
		// 设置日志级别为warn以上
		log.SetLevel(log.DebugLevel)
	} else {
		// 设置日志级别为warn以上
		log.SetLevel(log.InfoLevel)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 使用 Recovery 中间件
	router.Use(gin.Recovery())

	// 静态二进制文件
	fs := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
	router.StaticFS("/static", &fs)
	// 静态文件
	// router.StaticFS("/static", http.Dir("./tty/static"))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/index")
	})
	apiGroup := router.Group("/api", gin.BasicAuth(gin.Accounts{username: password}))
	apiGroup.GET("/ws", wstty)

	// 主页
	// 从内存取出然后渲染加载
	indexhtml := multitemplate.New()
	xterm3, err := Asset("xterm3.html")
	if err != nil {
		log.Error(err.Error())
		return
	}

	t, err := template.New("index").Parse(string(xterm3))
	if err != nil {
		log.Error(err.Error())
		return
	}

	indexhtml.Add("index", t)
	router.HTMLRender = indexhtml
	apiGroup.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", gin.H{"host": c.Request.RemoteAddr})
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		// log.Println("receive interrupt signal")
		// if err := server.Close(); err != nil {
		// 	log.Fatal("Server Close:", err)
		// }

		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}()

	log.Infof("Listening and serving HTTPS on %s:%s", host, port)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}

	log.Println("Server exiting")
}
