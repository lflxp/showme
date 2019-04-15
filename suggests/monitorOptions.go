package suggests

import "github.com/c-bata/go-prompt"

var MonitorOptions = []prompt.Suggest{
	{Text: "-B", Description: "Print Bytes received from/send to MySQL(Bytes_received,Bytes_sent)."},
	{Text: "-C", Description: "运行时间 默认无限"},
	{Text: "-H", Description: "Mysql连接主机，默认127.0.0.1 (default \"127.0.0.1\")"},
	{Text: "-L", Description: "Print to Logfile. (default \"none\")"},
	{Text: "-P", Description: "Mysql连接端口,默认3306 (default \"3306\")"},
	{Text: "-S", Description: "mysql socket连接文件地址 (default \"/tmp/mysql.sock\")"},
	{Text: "-T", Description: "Print Threads Status(Threads_running,Threads_connected,Threads_created,Threads_cached)."},
	{Text: "-c", Description: "打印Cpu 信息负载信息"},
	{Text: "-com", Description: "Print MySQL Status(Com_select,Com_insert,Com_update,Com_delete)."},
	{Text: "-d", Description: "打印Disk info (default \"none\")"},
	{Text: "-hit", Description: "Print Innodb Hit%."},
	{Text: "-i", Description: "STRING 时间间隔 默认1秒 (default \"1\")"},
	{Text: "-l", Description: "打印Load 信息"},
	{Text: "-lazy", Description: "Print Info  (include -t,-l,-c,-s,-n)."},
	{Text: "-n", Description: "打印net网络流量"},
	{Text: "-N", Description: "打印net网络详细流量"},
	{Text: "-nocolor", Description: "不显示颜"},
	{Text: "-s", Description: "打印swap 信息"},
	{Text: "-nocolor", Description: "不显示颜"},
	{Text: "-t", Description: "打印当前时间"},
}
