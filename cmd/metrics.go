/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/lflxp/djangolang/middlewares"
	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "prometheus metrics",
	Long:  `程序性能监控指标`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/abc", func(c *gin.Context) {
			c.Redirect(301, "/metrics")
		})

		middlewares.RegisterPrometheusMiddlewareBasic(r, false)
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metricsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metricsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
