package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chenjiandongx/ginprom"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	_ "github.com/lflxp/lflxp-api/pkg/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	HtmlTemp multitemplate.Render
)

// 跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
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

func Api(data *Apis) {
	boltDB = NewBolt()
	defer boltDB.Close()
	if data.Stats {
		// Statistics
		go func() {
			// Grab the initial stats.
			prev := boltDB.Bolt.Stats()

			for {
				// Wait for 10s.
				time.Sleep(10 * time.Second)

				// Grab the current stats and diff them.
				stats := boltDB.Bolt.Stats()
				diff := stats.Sub(&prev)

				// Encode stats to JSON and print to STDERR.
				json.NewEncoder(os.Stderr).Encode(diff)

				// Save stats for the next loop.
				prev = stats
			}
		}()
	}
	log.SetLevel(log.InfoLevel)
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	// 跨域问题
	r.Use(Cors())
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// prometheus
	r.Use(ginprom.PromMiddleware(nil))
	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	fs := assetfs.AssetFS{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	r.StaticFS("/static", &fs)

	htmlTemplate, err := Asset("main.html")
	if err != nil {
		panic(err)
	}

	HtmlTemp = multitemplate.New()
	t, _ := template.New("index").Parse(string(htmlTemplate))
	HtmlTemp.Add("index", t)
	r.HTMLRender = HtmlTemp

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index", nil)
	})

	// // 加载全局递归所有目标
	// router.LoadHTMLGlob("templates/**/*")
	// // router.LoadHTMLGlob("templates/*")

	RegisterAPI(r)
	RegisterOrmAPI(r)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", data.Host, data.Port),
		Handler: r,
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

	log.Infof("Listening and serving HTTPS on %s:%s", data.Host, data.Port)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect", err.Error())
		}
	}

	log.Println("Server exiting")
}
