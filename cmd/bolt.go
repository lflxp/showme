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
	"github.com/lflxp/showme/boltapi"
	"github.com/spf13/cobra"
)

var (
	bport string
	bhost string
	stats bool
)

// boltCmd represents the bolt command
var boltCmd = &cobra.Command{
	Use:   "api",
	Short: "快速本地DB CRUD API",
	Long: `1. 基于bolt+gin+swagger+web的简洁RESTFUL API服务
2. 提供本地监控数据、外部服务数据的存储和查询功能
3. 提供针对时间范围搜索的类时序数据库功能`,
	Run: func(cmd *cobra.Command, args []string) {
		boltapi.Api(bhost, bport, stats)
	},
}

func init() {
	rootCmd.AddCommand(boltCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// boltCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// boltCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	boltCmd.Flags().StringVarP(&bhost, "host", "H", "0.0.0.0", "http bind host")
	boltCmd.Flags().StringVarP(&bport, "port", "P", "8080", "http bind port")
	boltCmd.Flags().BoolVarP(&stats, "stats", "s", false, "output db stats")
}
