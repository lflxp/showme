/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	music "github.com/lflxp/lflxp-music/core"
	"github.com/spf13/cobra"
)

var (
	isHttps2 bool
)

var musicCmd = &cobra.Command{
	Use:   "music",
	Short: "本地在线音乐网站",
	Long:  `> mmPlayer 是由茂茂开源的一款在线音乐播放器，具有音乐搜索、播放、歌词显示、播放历史、查看歌曲评论、网易云用户歌单播放同步等功能`,
	Run: func(cmd *cobra.Command, args []string) {
		music.Run(isHttps2)
	},
}

func init() {
	rootCmd.AddCommand(musicCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// musicCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// musicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	musicCmd.Flags().BoolVarP(&isHttps2, "https", "s", false, "是否开启https")
}
