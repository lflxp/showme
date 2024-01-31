//go:build gopacket
// +build gopacket

/*
Copyright © 2020 Lixueping <382023823@qq.com>

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
	"log/slog"
	"net"
	"time"

	"github.com/spf13/cobra"

	"github.com/lflxp/lflxp-sflowtool/pkg"
)

var (
	Con pkg.Collected = pkg.Collected{
		SnapShotLen: 65535,
		Promiscuous: true,
		Timeout:     30 * time.Second,
	}
	sflowItem        string
	sflowProtocol    string
	sflowPort        string
	sflowEth         string
	sflowUDP         bool
	sflowUDPort      string
	sflowCounterport string
	sflowEsurl       string
	sflowIses        bool
	sflowDebug       bool
	sflowIndex       string
)

func SflowCounter(protocol, port string) {
	Con.ListenSflowCounter(protocol, port)
}

func SflowSample(protocol, port string) {
	Con.ListenSFlowSample(protocol, port)
}

// include SFlowSample and SflowCounter
func SflowAll(protocol, port string) {
	Con.ListenSflowAll(protocol, port)
}

func NetflowV5(protocol, port string) {
	Con.ListenNetFlowV5(protocol, port)
}

// sflowtoolCmd represents the sflowtool command
var sflowtoolCmd = &cobra.Command{
	Use:   "sflowtool",
	Short: "实时流量分析",
	Long: `elasticsearch + grafana + sflow
该工具是接受来自交换机的sflow流量进行解析，然后上传到es集群进行数据汇总，最后通过grafana进行数据查询聚合展示。
运行：sudo ./sflowtool -e enp1s0 -p 9999 -t all -i
注意：需要安装依赖pcap.h, yum install libpcap-devel`,
	Run: func(cmd *cobra.Command, args []string) {
		wait := make(chan int)

		Con.DeviceName = sflowEth
		Con.Host = sflowUDPort
		Con.Udpbool = sflowUDP
		Con.CounterHost = sflowCounterport
		Con.EsPath = sflowEsurl
		Con.IsEs = sflowIses
		Con.Index = sflowIndex

		// 设置日志级别
		if debug {
			lvl.Set(slog.LevelDebug)
		} else {
			lvl.Set(slog.LevelInfo)
		}

		// 初始化es index
		if Con.IsEs {
			slog.Info("开启es通道")
			pkg.InitEs(Con.EsPath, Con.Index)
		}

		// 是否开启udp数据转发
		if udp {
			Conn, err := net.Dial("udp", sflowUDPort)
			defer Conn.Close()
			if err != nil {
				panic(err)
			}
		}

		// 启动命令
		if sflowItem == "all" {
			SflowAll(sflowProtocol, sflowPort)
		} else if sflowItem == "counter" {
			SflowCounter(sflowProtocol, sflowPort)
		} else if sflowItem == "sample" {
			SflowSample(sflowProtocol, sflowPort)
		} else if sflowItem == "netflow" {
			NetflowV5(sflowProtocol, sflowPort)
		}
		<-wait
	},
}

func init() {
	rootCmd.AddCommand(sflowtoolCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sflowtoolCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sflowtoolCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sflowtoolCmd.Flags().StringVarP(&sflowItem, "item", "t", "all", "类型:all(sflowSample|Counter),counter(SflowCounter),sample(SflowSample),netflow")
	sflowtoolCmd.Flags().StringVarP(&sflowProtocol, "protocol", "l", "udp", "协议")
	sflowtoolCmd.Flags().StringVarP(&sflowPort, "port", "p", "6343", "监听端口")
	sflowtoolCmd.Flags().StringVarP(&sflowEth, "eth", "e", "en0", "接受sflow流量的物理网络设备名，eg: en0")
	sflowtoolCmd.Flags().BoolVarP(&sflowUDP, "udp", "u", false, "是否开启udp数据传输,默认不开启")
	sflowtoolCmd.Flags().StringVarP(&sflowUDPort, "udport", "y", "127.0.0.1:6666", "udp SFlowSample And Netflow 传输主机:端口")
	sflowtoolCmd.Flags().StringVarP(&sflowCounterport, "counterport", "c", "127.0.0.1:7777", "udp CounterSample 传输主机:端口")
	sflowtoolCmd.Flags().StringVarP(&sflowEsurl, "esurl", "s", "http://127.0.0.1:9200", "elasticsearch 5.6 接口地址")
	sflowtoolCmd.Flags().BoolVarP(&sflowIses, "ises", "i", false, "是否开启output到elasticsearch")
	sflowtoolCmd.Flags().BoolVarP(&sflowDebug, "debug", "d", false, "是否开启debug model")
	sflowtoolCmd.Flags().StringVarP(&sflowIndex, "index", "I", "sflow", "es index name, example: sflow-2019-09-06")
}
