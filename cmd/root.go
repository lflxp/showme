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
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/devopsxp/xp/module"
	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/protector"
	"github.com/lflxp/showme/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/util/homedir"
)

var cfgFile string
var debugs bool
var islog bool
var version string = "0.46"
var revision string = "devel"
var lvl slog.LevelVar

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
	* smart && fzf 快速命令补全工具
3. 无参数无命令
	* 运维本地快速问题排查工具
	3.1 排查工具介绍
	* 致力于解决人肉运维中想快速定位系统性能、数据库性能、网络包、快速文件传输服务器等基础但重要的功能
	3.2 目标
	* 单文件、无依赖、快速、信息丰富多样化的console terminal
	3.3 ZSH bindkey配置
	bindkey -s "^[\~" "showme^M"
	bindkey -s "^[1" "showme --height=40%^M"
	bindkey -s "^[2" "showme cmd^M"
	bindkey -s "^[3" "showme static^M"
	bindkey -s "^[4" "showme tty -w^M"
	bindkey -s "^[5" "showme watchw^M"
4. ~/.showme.yaml example`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 保证没有参数或者参数只有一个且为completion的时候执行cobra
	// 其余都走parseCmd
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "completion", "smart", "api", "cmd", "dashboard", "help", "k8s", "camera", "martix", "music", "playbook", "proxy", "static", "tty", "watch", "scan", "django", "metrics", "monitor", "screego":
			slog.Debug("进入cobra命令模式", slog.Any("args", os.Args[1]))
			err := rootCmd.Execute()
			if err != nil {
				slog.Error("执行cobra命令失败", slog.Any("err", err))
				os.Exit(1)
			}
		default:
			slog.Debug("进入completion模式")
			protector.Protect()
			fzf.Run(fzf.ParseOptions(), version, revision)
		}
	} else {
		slog.Debug("进入fzf模式")
		protector.Protect()
		fzf.Run(fzf.ParseOptions(), version, revision)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.showme.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugs, "debug", "G", false, "是否打印debug日志")
	rootCmd.PersistentFlags().BoolVarP(&islog, "log", "O", false, "是否文件输出")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// customFormatter := new(log.TextFormatter)
	// customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(customFormatter)
	// customFormatter.FullTimestamp = true // 显示时间
	// customFormatter.ForceQuote = true // 强制格式匹配
	// customFormatter.PadLevelText = true // 显示完整日志级别

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// 日志配置
	// lvl.Set(slog.LevelError)
	lvl.Set(slog.LevelInfo)
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     &lvl,
	}

	// slog.SetDefault(slog.New((slog.NewJSONHandler(os.Stdout, &opts))))

	if islog {
		file, err := os.OpenFile(fmt.Sprintf("%s.log", time.Now().Format("20060102150405")), os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// info, err := file.Stat()
		// if err != nil {
		// 	panic(err)
		// }

		// fileWriter := utils.LogFileWriter{file, info.Size()}
		// log.SetOutput(fileWriter)

		slog.SetDefault(slog.New((slog.NewTextHandler(file, &opts))))
	} else {
		slog.SetDefault(slog.New((slog.NewTextHandler(os.Stdout, &opts))))
		// log.SetOutput(os.Stdout)
	}

	//  log format setting
	if debugs {
		lvl.Set(slog.LevelDebug)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// 获取项目的执行路径
		// home, err := os.Getwd()
		home := homedir.HomeDir()
		target := filepath.Join(home, ".showme.yaml")
		if !utils.IsPathExists(target) {
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			_, err = file.WriteString(`account:
    admin:
        claim: '[{"id":1,"auth":"admin","type":"nav","value":"dashboard"}]'
        password: admin
admin: true
app:
    - test
global:
    Name: demo
    Pkg: demo
host: 0.0.0.0
log:
    level: info
meili:
    apikey: masterKey
    enable: false
    host: http://127.0.0.1:7700
port: 8000
snakemapper: admin_
`)
			if err != nil {
				panic(err)
			}

		}

		// Search config in home directory with name ".devopsxp" (without extension).
		viper.AddConfigPath(home)      // 设置读取文件的路径
		viper.SetConfigName(".showme") // 设置读取的文件名
		viper.SetConfigType("yaml")    // 设置文件的类型
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.ReadInConfig()
	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// fmt.Println("Using config file:", viper.ConfigFileUsed())
	// } else {
	// slog.Error("配置文件读取错误", "Using config file error: %s\n", err.Error())
	// }
}
