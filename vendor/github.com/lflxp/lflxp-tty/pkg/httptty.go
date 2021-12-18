// The function is used to provide the gin interface plug-in and dynamic parameter HTTP transfer
package pkg

import (
	"fmt"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chenjiandongx/ginprom"
	"github.com/creack/pty"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var httpXterm *XtermJs
var rootPath string

// isLocal 参数是用于判断是否为第三方引用，进而改变访问路径
func RegisterTty(router *gin.Engine, data *Tty, isLocal bool) {
	if data.IsDebug {
		// 设置日志级别为warn以上
		log.SetLevel(log.DebugLevel)
		gin.SetMode(gin.DebugMode)
	} else {
		// 设置日志级别为warn以上
		log.SetLevel(log.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	if data.IsAudit {
		InitBoltDB()
	}

	if data.EnableTLS {
		router.Use(TlsHandler(data.Host, data.Port))
	}

	// 初始化httpXtermJs全局属性配置
	connections := int64(0)
	httpXterm = &XtermJs{
		Options: Options{
			PermitWrite:    data.IsPermitWrite,
			CloseSignal:    9,
			MaxConnections: data.MaxConnections,
			Audit:          data.IsAudit,
			Xsrf:           data.IsXsrf,
			EnableTLS:      data.EnableTLS,
			CrtPath:        data.CrtPath,
			KeyPath:        data.KeyPath,
			IsReconnect:    data.IsReconnect,
			IsDebug:        data.IsDebug,
		},
		Title:       "LFLXP-TTY",
		Connections: &connections,
		Server:      NewServer(),
		XsrfToken:   sync.Map{},
		Cmds:        data.Cmds,
	}

	// 静态二进制文件
	fs := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	// router.StaticFS("/tty/static", &fs)

	if isLocal {
		rootPath = "/"
	} else {
		rootPath = "/tty"
	}
	var apiGroup *gin.RouterGroup
	// 是否密码登录
	if data.Username == "" && data.Password == "" {
		apiGroup = router.Group(rootPath)
	} else {
		apiGroup = router.Group(rootPath, gin.BasicAuth(gin.Accounts{data.Username: data.Password}))
	}

	// 添加html template
	// 主页
	// 从内存取出然后渲染加载
	indexhtml := multitemplate.New()
	httpXterm3, err := Asset("xterm3.html")
	if err != nil {
		log.WithField("tty.go", "198").Error(err.Error())
		return
	}

	t, err := template.New("index").Parse(string(httpXterm3))
	if err != nil {
		log.Error(err.Error())
		return
	}

	admin, err := Asset("admin.html")
	if err != nil {
		log.WithField("tty.go", "198").Error(err.Error())
		return
	}

	ta, err := template.New("admin").Parse(string(admin))
	if err != nil {
		log.Error(err.Error())
		return
	}

	indexhtml.Add("index", t)
	indexhtml.Add("admin", ta)
	router.HTMLRender = indexhtml

	if data.IsAudit {
		apiGroup.GET("/check", Check)
		apiGroup.GET("/who", Who)
		apiGroup.GET("/admin", Admin)
	}
	apiGroup.GET("/", Index)
	apiGroup.GET("/ws", Ws)
	apiGroup.StaticFS("/static", &fs)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	doMetrics()
	apiGroup.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
}

func Check(c *gin.Context) {
	defer func() {
		var path string
		if rootPath == "/" {
			path = "/check"
		} else {
			path = fmt.Sprintf("%s/check", rootPath)
		}
		who := &Whos{
			Remoteaddr: c.Request.RemoteAddr,
			Path:       path,
		}
		err := who.Save()
		if err != nil {
			log.Error(err)
		}
	}()
	name := c.DefaultQuery("name", "")
	data, err := GetAduit(name)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.JSONP(http.StatusOK, data)
	}
}

func Who(c *gin.Context) {
	defer func() {
		var path string
		if rootPath == "/" {
			path = "/who"
		} else {
			path = fmt.Sprintf("%s/who", rootPath)
		}
		who := &Whos{
			Remoteaddr: c.Request.RemoteAddr,
			Path:       path,
		}
		err := who.Save()
		if err != nil {
			log.Error(err)
		}
	}()
	name := c.DefaultQuery("name", "")
	data, err := GetWhos(name)
	if err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.JSONP(http.StatusOK, data)
	}
}

func Admin(c *gin.Context) {
	defer func() {
		var (
			path string
		)
		if rootPath == "/" {
			path = "/admin"
		} else {
			path = fmt.Sprintf("%s/admin", rootPath)
		}
		who := &Whos{
			Remoteaddr: c.Request.RemoteAddr,
			Path:       path,
		}
		err := who.Save()
		if err != nil {
			log.Error(err)
		}
	}()
	var (
		static  string
		metrics string
		who     string
		check   string
	)
	if rootPath == "/" {
		static = "/static"
		metrics = "/metrics"
		who = "/who"
		check = "/check"
	} else {
		static = fmt.Sprintf("%s/static", rootPath)
		metrics = fmt.Sprintf("%s/metrics", rootPath)
		who = fmt.Sprintf("%s/who", rootPath)
		check = fmt.Sprintf("%s/check", rootPath)
	}
	c.HTML(http.StatusOK, "admin", gin.H{
		"StaticPath": static,
		"Metrics":    metrics,
		"Who":        who,
		"Check":      check,
	})
}

func Ws(c *gin.Context) {
	defer func() {
		if httpXterm.Options.Audit {
			var path string
			if rootPath == "/" {
				path = "/ws"
			} else {
				path = fmt.Sprintf("%s/ws", rootPath)
			}
			who := &Whos{
				Remoteaddr: c.Request.RemoteAddr,
				Path:       path,
			}
			err := who.Save()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	conns := atomic.AddInt64(httpXterm.Connections, 1)
	connects.Set(float64(conns))
	if httpXterm.Options.MaxConnections != 0 {
		if conns > httpXterm.Options.MaxConnections {
			log.WithFields(log.Fields{
				"tty.go": "147",
			}).Printf("Max Connected: %d", httpXterm.Options.MaxConnections)
			atomic.AddInt64(httpXterm.Connections, -1)
			return
		}
	}

	// 升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	var cmd *exec.Cmd
	if len(httpXterm.Cmds) == 0 {
		cmd = exec.Command("bash")
	} else if len(httpXterm.Cmds) == 1 {
		cmd = exec.Command(httpXterm.Cmds[0])
	} else if len(httpXterm.Cmds) > 1 {
		cmd = exec.Command(httpXterm.Cmds[0], httpXterm.Cmds[1:]...)
	}

	//这里得到标准输出和标准错误输出的两个管道，此处获取了错误处理
	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Errorf("ptmx[52] %s", err.Error())
		return
	}

	if httpXterm.Options.MaxConnections != 0 {
		log.WithField("tty.go", "169").Printf("Command is running for client %s with PID %d (args=%q), connections: %d/%d",
			c.Request.RemoteAddr, cmd.Process.Pid, httpXterm.Cmds, conns, httpXterm.Options.MaxConnections)
	} else {
		log.WithField("tty.go", "172").Printf("Command is running for client %s with PID %d (args=%q), connections: %d",
			c.Request.RemoteAddr, cmd.Process.Pid, httpXterm.Cmds, conns)
	}

	httpXterm.Server.StartGo()

	context := &ClientContext{
		Xtermjs: httpXterm,
		Request: c.Request,
		WsConn:  ws,
		Cmd:     cmd,
		Pty:     ptmx,
		// Cache:      bytes.NewBuffer([]byte("")),
		// CacheMutex: &sync.Mutex{},
		WriteMutex: &sync.Mutex{},
	}

	context.HandleClient()
	httpXterm.Server.WaitGo()
}

func Index(c *gin.Context) {
	defer func() {
		if httpXterm.Options.Audit {
			who := &Whos{
				Remoteaddr: c.Request.RemoteAddr,
				Path:       rootPath,
			}
			err := who.Save()
			if err != nil {
				log.Error(err)
			}
		}
	}()

	// 修改Ws 全局初始化命令
	cmds, status := c.GetQueryArray("cmds")
	if status {
		if len(cmds) > 0 {
			httpXterm.Cmds = cmds
		}
	}

	var protocol, httproto string
	if httpXterm.Options.EnableTLS && IsPathExists(httpXterm.Options.CrtPath) && IsPathExists(httpXterm.Options.KeyPath) {
		protocol = "wss"
		httproto = "https"
	} else {
		protocol = "ws"
		httproto = "http"
	}
	newXsrf := GetRandomSalt()
	log.WithField("tty.go", "212").Debugf("%s xsrftoken %s", c.Request.RemoteAddr, newXsrf)
	if !httpXterm.Options.Xsrf {
		httpXterm.XsrfToken.Store(fmt.Sprintf("%s%s", newXsrf, strings.Split(c.Request.RemoteAddr, ":")[0]), time.Now().String())
	}

	var (
		path   string
		wspath string
		admin  string
	)
	if rootPath == "/" {
		path = "/static"
		wspath = "/ws"
		admin = "/admin"
	} else {
		path = fmt.Sprintf("%s/static", rootPath)
		wspath = fmt.Sprintf("%s/ws", rootPath)
		admin = fmt.Sprintf("%s/admin", rootPath)
	}

	c.HTML(http.StatusOK, "index", gin.H{
		"host":       c.Request.RemoteAddr,
		"Reconnect":  httpXterm.Options.IsReconnect,
		"Debug":      httpXterm.Options.IsDebug,
		"Write":      httpXterm.Options.PermitWrite,
		"MaxC":       httpXterm.Options.MaxConnections,
		"Conn":       *httpXterm.Connections + 1,
		"Cmd":        strings.Join(httpXterm.Cmds, " "),
		"Xsrf":       newXsrf,
		"Protocol":   protocol,
		"Httproto":   httproto,
		"isAduit":    httpXterm.Options.Audit,
		"StaticPath": path,
		"WsPath":     wspath,
		"Admin":      admin,
	})
}
