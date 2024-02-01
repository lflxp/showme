# 介绍

![first](assets/first.png)
![second](assets/second.png)
![three](assets/three.png)

`本项目的目的是利用django的admin设计模式 + gin + golang实现应用的快速搭建和使用，免去搭建底层框架的困扰，自动生成CRUD的界面，节能增效`

Django 是一款流行的 Python Web 框架，提供了丰富的功能和灵活的配置选项，使得开发者可以快速地构建高质量、可扩展的 Web 应用程序。而 Golang 版本的 Django 框架 Django-Golang，则是将 Django 框架用 Golang 语言重新实现的一个版本，具有更高的运行时性能和更好的并发支持。

Django-Golang 的特点：

1. 高性能：Django-Golang 使用 Golang 语言编写，具有更高的运行时性能，可以快速地处理请求和响应。

2. 并发支持：Django-Golang 支持并发编程，可以轻松地编写多线程程序，提高程序的并发性能。

3. 灵活的配置：Django-Golang 提供了灵活的配置选项，开发者可以根据自己的需求进行自定义。

4. 丰富的功能：Django-Golang 继承了 Django 框架强大的功能，包括 ORM、模板引擎、表单、管理员等功能。

安装 Django-Golang：

1. 安装 Golang：可以从官方网站(https://golang.org/dl/)下载适合自己操作系统的安装包，然后按照官方文档的指引进行安装。

2. 安装 Django-Golang：在命令行中输入以下命令即可：
```
go get https://github.com/lflxp/djangolangexamples
```
这将下载 Django-Golang 并将其安装到你的 Golang 环境中。

使用 Django-Golang：

1. djangolangexamples 项目包含一个Demo App，里面包含一个一对多的模型和完整的Swagger API：在命令行中输入以下命令即可：
```
cd djangolangexamples
go mod tidy
go run main.go
```
这将创建一个新的 Django-Golang 项目，其中 `myproject` 是你的项目名称。

2. 设置 Django-Golang 项目配置：`main.go`中配置如下：
```
	// 是否开启https访问
	isHttps bool = true
	// 设置服务host
	Host string = "0.0.0.0"
	// 设置服务端口
	Port string = "8000"
	// 是否开启自动打开浏览器
	OpenBrowser bool = true
```
这将启动 Django-Golang 项目，并在默认的端口 8000 上运行。

3. 创建 Django-Golang 接口注册和其它配置：在 `main.go` 文件中，可以定义 Django-Golang 配置，例如：
```
	// 设置跨域设置
	r.Use(middlewares.Cors())
	// 设置恢复策略
	r.Use(gin.Recovery())
	// 设置日志
	r.Use(gin.Logger())
	// 设置健康检查接口
	r.GET("/health", middlewares.RegisterHealthMiddleware)
	// 设置swagger
	middlewares.RegisterSwaggerMiddleware(r)
	// 注册demo和vpn接口
	demo.RegisterDemo(r)
	demo.RegisterVpn(r)
	// 设置自动跳转
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/admin/index")
	})

	// 注册admin接口
	djangolang.RegisterControllerAdmin(r)
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// 配置gin优雅启动
	var server *http.Server
```

4. Admin视图适配，表结构字段解析

以下是djangolangexamples包含的两个数据结构： Demo和Vpn

```go
type Demo struct {
	Id         int64     `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Country    string    `json:"country" xorm:"varchar(255) not null" search:"true"`
	Zoom       string    `json:"zoom" xorm:"varchar(255) not null"`
	Company    string    `json:"company" xorm:"varchar(255) not null"`
	Items      string    `json:"items" xorm:"varchar(255) not null"`
	Production string    `json:"production" xorm:"varchar(255) not null"`
	Count      string    `json:"count" xorm:"varchar(255) not null"`
	Serial     string    `json:"serial" xorm:"varchar(255) not null" search:"true"`
	Extend     string    `json:"extend" xorm:"varchar(255) not null"`
	Files      string    `xorm:"file" name:"file" verbose_name:"上传文件" colType:"file"`
	File2      string    `xorm:"file2" name:"file2" verbose_name:"上传文件2" colType:"file"`
	Type       string    `xorm:"type" name:"type" verbose_name:"类型" search:"false" colType:"textarea"`
	Detail     string    `xorm:"detail" name:"detail" verbose_name:"VPN信息" list:"false" search:"false" o2m:"vpn|id,vpn" colType:"o2m"`
	Times      time.Time `xorm:"times" name:"times" verbose_name:"时间" colType:"time" list:"true" search:"true"`
}

type Vpn struct {
	Id   int64  `xorm:"id notnull unique pk autoincr" name:"id"`
	Vpn  string `xorm:"vpn" name:"vpn" verbose_name:"Vpn字段测试" list:"true" search:"true"`
	Name string `xorm:"name" name:"name" verbose_name:"姓名" list:"true" search:"false"`
	Ip   string `xorm:"ip" name:"ip" verbose_name:"ip信息" list:"true" search:"false"`
}

```

其中注释字段解析字段如下：

|  字段名   | 字段描述 | 类型  | 是否必须 | 显示效果 | 备注 |
|  ----  | ----  | ---- | ---- | ---- | ---- |
| xorm  | 数据库字段 | string | 是 | 无 | xorm框架定义数据库字段，`id notnull unique pk autoincr` 表示id字段 不为空 唯一性 主键 字增字段 |
| json | json字段显示 | string | 否 | 无 | - |
| name  | 显示字段 | string | 否 | 无 | admin自动框架表单字段 |
| verbose_name | 表单显示名称 | string | 否 | 无 | form表单显示名称 |
| search  | 是否支持搜索 | bool | 否 | 无 | 设置是否是table页面搜索框支持字段 |
| colType  | 字段类型 | string | 无 | 表单字段类型，有：`textarea` `file` `o2m` `int,int16,int64` `string` `text` `select` `radio` `multiselect` `time` `o2o` `m2m` `password` |
| list  | 是否表格显示 | bool | 否 | 无 | 表格字段是否显示 |
| o2m  | 一对多关系设置 | string | 否 | 无 | 设置一对多的表及pk外键，如： vpn|id,vpn |

Django-Golang 的配置选项非常丰富，可以根据自己的需求进行自定义。此外，Django-Golang 还提供了丰富的第三方库支持，可以方便地集成第三方库到项目中，例如：数据库、缓存、队列等。

总结：Django-Golang 是一个将 Django 框架用 Golang 语言重新实现的高性能 Web 框架，具有更高的运行时性能和更好的并发支持，非常适合开发高质量、可扩展的 Web 应用程序。如果你想要使用 Golang 构建 Web 应用程序，Django-Golang 是一个不错的选择。

# Install 

## 快速安装

> cd examples/sample && go run main.go

## 全功能

> cd examples/allinone && go run example.go

## 完整演示项目

参考 [djangolangexamples](https://github.com/lflxp/djangolangexamples)

# 功能详情

首先使用github.com/lflxp/djangolang很简单，下面是一个Demo

```go
package main

import (
	"github.com/lflxp/djangolang"
	"time"

	"github.com/gin-gonic/gin"
)

type Demotest2 struct {
	Id         int64     `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Country    string    `json:"country" xorm:"varchar(255) not null" search:"true"`
	Zoom       string    `json:"zoom" xorm:"varchar(255) not null"`
	Company    string    `json:"company" xorm:"varchar(255) not null"`
	Items      string    `json:"items" xorm:"varchar(255) not null"`
	Production string    `json:"production" xorm:"varchar(255) not null"`
	Count      string    `json:"count" xorm:"varchar(255) not null"`
	Serial     string    `json:"serial" xorm:"varchar(255) not null" search:"true"`
	Extend     string    `json:"extend" xorm:"varchar(255) not null"`
	Files      string    `xorm:"file" name:"file" verbose_name:"上传文件" colType:"file"`
	Times      time.Time `xorm:"times" name:"times" verbose_name:"时间" colType:"time" list:"true" search:"true"`
}

func init() {
	djangolang.RegisterAdmin(new(Demotest2))
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/admin/index")
	})

	djangolang.RegisterControllerAdmin(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```

高级功能展示

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/lflxp/djangolang"
	"github.com/lflxp/djangolang/middlewares"
	ctls "github.com/lflxp/djangolang/tls"
	"github.com/lflxp/djangolang/utils"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/skratchdot/open-golang/open"
)

var GinEngine *gin.Engine

type Args struct {
	Host        string
	Port        string
	IsHttps     bool
	OpenBrowser bool
	Auth        struct {
		Url         string
		IdentityKey string
		Dev         bool
	}
}

func Run(args *Args) {
	// gin.SetMode(gin.ReleaseMode)
	GinEngine = gin.Default()

	// 注册路由

	GinEngine.Use(gin.Logger())
	GinEngine.Use(gin.Recovery())
	GinEngine.Use(middlewares.Cors())
	GinEngine.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPathsRegexs([]string{".*"})))
	// GinEngine.Use(middlewares.NoRouteHandler)
	GinEngine.GET("/health", middlewares.RegisterHealthMiddleware)
	GinEngine.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/admin/index")
	})
	Registertest(GinEngine)
	djangolang.RegisterControllerAdmin(GinEngine)
	slog.Info("ip %s port %s", args.Host, args.Port)

	if args.Host == "" {
		// instance.Fatal("Check Host or Port config already!!!")
		args.Host = "0.0.0.0"
	}

	if args.Port == "" {
		args.Port = "8002"
	}

	var server *http.Server
	if args.IsHttps {
		err := ctls.Refresh()
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		caCeretPath := "ca.crt"

		caCrt, err := os.ReadFile(caCeretPath)
		if err != nil {
			panic(err)
		}

		pool.AppendCertsFromPEM(caCrt)

		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", args.Host, args.Port),
			Handler: GinEngine,
			TLSConfig: &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequestClientCert,
			},
		}
	} else {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", args.Host, args.Port),
			Handler: GinEngine,
		}

	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		slog.Warn("receive interrupt signal")
		if err := server.Close(); err != nil {
			slog.Error("Server Close:", err)
		}
	}()

	var openUrl string
	for index, ip := range utils.GetIPs() {
		if args.IsHttps {
			slog.Info("Listening and serving HTTPS on https://%s:%s", ip, args.Port)
		} else {
			slog.Info("Listening and serving HTTPS on http://%s:%s", ip, args.Port)
		}

		if index == 0 {
			openUrl = fmt.Sprintf("%s:%s", ip, args.Port)
		}
	}
	if args.IsHttps {
		if args.OpenBrowser {
			open.Start(fmt.Sprintf("https://%s", openUrl))
		}
		if err := server.ListenAndServeTLS("ca.crt", "ca.key"); err != nil {
			if err == http.ErrServerClosed {
				slog.Warn("Server closed under request")
			} else {
				slog.Error("Server closed unexpect %s", err.Error())
			}
		}
	} else {
		if args.OpenBrowser {
			open.Start(fmt.Sprintf("http://%s", openUrl))
		}
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Warn("Server closed under request")
			} else {
				slog.Error("Server closed unexpect %s", err.Error())
			}
		}
	}

	slog.Warn("Server exiting")
}

func main() {
	args := Args{
		IsHttps:     true,
		OpenBrowser: true,
		Port:        "9000",
	}
	Run(&args)
}
```

