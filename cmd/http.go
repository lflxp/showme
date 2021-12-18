// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/http"
	"os"
	"os/signal"

	"gitee.com/lflxp/proxy"
	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/spf13/cobra"
)

var target string
var local string

// httpCmd represents the httpreverse command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http反向代理",
	Long: `1、http反向代理 类似nginx
2、自动添加https tls支持
3、技术探索和学习`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		funcs := proxy.NewHttpProxyByGinCustom(target, nil)
		r.Any("/*action", funcs)
		server := &http.Server{
			Addr:    local,
			Handler: r,
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

		log.Infof("Listening and serving HTTPS on %s", local)
		if err := server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Warn("Server closed under request")
			} else {
				log.Fatal("Server closed unexpect")
			}
		}

		log.Warn("Server exiting")
	},
}

func init() {
	proxyCmd.AddCommand(httpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	httpCmd.Flags().StringVarP(&target, "target", "t", "https://www.baidu.com", "七层反向代理的http地址")
	httpCmd.Flags().StringVarP(&local, "local", "L", "0.0.0.0:8888", "本地代理服务器地址")
}
