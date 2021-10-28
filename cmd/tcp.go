/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"gitee.com/lflxp/proxy"
	"github.com/spf13/cobra"
)

var (
	from string
	to   string
)

// tcpCmd represents the tcp command
var tcpCmd = &cobra.Command{
	Use:   "tcp",
	Short: "tcp 四层反向代理",
	Long:  `access from [0.0.0.0:9999] to [0.0.0.0:3306]`,
	Run: func(cmd *cobra.Command, args []string) {
		err := proxy.NewTCPProxy(from, to)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	proxyCmd.AddCommand(tcpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tcpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tcpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tcpCmd.Flags().StringVarP(&from, "from", "f", "0.0.0.0:6001", "本地代理服务器地址")
	tcpCmd.Flags().StringVarP(&to, "to", "t", "127.0.0.1:3306", "需要反向代理的服务器地址")
}
