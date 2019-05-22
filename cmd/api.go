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
	"github.com/lflxp/showme/api"
	"github.com/spf13/cobra"
)

var (
	swagger                     bool
	ip, port, dbName, defaultDb string
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "OPS REST API INTERFACE",
	Long: `* 基于Bbolt的Rest CRUD Api
* 本地主机rest api性能监控
* 本地主机prometheus性能监控`,
	Run: func(cmd *cobra.Command, args []string) {
		api.Api(swagger, ip, port, dbName, defaultDb)
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	apiCmd.Flags().BoolVarP(&swagger, "swagger", "s", false, "是否开启Swagger Http服务")
	apiCmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "绑定服务IP")
	apiCmd.Flags().StringVarP(&port, "port", "p", "8080", "绑定服务端口")
	apiCmd.Flags().StringVarP(&dbName, "dbname", "n", "bbolt.db", "数据库物理文件名")
	apiCmd.Flags().StringVarP(&defaultDb, "defaultdb", "d", "test", "数据库名")
}
