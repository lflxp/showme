/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "showme dashboard快速学习平台",
	Long:  `dashboard的目的是展示showme工具的功能，让使用者快速了解和使用`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dashboard called")
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dashboardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dashboardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
