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
	"github.com/lflxp/showme/executors/httpstatic"
	"github.com/spf13/cobra"
)

var (
	portHttpStatic string
	pathHttpStatic string
)

// httpstaticCmd represents the httpstatic command
var httpstaticCmd = &cobra.Command{
	Use:   "httpstatic",
	Short: "本地静态文件服务器",
	Long:  `通过本地http服务进行简单的文件传输和文件展示`,
	Run: func(cmd *cobra.Command, args []string) {
		httpstatic.HttpStaticServeForCorba(portHttpStatic, pathHttpStatic)
	},
}

func init() {
	rootCmd.AddCommand(httpstaticCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpstaticCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpstaticCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	httpstaticCmd.Flags().StringVarP(&portHttpStatic, "port", "p", "9090", "服务端口")
	httpstaticCmd.Flags().StringVarP(&pathHttpStatic, "path", "f", "./", "加载目录")
}
