package suggests

import "github.com/c-bata/go-prompt"

var MonitorOptions = []prompt.Suggest{
	{Text: "-B", Description: "Print Bytes received from/send to MySQL(Bytes_received,Bytes_sent)."},
	{Text: "-C", Description: "INT 运行时间 默认无限"},
	{Text: "-H", Description: "STRING Mysql连接主机，默认127.0.0.1 (default \"127.0.0.1\")"},
	{Text: "-L", Description: "STRING Print to Logfile. (default \"none\")"},
	{Text: "-P", Description: "STRING Mysql连接端口,默认3306 (default \"3306\")"},
	{Text: "-S", Description: "STRING mysql socket连接文件地址 (default \"/tmp/mysql.sock\")"},
	{Text: "-T", Description: "Print Threads Status(Threads_running,Threads_connected,Threads_created,Threads_cached)."},
	{Text: "-c", Description: "打印Cpu info"},
	{Text: "-com", Description: "Print MySQL Status(Com_select,Com_insert,Com_update,Com_delete)."},
	{Text: "-d", Description: "STRING 打印Disk info (default \"none\")"},
	{Text: "-hit", Description: "Print Innodb Hit%."},
	{Text: "-i", Description: "STRING 时间间隔 默认1秒 (default \"1\")"},
	{Text: "-l", Description: "打印Load info"},
	{Text: "-lazy", Description: "Print Info  (include -t,-l,-c,-s,-com,-hit)."},
	{Text: "-n", Description: "STRING 打印net info (default \"none\")"},
	{Text: "-nocolor", Description: "不显示颜"},
	{Text: "-s", Description: "打印swap info"},
	{Text: "-nocolor", Description: "不显示颜"},
	{Text: "-t", Description: "打印当前时间"},
}
