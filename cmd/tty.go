/*
Copyright © 2020 NAME lflxp <382023823@qq.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/lflxp/lflxp-tty/pkg"
	"github.com/spf13/cobra"
)

var (
	enableTLS      bool
	crtPath        string
	keyPath        string
	isProf         bool
	isXsrf         bool
	isAudit        bool
	isPermitWrite  bool
	MaxConnections int64
	isReconnect    bool
	isDebug        bool
	username       string
	password       string
	port           string
	host           string
)

// ttyCmd represents the tty command
var ttyCmd = &cobra.Command{
	Use:   "tty",
	Short: "web terminial",
	Long: `showme tty [flags] [command] [args]
eg: showme tty -w -r showme proxy http`,
	Run: func(cmd *cobra.Command, args []string) {
		tty := pkg.Tty{
			EnableTLS:      enableTLS,
			CrtPath:        crtPath,
			KeyPath:        keyPath,
			IsProf:         isProf,
			IsXsrf:         isXsrf,
			IsAudit:        isAudit,
			IsPermitWrite:  isPermitWrite,
			MaxConnections: MaxConnections,
			IsReconnect:    isReconnect,
			IsDebug:        isDebug,
			Username:       username,
			Password:       password,
			Port:           port,
			Host:           host,
			Cmds:           args,
		}
	
		err := tty.Execute()
		if err != nil {
			panic(err)
		}
		// tty.ServeGin(host, port, username, password, crtPath, keyPath, args, isDebug, isReconnect, isPermitWrite, isAudit, isXsrf, isProf, enableTLS, MaxConnections)
	},
}

func init() {
	rootCmd.AddCommand(ttyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ttyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ttyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ttyCmd.Flags().StringVarP(&username, "username", "u", "", "BasicAuth 用户名")
	ttyCmd.Flags().StringVarP(&password, "password", "p", "", "BasicAuth 密码")
	ttyCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "http bind host")
	ttyCmd.Flags().StringVarP(&port, "port", "P", "8080", "http bind port")
	ttyCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "debug log mode")
	ttyCmd.Flags().BoolVarP(&isReconnect, "reconnect", "r", false, "是否自动重连")
	ttyCmd.Flags().BoolVarP(&isPermitWrite, "write", "w", false, "是否开启写入模式")
	ttyCmd.Flags().BoolVarP(&isAudit, "audit", "a", false, "是否开启审计")
	ttyCmd.Flags().BoolVarP(&isXsrf, "xsrf", "x", false, "是否开启xsrf,默认开启")
	ttyCmd.Flags().BoolVarP(&isProf, "prof", "f", false, "是否开启pprof性能分析")
	ttyCmd.Flags().BoolVarP(&enableTLS, "tls", "t", false, "是否开启https")
	ttyCmd.Flags().StringVarP(&crtPath, "crt", "c", "./server.crt", "*.crt文件路径")
	ttyCmd.Flags().StringVarP(&keyPath, "key", "k", "./server.key", "*.key文件路径")
	ttyCmd.Flags().Int64VarP(&MaxConnections, "maxconnect", "m", 0, "最大连接数")
}
