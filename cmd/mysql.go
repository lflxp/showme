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
	"github.com/lflxp/showme/proxy"
	"github.com/spf13/cobra"
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql proxy",
	Long: `1、mysql tcp端口转发
2、灵感：负载均衡
3、灵感：读写分类
4、灵感：mysql proxy代理
5、灵感：分布式调度`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy.RunProxy("mysql")
	},
}

func init() {
	proxyCmd.AddCommand(mysqlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mysqlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mysqlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
