/*
Copyright © 2024 lflxp <382023823@qq.com>
*/
package cmd

import (
	"log/slog"
	"os"

	fzf "github.com/junegunn/fzf/src"
	"github.com/junegunn/fzf/src/protector"
	"github.com/spf13/cobra"
)

var version string = "0.46"
var revision string = "devel"
var lvl slog.LevelVar

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "smkubectl",
	Short: "A brief descriptsion of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// parseCmd(args[1:])
	// 	fmt.Println(strings.Join(args, " "))
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 日志配置
	lvl.Set(slog.LevelError)
	opts := slog.HandlerOptions{
		AddSource: true,
		Level:     &lvl,
	}

	// slog.SetDefault(slog.New((slog.NewJSONHandler(os.Stdout, &opts))))
	slog.SetDefault(slog.New((slog.NewTextHandler(os.Stdout, &opts))))

	// 保证没有参数或者参数只有一个且为completion的时候执行cobra
	// 其余都走parseCmd
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "completion", "smart":
			slog.Debug("进入cobra命令模式", slog.Any("args", os.Args[1]))
			err := rootCmd.Execute()
			if err != nil {
				os.Exit(1)
			}
		default:
			slog.Debug("进入completion模式")
			protector.Protect()
			fzf.Run(fzf.ParseOptions(), version, revision)
		}
	} else {
		slog.Debug("进入fzf模式")
		protector.Protect()
		fzf.Run(fzf.ParseOptions(), version, revision)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.demo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
