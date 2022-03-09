package suggests

import "github.com/c-bata/go-prompt"

var TtyOptions = []prompt.Suggest{
	{Text: "-a", Description: "是否开启审计"},
	{Text: "-c", Description: "*.crt文件路径 (default ./server.crt)"},
	{Text: "-k", Description: "*.key文件路径 (default ./server.key)"},
	{Text: "-h", Description: "help for tty"},
	{Text: "-H", Description: "http bind host (default 0.0.0.0)"},
	{Text: "-m", Description: "最大连接数"},
	{Text: "-p", Description: "BasicAuth 密码"},
	{Text: "-P", Description: "http bind port (default 8080)"},
	{Text: "-f", Description: "是否开启pprof性能分析"},
	{Text: "-r", Description: "是否自动重连"},
	{Text: "-t", Description: "是否开启https"},
	{Text: "-u", Description: "BasicAuth 用户名"},
	{Text: "-w", Description: "是否开启写入模式"},
	{Text: "-x", Description: "是否开启xsrf,默认开启"},
	{Text: "-d", Description: "是否打印debug日志"},
	{Text: "-l", Description: "日志是否文件输出"},
}
