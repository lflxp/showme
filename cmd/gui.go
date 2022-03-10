/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/lflxp/showme/pkg/vcls"
	"github.com/lflxp/showme/pkg/vcls/form"
	"github.com/spf13/cobra"
)

var sample bool

// guiCmd represents the gui command
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "GUI图形化界面",
	Long: `基于GoVCL进行开发，跨平台部署
	参考：https://gitee.com/ying32/govcl/wikis/pages?sort_id=1540640&doc_id=102420
	依赖包: https://gitee.com/lflxp/govcl/attach_files/992338/download/liblcl-2.2.0.zip`,
	Run: func(cmd *cobra.Command, args []string) {
		if sample {
			vcls.Run()
		} else {
			form.Run()
		}

	},
}

func init() {
	rootCmd.AddCommand(guiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// guiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// guiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	guiCmd.Flags().BoolVarP(&sample, "sample", "s", false, "sample or full function")
}
