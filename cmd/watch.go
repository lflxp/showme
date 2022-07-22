/*
Copyright © 2021 lflxp <382023823@qq.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cosmtrek/air/runner"
	"github.com/spf13/cobra"
)

var (
	cfgPath     string
	debugMode   bool
	showVersion bool
	// cmdArgs     map[string]runner.TomlInfo
	// runArgs     []string
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "go web热加载工具",
	Long:  `Air is yet another live-reloading command line utility for Go applications in development. Just air in your project root directory, leave it alone, and focus on your code.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
  __    _   ___  
 / /\  | | | |_) 
/_/--\ |_| |_| \_ %s, built with Go %s
`, "v1.40.4", "1.19")

		if showVersion {
			return
		}

		if debugMode {
			fmt.Println("[debug] mode")
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		var err error
		// cfg, err := runner.InitConfig(cfgPath)
		// if err != nil {
		// 	log.Fatal(err)
		// 	return
		// }

		// cfg.WithArgs(cmdArgs)
		// r, err := runner.NewEngineWithConfig(cfg, debugMode)
		r, err := runner.NewEngine(cfgPath, debugMode)
		if err != nil {
			log.Fatal(err)
			return
		}
		go func() {
			<-sigs
			r.Stop()
		}()

		defer func() {
			if e := recover(); e != nil {
				log.Fatalf("PANIC: %+v", e)
			}
		}()

		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	watchCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "debug mode")
	watchCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show version")
	watchCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "config path")
}
