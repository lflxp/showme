/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	anywayPort string
)

// anywayCmd represents the anyway command
var anywayCmd = &cobra.Command{
	Use:   "anyway",
	Short: "测试校验接口",
	Long:  `任意接口请求方式或者路径均返回200｜success`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.Any("/*any", func(c *gin.Context) {
			c.JSON(200, "success")
		})

		r.Run(fmt.Sprintf("0.0.0.0:%s", anywayPort))
	},
}

func init() {
	rootCmd.AddCommand(anywayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// anywayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// anywayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	anywayCmd.Flags().StringVarP(&anywayPort, "port", "p", "8080", "本地启动端口")
}
