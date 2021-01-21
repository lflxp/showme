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
	"fmt"
	"os"
	"time"

	_ "github.com/devopsxp/xp/module"
	"github.com/lflxp/showme/executors"
	"github.com/lflxp/showme/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugs bool
var islog bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "showme",
	Short: "运维快速问题排查工具兼运维自动化工具",
	Long: `1. 排查工具包括：
	* IP扫描
	* 端口扫描
	* 性能实时监控
	* 静态文件传输服务器
	* mysql性能监控及排查工具
	* 网络流量抓包工具
	* shadowsocks local & server client
2. 运维自动化工具
	* 本地API接口（远程监控、远程数据bbolt存储）
	* 运维自动化Agent（RPCX远程过程调用）
	* web terminal 快速服务器登录
3. 无参数无命令
	* 运维本地快速问题排查工具
	3.1 排查工具介绍
	* 致力于解决人肉运维中想快速定位系统性能、数据库性能、网络包、快速文件传输服务器等基础但重要的功能
	3.2 目标
	* 单文件、无依赖、快速、信息丰富多样化的console terminal`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		executors.AllInOne()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.showme.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugs, "debug", "d", false, "是否打印debug日志")
	rootCmd.PersistentFlags().BoolVarP(&islog, "log", "l", false, "是否文件输出")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true // 显示时间
	// customFormatter.ForceQuote = true // 强制格式匹配
	// customFormatter.PadLevelText = true // 显示完整日志级别

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if islog {
		file, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now().Format("20060102150405")), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			panic(err)
		}

		info, err := file.Stat()
		if err != nil {
			panic(err)
		}

		fileWriter := utils.LogFileWriter{file, info.Size()}
		log.SetOutput(fileWriter)
	} else {
		log.SetOutput(os.Stdout)
	}

	//  log format setting
	if debugs {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// 获取项目的执行路径
		home, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		// Search config in home directory with name ".devopsxp" (without extension).
		viper.AddConfigPath(home)       // 设置读取文件的路径
		viper.SetConfigName("devopsxp") // 设置读取的文件名
		viper.SetConfigType("yaml")     // 设置文件的类型
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Errorf("Using config file error: %s\n", err.Error())
	}
}
