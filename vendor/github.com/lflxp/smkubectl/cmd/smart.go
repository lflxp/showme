/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

var debugLevel bool

// smartCmd represents the smart command
var smartCmd = &cobra.Command{
	Use:   "smart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 动态调整日志输出级别
		if debugLevel {
			lvl.Set(slog.LevelDebug)
		}
		slog.Debug("smart cmd", slog.Any("args", args))
		ParseCmd(strings.Join(args, " "))
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
	smartCmd.Flags().BoolVarP(&debugLevel, "debug", "d", false, "Log Level")
}
