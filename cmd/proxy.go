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
	"log"

	"github.com/lflxp/showme/proxy"
	"github.com/spf13/cobra"
)

var (
	localPort string
	targetUrl string
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "代理服务器",
	Long: `* http正向代理
* http 反向代理
* mysql tcp代理（负载均衡、读写分离、分布式调度）
* socket5 代理
* ss fq代理`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("本地服务 0.0.0.0:%s To %s\n", localPort, targetUrl)
		proxy.Run(localPort, targetUrl)
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	proxyCmd.Flags().StringVarP(&localPort, "port", "p", "9090", "本地代理服务器端口")
	proxyCmd.Flags().StringVarP(&targetUrl, "target", "t", "http://127.0.0.1:8888", "需要代理的服务器")
}
