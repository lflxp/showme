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

// httpreverseCmd represents the httpreverse command
var httpreverseCmd = &cobra.Command{
	Use:   "httpreverse",
	Short: "http反向代理",
	Long: `1、http反向代理 类似nginx
2、自动添加https tls支持
3、技术探索和学习`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy.RunProxy("httprp")
	},
}

func init() {
	proxyCmd.AddCommand(httpreverseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpreverseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// httpreverseCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
