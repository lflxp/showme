# 介绍

linux fast command for ops & dev

# Usage

tags:

- gopacket

```
go get github.com/lflxp/showme
go run main.go
go build -tags=gopacket main.go
```

# Operation

1. 命令提示

- monitor args...
- dashboard args...
- gocui active
- gopacket args...
- httpstatic

![command.png](https://github.com/lflxp/showme/blob/master/command.png)

2. monitor 监控展示

![monitor.png](https://github.com/lflxp/showme/blob/master/monitor.png)

- {Text: "-L", Description: "Print to Logfile. (default \"none\")"}
- {Text: "-c", Description: "打印Cpu 信息负载信息"}
- {Text: "-d", Description: "打印Disk info (default \"none\")"}
- {Text: "-i", Description: "STRING 时间间隔 默认1秒 (default \"1\")"}
- {Text: "-l", Description: "打印Load 信息"}
- {Text: "-lazy", Description: "Print Info  (include -t,-l,-c,-s,-n)."}
- {Text: "-n", Description: "打印net网络流量"}
- {Text: "-N", Description: "打印net网络详细流量"}
- {Text: "-s", Description: "打印swap 信息"}
- {Text: "-t", Description: "打印当前时间"}

3. scan ip and port

![scan.png](https://github.com/lflxp/showme/blob/master/scan.png)

- Tab: Next View
- Enter: Select IP/Commit Input
- F5: Input New Scan IP or Port range
- ↑ ↓: Move View
- ^c: Exit
- F1: Help
- Space: search result with ip view and port view

4. static http server

> httpstatic -port 9999

- port: static http port, default: 9090

5. gopacket 网络流量监控

command: 
- gopacket interface eth0
- gopacket screen eth0

> TIPS: 条件编译，需要该功能需要在build的指定tags
> eg: go build -tags=gopacket main.go

# feature

1. 非网络连接
2. 非vue类界面操作
3. 速度
4. 本地性能、监控、操作
5. 快速定位
6. command visio
7. 网络监听（原始报文解析和展示）

# functions

1. go-prompt
2. gocui/tcell/tview/ncurses/goncurses
3. 提示选项分为两种： 一、命令参数 dashboard status 二、配置参数 dashboard --status
4. github.com/jroimartin/gocui
5. github.com/gdamore/tcell

# 新增操作

0. suggests/commands.go 添加首字符命令添加提示
1. completers/options.go 添加含【-】的参数
2. completers/common.go -> FirstCommandFunc 添加命令提示 添加基于首字符的二级字符命令提示
3. executors 添加目录实现命令gocui展示
4. executors/executors.go 添加command对应的执行命令

# 多线程改造

https://blog.csdn.net/lengyuezuixue/article/details/79664409

# todo

- mysql 解析
- 微服务管理和功能测试