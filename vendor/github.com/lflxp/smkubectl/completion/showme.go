package completion

var showme = Completion{
	Level: "showme",
	Cmd: []string{
		"COMMAND     DESCRIPTION",
		"api         快速本地DB CRUD API",
		"cmd         命令聚合工具",
		"completion  A brief description of your command",
		"dashboard   showme dashboard快速学习平台",
		"help        显示showme帮助文档",
		"k8s         k8s dashboard",
		"martix      黑客帝国字母雨特效",
		"music       本地在线音乐网站",
		"playbook    批量主机任务编排脚本执行器",
		"proxy       代理服务器",
		"scan        A brief description of your command",
		"smart       A brief description of your command",
		"static      本地静态文件服务器",
		"tty         web terminial",
		"watch       go web热加载工具",
	},
	Daughter: map[string]Completion{
		"tty": Completion{
			Level: "tty",
			Cmd: []string{
				"ARGS                    DESCRIPTION",
				"-a  --audit             是否开启审计",
				"-c  --crt string        *.crt文件路径 (default ./server.crt)",
				"-d  --debug             debug log mode",
				"-h  --help              help for tty",
				"-H  --host string       http bind host (default \"0.0.0.0\")",
				"-k  --key string        *.key文件路径 (default \"./server.key\")",
				"-m  --maxconnect int    最大连接数",
				"-p  --password string   BasicAuth 密码",
				"-P  --port string       http bind port (default 8080)",
				"-f  --prof              是否开启pprof性能分析",
				"-r  --reconnect         是否自动重连",
				"-t  --tls               是否开启https",
				"-u  --username string   BasicAuth 用户名",
				"-w  --write             是否开启写入模式",
				"-x  --xsrf              是否开启xsrf,默认开启",
			},
		},
		"watch": Completion{
			Level: "watch",
			Cmd: []string{
				"ARGS                    DESCRIPTION",
				"-c  --config string     config path",
				"-d  --debug             debug mode",
				"-h  --help              help for watch",
				"-v  --version           show version",
			},
		},
		"static": Completion{
			Level: "static",
			Cmd: []string{
				"ARGS                      DESCRIPTION",
				"-C  --clean string        删除文件或文件夹,如：/tmp",
				"-T  --contains string     删除文件名包含的内容，如：.temp",
				"-d  --debug               是否开启断点续传debug日志",
				"-D  --dest string         复制文件目标文件或文件夹",
				"-h  --help                help for static",
				"-c  --pagesize int        每页显示视频数 (default 20)",
				"-f  --path string         加载目录 (default ./)",
				"-p  --port string         服务端口 (default 9090)",
				"-r  --raw                 是否切换为无html页面状态",
				"-S  --src string          复制文件原文件或文件夹",
				"-P  --staticPort string   文件服务端口 (default 9091)",
				"-t  --types string        过滤视频类型，多个用逗号隔开 (default .avi,.wma,.rmvb,.rm,.mp4,.mov,.3gp,.mpeg,.mpg,.mpe,.m4v,.mkv,.flv,.vob,.wmv,.asf,.asx)",
				"-v  --video               是否切换为视频模式",
			},
		},
		"completion": Completion{
			Level: "completion",
			Cmd: []string{
				"COMMAND",
				"bash",
				"zsh",
			},
		},
		"scan": Completion{
			Level: "scan",
			Cmd: []string{
				"ARGS                DESCRIPTION",
				"-d  --debug         Debug Log Info",
				"-h  --help          help for scan",
				"-i  --ip string     ip范围,例如:127.0.0.1 (default 127.0.0.1)",
				"-p  --port string   端口范围，例如: 22,80,8080-9000 (default 1-65535)",
			},
		},
		"cmd": Completion{
			Level: "cmd",
			Cmd: []string{
				"COMMAND             DESCRIPTION",
				"dashboard           computer configeration",
				"gocui               https://github.com/jroimartin/gocui ",
				"monitor             monitoring Linux/Unix or MacOs status runtime  ",
				"scan                IP端扫工            ",
				"mysql               monitor mysql info     ",
				"help                List All Menu    ",
				"tty                 web terminal     ",
				"sw                  golang持热新具",
			},
		},
		"proxy": Completion{
			Level: "proxy",
			Cmd: []string{
				"COMMAND             DESCRIPTION",
				"http                http反向代理",
				"socket5             socket5 http代理服务器",
				"ss                  shadowsocks",
				"tcp                 tcp 四层反向代理",
			},
		},
	},
}