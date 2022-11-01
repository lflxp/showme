/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	k8score "github.com/lflxp/lflxp-k8s/core"
	"github.com/lflxp/lflxp-k8s/core/middlewares/jwt/framework"
	"github.com/spf13/cobra"
)

var (
	isHttps bool
)

var kubernetesCmd = &cobra.Command{
	Use:   "k8s",
	Short: "k8s dashboard",
	Long:  `类似rancher、kubesphere paas平台快速查看k8s集群资源工具，后期会添加应用商店和kubevela支持`,
	Run: func(cmd *cobra.Command, args []string) {
		k8score.Run(isHttps)
	},
}

func init() {
	rootCmd.AddCommand(kubernetesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kubernetesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kubernetesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	kubernetesCmd.Flags().BoolVarP(&isHttps, "https", "s", false, "是否开启https")
	kubernetesCmd.Flags().BoolVarP(&framework.IsRancherLogin, "rancher", "r", false, "是否切换为Rancher登录")
}
