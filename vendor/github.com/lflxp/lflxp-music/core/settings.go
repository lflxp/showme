package core

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"

	ctls "github.com/lflxp/lflxp-music/core/tls"

	"github.com/lflxp/lflxp-music/core/utils"

	"github.com/lflxp/lflxp-music/core/router"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

func Run(ishttps bool) {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 注册路由

	router.PreGinServe(r)
	ip := viper.GetString("host")
	port := viper.GetString("port")
	log.Infof("ip %s port %s", ip, port)

	if ip == "" || port == "" {
		// instance.Fatal("Check Host or Port config already!!!")
		ip = "0.0.0.0"
		port = "8000"
	}

	var server *http.Server
	if ishttps {
		err := ctls.Refresh()
		if err != nil {
			panic(err)
		}

		pool := x509.NewCertPool()
		caCeretPath := "ca.crt"

		caCrt, err := ioutil.ReadFile(caCeretPath)
		if err != nil {
			panic(err)
		}

		pool.AppendCertsFromPEM(caCrt)

		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", ip, port),
			Handler: r,
			TLSConfig: &tls.Config{
				ClientCAs:  pool,
				ClientAuth: tls.RequestClientCert,
			},
		}
	} else {
		server = &http.Server{
			Addr:    fmt.Sprintf("%s:%s", ip, port),
			Handler: r,
		}

	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Warn("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	var openUrl string
	for index, ip := range utils.GetIPs() {
		if ishttps {
			log.Infof("Listening and serving HTTPS on https://%s:%s", ip, port)
		} else {
			log.Infof("Listening and serving HTTPS on http://%s:%s", ip, port)
		}

		if index == 0 {
			openUrl = fmt.Sprintf("%s:%s", ip, port)
		}
	}
	if ishttps {
		open.Start(fmt.Sprintf("https://%s", openUrl))
		if err := server.ListenAndServeTLS("ca.crt", "ca.key"); err != nil {
			if err == http.ErrServerClosed {
				log.Warn("Server closed under request")
			} else {
				log.Fatalf("Server closed unexpect %s", err.Error())
			}
		}
	} else {
		open.Start(fmt.Sprintf("http://%s", openUrl))
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Warn("Server closed under request")
			} else {
				log.Fatalf("Server closed unexpect %s", err.Error())
			}
		}
	}

	log.Warn("Server exiting")
}
