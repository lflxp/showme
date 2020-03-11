package tty

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/DeanThompson/ginpprof"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"github.com/lflxp/showme/utils"
	log "github.com/sirupsen/logrus"
)

var xterm *XtermJs

func init() {
	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	log.SetReportCaller(false)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeGin(port, username, password string, cmds []string, isdebug, isReconnect, isPermitWrite, isAudit, isXsrf, isProf bool, MaxConnections int64) {
	if isdebug {
		// 设置日志级别为warn以上
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		// 设置日志级别为warn以上
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	if isAudit {
		utils.InitSqlite()
		utils.Engine.Sync2(new(Aduit))
	}

	router := gin.Default()

	// 使用 Recovery 中间件
	router.Use(gin.Recovery())

	// 判断cmds输入，为空默认设置为bash
	if len(cmds) == 0 {
		cmds = append(cmds, "bash")
	}

	// 初始化XtermJs全局属性配置
	connections := int64(0)
	xterm = &XtermJs{
		Options: Options{
			PermitWrite:    isPermitWrite,
			CloseSignal:    9,
			MaxConnections: MaxConnections,
			Audit:          isAudit,
			Xsrf:           isXsrf,
		},
		Title:       "Showme",
		Connections: &connections,
		Server:      NewServer(),
		XsrfToken:   sync.Map{},
	}

	// 静态二进制文件
	fs := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
	router.StaticFS("/static", &fs)
	// 静态文件
	// router.StaticFS("/static", http.Dir("./tty/static"))

	var apiGroup *gin.RouterGroup
	// 是否密码登录
	if username == "" && password == "" {
		apiGroup = router.Group("/")
	} else {
		apiGroup = router.Group("/", gin.BasicAuth(gin.Accounts{username: password}))
	}

	// 添加审计查询接口
	if isAudit {
		apiGroup.GET("/check", func(c *gin.Context) {
			name := c.DefaultQuery("name", "")
			data, err := GetAduit(name)
			if err != nil {
				c.String(http.StatusOK, err.Error())
			} else {
				c.JSONP(http.StatusOK, data)
			}
		})
	}

	// 后端websocket服务
	apiGroup.GET("/ws", func(c *gin.Context) {
		conns := atomic.AddInt64(xterm.Connections, 1)
		if xterm.Options.MaxConnections != 0 {
			if conns > xterm.Options.MaxConnections {
				log.Printf("Max Connected: %d", xterm.Options.MaxConnections)
				atomic.AddInt64(xterm.Connections, -1)
				return
			}
		}
		// 升级get请求为webSocket协议
		ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		cmd := exec.Command(cmds[0], cmds[1:]...)
		//这里得到标准输出和标准错误输出的两个管道，此处获取了错误处理
		ptmx, err := pty.Start(cmd)
		if err != nil {
			log.Errorf("ptmx[52] %s", err.Error())
			return
		}

		if xterm.Options.MaxConnections != 0 {
			log.Printf("Command is running for client %s with PID %d (args=%q), connections: %d/%d",
				c.Request.RemoteAddr, cmd.Process.Pid, cmds, conns, xterm.Options.MaxConnections)
		} else {
			log.Printf("Command is running for client %s with PID %d (args=%q), connections: %d",
				c.Request.RemoteAddr, cmd.Process.Pid, cmds, conns)
		}

		xterm.Server.StartGo()

		context := &ClientContext{
			Xtermjs: xterm,
			Request: c.Request,
			WsConn:  ws,
			Cmd:     cmd,
			Pty:     ptmx,
			// Cache:      bytes.NewBuffer([]byte("")),
			// CacheMutex: &sync.Mutex{},
			WriteMutex: &sync.Mutex{},
		}

		context.HandleClient()
		xterm.Server.WaitGo()
	})

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
	apiGroup.GET("/", func(c *gin.Context) {
		newXsrf := utils.GetRandomSalt()
		log.Debugf("%s xsrftoken %s", c.Request.RemoteAddr, newXsrf)
		if !xterm.Options.Xsrf {
			xterm.XsrfToken.Store(fmt.Sprintf("%s%s", newXsrf, strings.Split(c.Request.RemoteAddr, ":")[0]), time.Now().String())
		}
		c.HTML(http.StatusOK, "index", gin.H{
			"host":      c.Request.RemoteAddr,
			"Reconnect": isReconnect,
			"Debug":     isdebug,
			"Write":     isPermitWrite,
			"MaxC":      MaxConnections,
			"Conn":      *xterm.Connections + 1,
			"Cmd":       strings.Join(cmds, " "),
			"Xsrf":      newXsrf,
		})
	})

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	if isProf {
		ginpprof.Wrapper(router)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
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

	ips := utils.GetIPs()
	for _, ip := range ips {
		log.Infof("Listening and serving HTTPS on %s:%s", ip, port)
	}

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}

	log.Println("Server exiting")
}
