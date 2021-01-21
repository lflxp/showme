/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/devopsxp/xp/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// playbookCmd represents the playbook command
var playbookCmd = &cobra.Command{
	Use:   "playbook",
	Short: "批量主机任务编排脚本执行器",
	Long:  `测试ansile-playbook功能和pipeline流程管控`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("playbook called")

		// 根据yaml解析shell等模块，进行动态匹配，进行顺序执行
		config := pipeline.DefaultPipeConfig("shell").
			WithInputName("localyaml").
			WithFilterName("shell").
			WithOutputName("console")

		p := pipeline.Of(*config)
		p.Init()
		p.Start()
		p.Exec()
		p.Stop()
	},
}

func init() {
	rootCmd.AddCommand(playbookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playbookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playbookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
