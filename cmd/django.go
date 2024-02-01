/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lflxp/djangolang"
	"github.com/lflxp/djangolang/middlewares"
	"github.com/spf13/cobra"
)

type Demotest2 struct {
	Id         int64     `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Country    string    `json:"country" xorm:"varchar(255) not null" search:"true"`
	Zoom       string    `json:"zoom" xorm:"varchar(255) not null"`
	Company    string    `json:"company" xorm:"varchar(255) not null"`
	Items      string    `json:"items" xorm:"varchar(255) not null"`
	Production string    `json:"production" xorm:"varchar(255) not null"`
	Count      string    `json:"count" xorm:"varchar(255) not null"`
	Serial     string    `json:"serial" xorm:"varchar(255) not null" search:"true"`
	Extend     string    `json:"extend" xorm:"varchar(255) not null"`
	Files      string    `xorm:"file" name:"file" verbose_name:"上传文件" colType:"file"`
	File2      string    `xorm:"file2" name:"file2" verbose_name:"上传文件2" colType:"file"`
	Type       string    `xorm:"type" name:"type" verbose_name:"类型" search:"false" colType:"textarea"`
	Detail     string    `xorm:"detail" name:"detail" verbose_name:"VPN信息" list:"false" search:"false" o2m:"vpn|id,vpn" colType:"o2m"`
	Times      time.Time `xorm:"times" name:"times" verbose_name:"时间" colType:"time" list:"true" search:"true"`
}

// djangoCmd represents the django command
var djangoCmd = &cobra.Command{
	Use:   "django",
	Short: "golang版本django的web框架",
	Long:  `基于golang的web框架，类似于python的django框架，但是更加简单易用，更加高效`,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				// 打印异常，关闭资源，退出此函数
				fmt.Println(err)
			}
		}()
		djangolang.RegisterAdmin(new(Demotest2))

		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/admin/index")
		})

		djangolang.RegisterControllerAdmin(r, false)
		middlewares.RegisterPrometheusMiddlewareBasic(r, false)
		r.Run() // listen and serve on
	},
}

func init() {
	rootCmd.AddCommand(djangoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// djangoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// djangoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
