package pkg

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/DeanThompson/ginpprof"
	"github.com/unrolled/secure"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket" // https://blog.csdn.net/u014029783/article/details/80001251 教程
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func ServeGin(data.Host,data.Port, username, password, crtpath, keypath string, cmds []string, isdebug, isReconnect, isPermitWrite, isAudit, isXsrf, isProf, enabletls bool, MaxConnections int64) {
func ServeGin(data *Tty) {
	router := gin.Default()

	// 使用 Recovery 中间件
	router.Use(gin.Recovery())

	// 添加prometheus监控
	// use prometheus metrics exporter middleware.
	//
	// ginprom.PromMiddleware() expects a ginprom.PromOpts{} poniter.
	// It was used for filtering labels with regex. `nil` will pass every requests.
	//
	// ginprom promethues-labels:
	//   `status`, `endpoint`, `method`
	//
	// for example:
	// 1). I want not to record the 404 status request. That's easy for it.
	// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexStatus: "404"})
	//
	// 2). And I wish to ignore endpoint start with `/prefix`.
	// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexEndpoint: "^/prefix"})
	router.Use(ginprom.PromMiddleware(nil))

	RegisterTty(router, data, true)

	// automatically add routers for net/http/pprof
	// e.g. /debug/pprof, /debug/pprof/heap, etc.
	if data.IsProf {
		ginpprof.Wrapper(router)
	}

	var server *http.Server

	if httpXterm.Options.EnableTLS {
		// if !IsPathExists(httpXterm.Options.CrtPath) || !IsPathExists(httpXterm.Options.KeyPath) {
		// 			err := GenerateRSAKey(1024)
		// 			if err != nil {
		// 				panic(err)
		// 			}
		// 		}

		pool := x509.NewCertPool()
		caCeretPath := httpXterm.Options.CrtPath

		caCrt, err := ioutil.ReadFile(caCeretPath)
		if err != nil {
			panic(err)
		}

		pool.AppendCertsFromPEM(caCrt)

		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", data.Host, data.Port),
			Handler: router,
			TLSConfig: &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequestClientCert,
			},
		}
	} else {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", data.Host, data.Port),
			Handler: router,
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		// log.Println("receive interrupt signal")
		// if err := server.Close(); err != nil {
		// 	log.Fatal("Server Close:", err)
		// }

		slog.Info("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Server Shutdown:", err)
		}
		slog.Info("Server exiting")
	}()

	if data.Host == "0.0.0.0" {
		ips := GetIPs()
		for _, ip := range ips {
			slog.Info("Listening and serving HTTPS on http://%s:%s", ip, data.Port)
		}
	} else {
		slog.Info("Listening and serving HTTPS on http://%s:%s", data.Host, data.Port)
	}

	if httpXterm.Options.EnableTLS {
		if IsPathExists(httpXterm.Options.CrtPath) && IsPathExists(httpXterm.Options.KeyPath) {
			if err := server.ListenAndServeTLS(httpXterm.Options.CrtPath, httpXterm.Options.KeyPath); err != nil {
				if err == http.ErrServerClosed {
					slog.Info("Server closed under request")
				} else {
					slog.Error("Server closed unexpect", err.Error())
				}
			}
		} else {
			slog.Error("EnableTLS is true,but crt or key path is not exists")
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				slog.Info("Server closed under request")
			} else {
				slog.Error("Server closed unexpect", err.Error())
			}
		}
	}

	slog.Info("Server exiting")
}

func TlsHandler(host, port string) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     fmt.Sprintf("%s:%s", host, port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
