package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jroimartin/gocui"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Upload(c *gin.Context) {
	// 多文件
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {

		// 上传文件到指定的路径
		c.SaveUploadedFile(file, fmt.Sprintf("./%s", file.Filename))
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func Video(c *gin.Context) {
	data, _ := GetAllFiles(".", types)
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
		pagestring = append(pagestring, fmt.Sprintf("/video?page=%d", i))
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
}

func DecorderHandler(h http.Handler, g *gocui.Gui) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, _ := g.View("history")
		fmt.Fprintln(v, fmt.Sprintf("%s - %s - %s - http://%s%s", time.Now().Format("2006-01-02 15:04:05"), r.RemoteAddr, r.Method, r.Host, r.RequestURI))
		h.ServeHTTP(w, r)
	})
}

// 跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("%s - %s - %s - http://%s%s\n", time.Now().Format("2006-01-02 15:04:05"), c.Request.RemoteAddr, c.Request.Method, c.Request.Host, c.Request.RequestURI)

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
		c.Header("Cache-Control", "no-store")
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

func staticServer(ctx context.Context, port string) {
	r := gin.Default()
	r.Use(Cors())
	r.Use(gin.Recovery())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}
	r.StaticFS("/", http.Dir(path))

	go func(ctx context.Context) {
		<-ctx.Done()
		log.Println("退出文件服务器 ...")

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("退出文件服务器 Server exiting")
	}(ctx)

	log.Printf("开启Static Server 0.0.0.0:%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}
}

func serverGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(Cors())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	}))
	// 使用 Recovery 中间件
	router.Use(gin.Recovery())
	router.Use(ginprom.PromMiddleware(nil))

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

	// fs := assetfs.AssetFS{
	// 	Asset:    Asset,
	// 	AssetDir: AssetDir,
	// }

	if !raw {
		// 加载html template模板
		router.HTMLRender = HtmlTemp
		router.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

		// automatically add routers for net/http/pprof
		// e.g. /debug/pprof, /debug/pprof/heap, etc.
		ginpprof.Wrapper(router)

		// router.StaticFS("/dist", &fs)
		// 前端文件
		// https://blog.csdn.net/LeoForBest/article/details/121041743
		router.Any("/static/*filepath", func(c *gin.Context) {
			staticServer := http.FileServer(http.FS(Static))
			staticServer.ServeHTTP(c.Writer, c.Request)
		})
		// router.StaticFS("/dist", http.FS(Static))

		// curl -X POST http://127.0.0.1:9090/upload -F "file=@/Users/lxp/123.mp4" -H "Content-Type:multipart/form-data"
		router.POST("/upload", Upload)

		// 美化static
		// router.GET("/x", func(c *gin.Context) {
		// 	c.Redirect(302, "/index.html")
		// })

		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "element", gin.H{
				"isVideo":    isvideo,
				"staticPort": staticPort,
			})
		})

		// swagger
		// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if isvideo {
			router.GET("/video", Video)
		}
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go staticServer(ctx, staticPort)

	go func() {
		<-closeChan
		log.Println("退出页面服务器 ...")

		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		log.Println("退出页面服务器 Server exiting")
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}
}
