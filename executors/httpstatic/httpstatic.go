package httpstatic

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/jroimartin/gocui"
	"github.com/lflxp/showme/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	types     []string
	pageSize  int
	isvideo   bool
	path      string
	port      string
	closeChan chan os.Signal
	uri       string
	initnum   int
)

const maxUploadSize = 2000 * 1024 * 2014 // 2 MB
const uploadPath = "/tmp"

func init() {
	initnum = 0
	// path = utils.GetCurrentDirectory()
	// port = "9090"
	closeChan = make(chan os.Signal)
}

func DecorderHandler(h http.Handler, g *gocui.Gui) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, _ := g.View("history")
		fmt.Fprintln(v, fmt.Sprintf("%s - %s - %s - http://%s%s", time.Now().Format("2006-01-02 15:04:05"), r.RemoteAddr, r.Method, r.Host, r.RequestURI))
		h.ServeHTTP(w, r)
	})
}

// 跨域设置
func Cors(g *gocui.Gui) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, _ := g.View("history")
		fmt.Fprintln(v, fmt.Sprintf("%s - %s - %s - http://%s%s", time.Now().Format("2006-01-02 15:04:05"), c.Request.RemoteAddr, c.Request.Method, c.Request.Host, c.Request.RequestURI))

		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		// headerStr := strings.Join(headerKeys, ", ")
		// if headerStr != "" {
		// 	headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		// } else {
		// 	headerStr = "access-control-allow-origin, access-control-allow-headers"
		// }
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		//if method == "OPTIONS" {
		//    c.JSON(http.StatusOK, "Options Request!")
		//}
		if method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		// 处理请求
		c.Next() //  处理请求
	}
}

func serverGin(g *gocui.Gui) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors(g))
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	}))
	// 使用 Recovery 中间件
	router.Use(gin.Recovery())
	router.Use(ginprom.PromMiddleware(nil))
	router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	router.GET("/", func(c *gin.Context) {
		if isvideo {
			c.Redirect(http.StatusMovedPermanently, "/video")
		} else {
			c.Redirect(http.StatusMovedPermanently, "/static")
		}
	})
	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	ginpprof.Wrapper(router)

	// LoggerWithFormatter 中间件会将日志写入 gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
		// 你的自定义格式
		// return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		//         param.ClientIP,
		//         param.TimeStamp.Format(time.RFC1123),
		//         param.Method,
		//         param.Path,
		//         param.Request.Proto,
		//         param.StatusCode,
		//         param.Latency,
		//         param.Request.UserAgent(),
		//         param.ErrorMessage,
		// )
	}))

	router.StaticFS("/static", http.Dir(path))
	// curl -X POST http://127.0.0.1:9090/upload -F "file=@/Users/lxp/123.mp4" -H "Content-Type:multipart/form-data"
	router.POST("/upload", func(c *gin.Context) {
		// 多文件
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for _, file := range files {

			// 上传文件到指定的路径
			c.SaveUploadedFile(file, fmt.Sprintf("./%s", file.Filename))
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	if isvideo {
		data, _ := utils.GetAllFiles(".", types)
		indexhtml := multitemplate.New()
		htmlTemplate := `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style type="text/css">
	li {
		float: left;
		display: inline;
		list-style: none;
		border: 0px solid #161EE8FF;
		text-align: center;
		line-height: 200px;
		padding: 9px;
		margin-top: 10px;
		height: 220px;
		width: 220px;
		margin: 10 10px;
	}
</style>
</head>
<body>
<div id="nav">
  <ul>
	{{ range $src := .data }}
	<li><video src="{{ $src }}" controls="controls" height="200" width="200" preload="metadata" loop="loop">{{ $src }}</video></li>
	{{ end }}
  </ul>
  {{ .page }}
  <ul>
	{{ range $n,$page := .pages }}
		<li><a href="{{$page}}" target="_self">{{ $n }}</a></li>
	{{ end }}
  </ul>
</div>
</body>
</html>`
		t, _ := template.New("index").Parse(htmlTemplate)
		indexhtml.Add("index", t)
		router.HTMLRender = indexhtml
		router.GET("/video", func(c *gin.Context) {
			var currentPage int
			page := c.DefaultQuery("page", "")
			if page == "" {
				currentPage = 1
			} else {
				var err error
				currentPage, err = strconv.Atoi(page)
				if err != nil {
					c.String(http.StatusBadRequest, err.Error())
					return
				}
			}

			var pages int
			if len(data)%pageSize > 0 {
				pages = len(data)/pageSize + 1
			} else {
				pages = len(data) / pageSize
			}
			pagestring := []string{}
			for i := 1; i <= pages; i++ {
				pagestring = append(pagestring, fmt.Sprintf("/?page=%d", i))
			}
			if pageSize*currentPage < len(data)-1 {
				c.HTML(http.StatusOK, "index", gin.H{"data": data[pageSize*(currentPage-1) : pageSize*currentPage], "pages": pagestring})
			} else {
				if len(data) == 0 {
					c.String(http.StatusOK, "这小子什么都没留下！")
				} else {
					c.HTML(http.StatusOK, "index", gin.H{"data": data[pageSize*(currentPage-1) : len(data)-1], "pages": pagestring})
				}
			}

		})
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	signal.Notify(closeChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		os.Interrupt,
		os.Kill,
	)

	go func() {
		<-closeChan
		log.Println("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("Server exiting")
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}
}

func HttpStaticServeForCorba(ports, paths, typesed string, isVideo bool, pagesize int) {
	// httpstatic -port 9090 -path ./
	port = ports
	path = paths
	isvideo = isVideo
	pageSize = pagesize
	types = strings.Split(typesed, ",")

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	d := time.Duration(time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				g.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func HttpStaticServe(in string) {
	// httpstatic -port 9090 -path ./
	tmp := strings.Split(in, " ")
	for n, x := range tmp {
		if x == "-port" {
			port = tmp[n+1]
		}
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(dlayout)

	d := time.Duration(time.Second)
	t := time.NewTicker(d)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-t.C:
				g.Update(func(g *gocui.Gui) error { return nil })
				// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
				// fmt.Fprintln(v, )
			}
		}
	}()

	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
	// 	log.Panicln(err)
	// }

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func keybindings(g *gocui.Gui) error {
	// 清空side缓存
	// if err := g.SetKeybinding("help", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, dquit); err != nil {
		return err
	}
	return nil
}

func dquit(g *gocui.Gui, v *gocui.View) error {
	closeChan <- syscall.SIGINT
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func dlayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("history", 0, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "访问记录"
		v.Wrap = true
		v.Frame = false
		// v.Highlight = true
		v.Autoscroll = true
		v.SelFgColor = gocui.ColorYellow
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))
	}
	ips := utils.GetIPs()
	if v, err := g.SetView("top", maxX/2-80, maxY/2, maxX/2+80, maxY/2+2*len(ips)+4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if isvideo {
			v.Title = "视频服务器地址"
		} else {
			v.Title = "静态服务器地址"
		}

		v.Wrap = true
		// v.Highlight = true
		// v.Autoscroll = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		// v.Editable = true
		// fmt.Fprintf(v, time.Now().Format("2006-01-02 15:04:05"))
		// uri = fmt.Sprintf("/a%s", time.Now().Format("150405"))

		urls := []string{}

		for _, ip := range ips {
			urls = append(urls, fmt.Sprintf("UploadURL: => %s:%s <= PATH: => /upload <= ", ip, port))
		}
		for _, ip := range ips {
			urls = append(urls, fmt.Sprintf("DownURL: => %s:%s <= PATH: => / <= ", ip, port))
		}
		dir, _ := os.Getwd()
		urls = append(urls, fmt.Sprintf("CurrentPWD: => %s <= ", dir))
		urls = append(urls, fmt.Sprintf("UploadDIR: => %s <= ", dir))
		urls = append(urls, "curl -X POST http://127.0.0.1:9090/upload -F \"file=@/Users/lxp/123.mp4\" -H \"Content-Type:multipart/form-data\"")
		fmt.Fprintln(v, strings.Join(urls, "\n"))
		go serverGin(g)

		if _, err = setCurrentViewOnTop(g, "top"); err != nil {
			return err
		}
	}
	return nil
}
