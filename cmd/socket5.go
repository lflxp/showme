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
	"github.com/lflxp/showme/pkg/proxy"
	"github.com/spf13/cobra"
)

// socket5Cmd represents the socket5 command
var socket5Cmd = &cobra.Command{
	Use:   "socket5",
	Short: "socket5 http代理服务器",
	Long: `1、socket5 protocol support
2、非加密
3、基于socket5 代理，local直连，无server代理`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy.RunProxy("socket5")
	},
}

func init() {
	proxyCmd.AddCommand(socket5Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// socket5Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// socket5Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
