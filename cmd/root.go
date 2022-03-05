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
	fzf "github.com/junegunn/fzf/src"
	"github.com/lflxp/showme/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugs bool
var islog bool

var version string = "0.29.1"
var revision string = "cobra-dev"

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
	* 单文件、无依赖、快速、信息丰富多样化的console terminal
	3.3 ZSH bindkey配置
	bindkey -s "^[\~" "showme^M"
bindkey -s "^[1" "showme watch^M"
bindkey -s "^[2" "showme static^M"
bindkey -s "^[3" "showme tty -w^M"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fzf.Run(fzf.ParseOptions(), version, revision)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Execute Error: ", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.showme.yaml)")
	// rootCmd.PersistentFlags().BoolVarP(&debugs, "debug", "d", false, "是否打印debug日志")
	rootCmd.PersistentFlags().BoolVarP(&islog, "log", "l", false, "是否文件输出")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// =============================fzf========================
	rootCmd.Flags().BoolVarP(&extended, "extended", "x", false, "Extended-search mode (enabled by default; +x or --no-extended to disable)")
	rootCmd.Flags().BoolVarP(&exact, "exact", "e", false, "Enable Exact-match")
	rootCmd.Flags().StringVarP(&query, "query", "q", "", "Start the finder with the given query")
	rootCmd.Flags().StringVar(&algo, "algo", "v2", "Start the finder with the given query")
	rootCmd.Flags().StringVarP(&filter, "filter", "f", "", "Filter mode. Do not start interactive finder.")
	rootCmd.Flags().StringVarP(&nth, "nth", "n", "", "Comma-separated list of field index expressions")
	rootCmd.Flags().StringVarP(&delimiter, "delimiter", "d", "AWK-style", "Field delimiter regex (default: AWK-style)")
	rootCmd.Flags().StringVar(&withnth, "with-nth", "", "Transform the presentation of each line using field index expressions")
	rootCmd.Flags().BoolVar(&literal, "literal", false, "Do not normalize latin script letters before matching")
	rootCmd.Flags().BoolVar(&nosort, "no-sort", false, "Do not sort the result")
	rootCmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "Do not normalize latin script letters before matching")
	rootCmd.Flags().BoolVar(&phony, "phony", false, "Enable Exact-match")
	rootCmd.Flags().StringVar(&tiebreak, "tiebreak", "length", "omma-separated list of sort criteria to apply 	when the scores are tied [length|begin|end|index] 	(default: length)")
	rootCmd.Flags().BoolVar(&enabled, "enabled", false, "Do not perform search")
	rootCmd.Flags().BoolVar(&disabled, "disabled", false, "Do not perform search")
	rootCmd.Flags().IntVarP(&sort, "sort", "s", 0, "Do not sort the result")
	rootCmd.Flags().BoolVar(&tac, "tac", false, "Reverse the order of the input")
	rootCmd.Flags().IntVarP(&multi, "multi", "m", 0, "Enable multi-select with tab/shift-tab")
	rootCmd.Flags().BoolVar(&ansi, "ansi", false, "Enable processing of ANSI color codes")
	rootCmd.Flags().BoolVar(&mouse, "mouse", false, "Enable mouse")
	rootCmd.Flags().BoolVar(&mouse, "no-mouse", false, "Disabled mouse")
	rootCmd.Flags().BoolVar(&black, "black", false, "black")
	rootCmd.Flags().BoolVar(&bold, "bold", false, "bold")
	rootCmd.Flags().BoolVar(&cycle, "cycle", false, "cycle")
	rootCmd.Flags().BoolVar(&keepright, "keep-right", false, "keepright")
	rootCmd.Flags().BoolVar(&hscroll, "hscroll", false, "hscroll")
	rootCmd.Flags().BoolVar(&hscroll, "no-hscroll", false, "Disable horizontal scroll")
	rootCmd.Flags().IntVar(&hscrolloff, "hscroll-off", 0, "hscrolloff")
	rootCmd.Flags().IntVar(&scrolloff, "scroll-off", 0, "scrolloff")
	rootCmd.Flags().BoolVar(&fileword, "filepath-word", false, "Make word-wise movements respect path separator")
	rootCmd.Flags().StringVar(&jumplabels, "jump-labels", "", "Label characters for jump and jump-accept")
	rootCmd.Flags().BoolVarP(&select1, "select-1", "1", false, "Automatically select the only match")
	rootCmd.Flags().BoolVarP(&exit0, "exit-0", "0", false, "Exit immediately when there's no match")
	rootCmd.Flags().BoolVar(&readzero, "read0", false, "Read input delimited by ASCII NUL characters")
	rootCmd.Flags().BoolVar(&print0, "print0", false, "Print output delimited by ASCII NUL character")
	rootCmd.Flags().StringVar(&prompt, "expect", "", "Comma-separated list of keys to complete fzf")
	rootCmd.Flags().BoolVar(&printquery, "print-query", false, "Print query as the first line")
	rootCmd.Flags().StringVar(&pointer, "pointer", "", "Pointer to the current line (default: '>'")
	rootCmd.Flags().StringVar(&marker, "marker", "", "Multi-select marker (default: '>')")
	rootCmd.Flags().BoolVar(&sync, "sync", false, "Synchronous search for multi-staged filtering")
	rootCmd.Flags().StringVar(&preview, "preview", "", "Command to preview highlighted line ({})")
	rootCmd.Flags().StringVar(&preview, "preview-window", "right:50%", "Preview window layout (default: right:50%)")
	rootCmd.Flags().StringVar(&height, "height", "40%", "Display fzf window below the cursor with the give")
	rootCmd.Flags().IntVar(&minheight, "min-height", 0, "Minimum height when --height is given in percent default: 10")
	rootCmd.Flags().BoolVar(&unicode, "unicode", false, "unicode")
	rootCmd.Flags().IntVar(&tabstop, "tabstop", 8, "Number of spaces for a tab character (default: 8)")
	rootCmd.Flags().BoolVar(&clearonexit, "clear", false, "clearonexit")
	rootCmd.Flags().BoolVar(&versions, "versions", false, "versions")
	rootCmd.Flags().StringVar(&bind, "bind", "", "Custom key bindings. Refer to the man page.")
	rootCmd.Flags().StringVar(&layout, "layout", "default", "Choose layout: [default|reverse|reverse-list]")
	rootCmd.Flags().StringVar(&border, "border", "", "Draw border around the finder")
	rootCmd.Flags().StringVar(&border, "margin", "rounded", "Screen margin (TRBL | TB,RL | T,RL,B | T,R,B,L)")
	rootCmd.Flags().StringVar(&border, "padding", "", "Padding inside border (TRBL | TB,RL | T,RL,B | T,R,B,L)")
	rootCmd.Flags().StringVar(&border, "info", "", "Finder info style [default|inline|hidden]")
	rootCmd.Flags().StringVar(&border, "prompt", ">", "Input prompt (default: '> ')")
	rootCmd.Flags().StringVar(&border, "header", "", "String to print as header")
	rootCmd.Flags().IntVar(&headerlines, "header-lines", 0, "The first N lines of the input are treated as heade")
	rootCmd.Flags().BoolVar(&headerfirst, "header-first", false, "Print header before the prompt line")
	rootCmd.Flags().StringVar(&color, "color", "", "Base scheme (dark|light|16|bw) and/or custom colors")
	rootCmd.Flags().BoolVar(&versions, "no-bold", false, "Do not use bold text")
	rootCmd.Flags().StringVar(&color, "history", "", "History file")
	rootCmd.Flags().StringVar(&color, "history-size", "1000", "Maximum number of history entries (default: 1000)")
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
