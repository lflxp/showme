package suggests

import "github.com/c-bata/go-prompt"

var MonitorOptions = []prompt.Suggest{
	{Text: "-C", Description: "运行时间 默认无限"},
	{Text: "-L", Description: "Print to Logfile. (default \"none\")"},
	{Text: "-c", Description: "打印Cpu 信息负载信息"},
	{Text: "-d", Description: "打印Disk info (default \"none\")"},
	{Text: "-i", Description: "STRING 时间间隔 默认1秒 (default \"1\")"},
	{Text: "-l", Description: "打印Load 信息"},
	{Text: "-lazy", Description: "Print Info  (include -t,-l,-c,-s,-n)."},
	{Text: "-n", Description: "打印net网络流量"},
	{Text: "-N", Description: "打印net网络详细流量"},
	{Text: "-s", Description: "打印swap 信息"},
	{Text: "-nocolor", Description: "不显示颜"},
	{Text: "-t", Description: "打印当前时间"},
}
