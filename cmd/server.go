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
	"github.com/lflxp/showme/pkg/proxy/server"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"github.com/spf13/cobra"
)

var (
	configFile  string
	cmdConfig   ss.Config
	printVer    bool
	core        int
	sanitizeIps bool
	udp         bool
	managerAddr string
	debug       ss.DebugLog
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "ss server",
	Long:  `usage: ./showme proxy ss server -c config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		server.RunServer(printVer, sanitizeIps, udp, configFile, managerAddr, cmdConfig, debug, core)
	},
}

func init() {
	ssCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().BoolVarP(&printVer, "version", "v", false, "print version")
	serverCmd.Flags().StringVarP(&configFile, "config", "c", "config.json", "specify config file")
	serverCmd.Flags().StringVarP(&cmdConfig.Password, "pwd", "k", "", "password")
	serverCmd.Flags().IntVarP(&cmdConfig.ServerPort, "port", "p", 0, "server port")
	serverCmd.Flags().IntVarP(&cmdConfig.Timeout, "timeout", "t", 300, "timeout in seconds")
	serverCmd.Flags().StringVarP(&cmdConfig.Method, "method", "m", "", "encryption method, default: aes-256-cfb")
	serverCmd.Flags().IntVarP(&core, "core", "N", 0, "maximum number of CPU cores to use, default is determinied by Go runtime")
	serverCmd.Flags().BoolVarP((*bool)(&debug), "debug", "d", false, "print debug message")
	serverCmd.Flags().BoolVarP((*bool)(&sanitizeIps), "sanitizeIps", "A", false, "anonymize client ip addresses in all output")
	serverCmd.Flags().BoolVarP(&udp, "udp", "u", false, "UDP Relay")
	serverCmd.Flags().StringVarP(&managerAddr, "manager-address", "M", "", "shadowsocks manager listening address")
}
