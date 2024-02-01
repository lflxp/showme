/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"strings"

	sm "github.com/lflxp/smkubectl/cmd"
	"github.com/spf13/cobra"
)

var smartDebugLevel bool

// smartCmd represents the smart command
var smartCmd = &cobra.Command{
	Use:   "smart",
	Short: "智能命令数据补全工具",
	Long:  `需配合completion，功能类似：kubectl + fzf + zsh-completion`,
	Run: func(cmd *cobra.Command, args []string) {
		// 动态调整日志输出级别
		if smartDebugLevel {
			lvl.Set(slog.LevelDebug)
		}
		slog.Debug("smart cmd", slog.Any("args", args))
		sm.ParseCmd(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(smartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// smartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// smartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	smartCmd.Flags().BoolVarP(&smartDebugLevel, "debug", "d", false, "Log Level")
}
