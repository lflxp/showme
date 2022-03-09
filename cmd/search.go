/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/lflxp/showme/pkg/search"
	"github.com/spf13/cobra"
)

var uninclude []string
var show bool

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "全局模糊搜索工具",
	Long: `模仿fzf进行全局资源查找，包括：
	1. 文件和文件夹
	2. todo: 历史命令
	3. todo: 远程gitee备份
	4. todo: 环境变量
	5. todo: 文件预览？
	6. todo: showme内置功能预览和快速导航（scan、monitor等）
	7. todo: linux内核参数
	8. todo: zsh、vim、go、rust配置等
	9. todo: alias快捷键管理
	10. todo: 命令管理
	11. todo: 排序，自动识别时间字段、数字字段 sort.SliceStable`,
	Run: func(cmd *cobra.Command, args []string) {
		search.Run(uninclude, show)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	searchCmd.Flags().StringSliceVarP(&uninclude, "uninclude", "e", []string{".git"}, "不扫描目录,默认: ['.git']")
	searchCmd.Flags().BoolVarP(&show, "show", "s", false, "Show loading ui")
}
