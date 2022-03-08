/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/lflxp/showme/pkg/gomatrix"
	"github.com/spf13/cobra"
)

var (
	ascii   bool
	logging bool
	profile string
	fps     int
)

// martixCmd represents the martix command
var martixCmd = &cobra.Command{
	Use:   "martix",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gomatrix.Run(ascii, logging, profile, fps)
	},
}

func init() {
	rootCmd.AddCommand(martixCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// martixCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// martixCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	martixCmd.Flags().BoolVarP(&ascii, "ascii", "a", false, "Use ascii/alphanumeric characters instead of japanese kana's.")
	martixCmd.Flags().BoolVarP(&logging, "logging", "L", false, "Enable logging debug messages to ~/.gomatrix-log.")
	martixCmd.Flags().StringVarP(&profile, "profile", "p", "", "Write profile to given file path")
	martixCmd.Flags().IntVarP(&fps, "fps", "f", 25, "required FPS, must be somewhere between 1 and 60 default: 25")
}
