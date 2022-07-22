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
	"github.com/lflxp/lflxp-static/pkg"
	staticsync "github.com/lflxp/lflxp-static/pkg/sync"
	"github.com/spf13/cobra"
)

var (
	portHttpStatic                string
	pathHttpStatic                string
	isVideo                       bool
	raw                           bool
	staticPort                    string
	pagesize                      int
	types                         string
	debug2                        bool
	src, dest, clean, cleanString string
)

// staticCmd represents the static command
var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "本地静态文件服务器",
	Long:  `通过本地http服务进行简单的文件传输和文件展示`,
	Run: func(cmd *cobra.Command, args []string) {
		if src != "" && dest != "" {
			err := staticsync.LocalDirSync(src, dest, debug2)
			if err != nil {
				panic(err)
			}
		} else if clean != "" && cleanString != "" {
			err := staticsync.Clean(clean, cleanString)
			if err != nil {
				panic(err)
			}
		} else {
			api := pkg.Apis{
				Port:       portHttpStatic,
				Path:       pathHttpStatic,
				Types:      types,
				IsVideo:    isVideo,
				PageSize:   pagesize,
				Raw:        raw,
				StaticPort: staticPort,
			}

			err := api.Check()
			if err != nil {
				panic(err)
			}

			err = api.Execute()
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(staticCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// staticCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// staticCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	staticCmd.Flags().StringVarP(&types, "types", "t", ".avi,.wma,.rmvb,.rm,.mp4,.mov,.3gp,.mpeg,.mpg,.mpe,.m4v,.mkv,.flv,.vob,.wmv,.asf,.asx", "过滤视频类型，多个用逗号隔开")
	staticCmd.Flags().StringVarP(&portHttpStatic, "port", "p", "9090", "服务端口")
	staticCmd.Flags().StringVarP(&staticPort, "staticPort", "P", "9091", "文件服务端口")
	staticCmd.Flags().StringVarP(&pathHttpStatic, "path", "f", "./", "加载目录")
	staticCmd.Flags().BoolVarP(&isVideo, "video", "v", false, "是否切换为视频模式")
	staticCmd.Flags().BoolVarP(&raw, "raw", "r", false, "是否切换为无html页面状态")
	staticCmd.Flags().BoolVarP(&debug2, "debug", "d", false, "是否开启断点续传debug日志")
	staticCmd.Flags().IntVarP(&pagesize, "pagesize", "c", 20, "每页显示视频数")
	staticCmd.Flags().StringVarP(&src, "src", "S", "", "复制文件原文件或文件夹")
	staticCmd.Flags().StringVarP(&dest, "dest", "D", "", "复制文件目标文件或文件夹")
	staticCmd.Flags().StringVarP(&clean, "clean", "C", "", "删除文件或文件夹,如：/tmp")
	staticCmd.Flags().StringVarP(&cleanString, "contains", "T", "", "删除文件名包含的内容，如：.temp")
}
