/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/lflxp/showme/tty"
	"github.com/spf13/cobra"
)

var (
	isDebug  bool
	username string
	password string
	host     string
	port     string
)

// ttyCmd represents the tty command
var ttyCmd = &cobra.Command{
	Use:   "tty",
	Short: "web terminial",
	Long:  `gin web + websocket + xterm.js`,
	Run: func(cmd *cobra.Command, args []string) {
		tty.ServeGin(host, port, username, password, isDebug)
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
	ttyCmd.Flags().StringVarP(&username, "username", "u", "admin", "BasicAuth 用户名")
	ttyCmd.Flags().StringVarP(&password, "password", "p", "system", "BasicAuth 密码")
	ttyCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "http host")
	ttyCmd.Flags().StringVarP(&port, "port", "P", "8080", "http port")
	ttyCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "debug日志输出")
}
