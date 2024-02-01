/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/lflxp/lflxp-monitor/pkg"
	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "实时性能监控工具",
	Long: `监控指标包括：CPU、内存、磁盘、网络、时间
默认刷新时间1秒钟`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) > 2 {
			pkg.Run(strings.Join(os.Args[2:], " "))
		} else {
			pkg.Run("--lazy")
		}

	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// monitorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	monitorCmd.Flags().Bool("lazy", false, "Print Info  (include -t,-l,-c,-s,-com,-hit).")
	monitorCmd.Flags().BoolP("load", "l", false, "load info")
	monitorCmd.Flags().BoolP("cpu", "c", false, "cpu info")
	monitorCmd.Flags().BoolP("swap", "s", false, "swap info")
	monitorCmd.Flags().BoolP("net", "n", false, "基础网络信息")
	monitorCmd.Flags().BoolP("netdetail", "N", false, "详细网络信息")
	monitorCmd.Flags().BoolP("disk", "d", false, "disk info")
	monitorCmd.Flags().BoolP("time", "t", false, "打印当前时间")
}
