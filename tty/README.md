showme tty功能模块提供基于Web的Terminial工具，通过websocket实现远程终端的输入输出的全双工通信机制。

## 特点

`功能点`

* 独立二进制文件运行，无其它依赖
* 提供安全审计
    * 访问统计
    * 命令统计
* 提供安全https连接
* 提供xsrf安全辅助
* 提供性能分析
    * pprof
* 提供接口监控
    * prometheus
* 提供重连机制
* 提供限流机制
* 跨平台支撑
* 提供后台Web可视化管理
    * 查看访问记录
    * 查看操作记录
    * 查看监控记录（prometheus metrics）
    * 全屏显示
* 提供BasicAuth模式
* 提供主机:端口绑定模式
* 提供DEBUG日志模式
* 支持多种远程方式
    * 读/写模式
    * bash模式（默认）
    * 自定义模式
        * tmux
        * docker run --rm -it -p 6379:6379 redis
        * top

`技术栈`

* xterm.js
* pty/tty
* 后端
    * golang
    * websocket
    * prometheus
    * gin
    * restful
* 前端
    * html
    * javascript
    * vue
    * element-ui

![主界面](https://github.com/lflxp/showme/blob/master/img/tty.png)

![后台界面](https://github.com/lflxp/showme/blob/master/img/ttyadmin.png)

## 安装

`环境准备`

```
go get github.com/jteeuwen/go-bindata/...
go get github.com/elazarl/go-bindata-assetfs/...
```

`自动安装`

```shell
git clone https://github.com/lflxp/showme
make install
showme -h
```

`注意`：提供功能选配，比如：go build -tags=gopacket 

## Usage

```
showme tty [flags] [command] [args]
eg: showme tty -w -r showme proxy http

Usage:
  showme tty [flags]

Flags:
  -a, --audit             是否开启审计
  -c, --crt string        *.crt文件路径 (default "./server.crt")
  -d, --debug             debug log mode
  -h, --help              help for tty
  -H, --host string       http bind host (default "0.0.0.0")
  -k, --key string        *.key文件路径 (default "./server.key")
  -m, --maxconnect int    最大连接数
  -p, --password string   BasicAuth 密码
  -P, --port string       http bind port (default "8080")
  -f, --prof              是否开启pprof性能分析
  -r, --reconnect         是否自动重连
  -t, --tls               是否开启https
  -u, --username string   BasicAuth 用户名
  -w, --write             是否开启写入模式
  -x, --xsrf              是否开启xsrf,默认开启

Global Flags:
      --config string   config file (default is $HOME/.showme.yaml)
```

### 使用案例：

* 只读top命令
    > showme tty top
* 读写模式
    > showme tty -w
* 审计、可写、安全、复杂命令混合
    > showme tty -w -a -t -c ./server.crt -k ./server.key -m 10 -u admin -p admin -f -d 'docker run --rm -it -p 6379:7379 redis'

```bash
➜  tls git:(master) ✗ showme tty -w -a -t -c ./server.crt -k ./server.key -m 10 -u admin -p admin -f -d 'docker run --rm -it -p 6379:7379 redis'
INFO[0000] 初始化sqlite数据库 /Users/lxp/.showme.db           
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
[GIN-debug] HEAD   /static/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
[GIN-debug] GET    /metrics                  --> github.com/chenjiandongx/ginprom.PromHandler.func1 (6 handlers)
[GIN-debug] GET    /check                    --> github.com/lflxp/showme/tty.ServeGin.func1 (6 handlers)
[GIN-debug] GET    /who                      --> github.com/lflxp/showme/tty.ServeGin.func2 (6 handlers)
[GIN-debug] GET    /ws                       --> github.com/lflxp/showme/tty.ServeGin.func3 (6 handlers)
[GIN-debug] GET    /                         --> github.com/lflxp/showme/tty.ServeGin.func4 (6 handlers)
[GIN-debug] GET    /debug/pprof/             --> github.com/DeanThompson/ginpprof.IndexHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/heap         --> github.com/DeanThompson/ginpprof.HeapHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/goroutine    --> github.com/DeanThompson/ginpprof.GoroutineHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/allocs       --> github.com/DeanThompson/ginpprof.AllocsHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/block        --> github.com/DeanThompson/ginpprof.BlockHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/threadcreate --> github.com/DeanThompson/ginpprof.ThreadCreateHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/cmdline      --> github.com/DeanThompson/ginpprof.CmdlineHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/profile      --> github.com/DeanThompson/ginpprof.ProfileHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/symbol       --> github.com/DeanThompson/ginpprof.SymbolHandler.func1 (6 handlers)
[GIN-debug] POST   /debug/pprof/symbol       --> github.com/DeanThompson/ginpprof.SymbolHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/trace        --> github.com/DeanThompson/ginpprof.TraceHandler.func1 (6 handlers)
[GIN-debug] GET    /debug/pprof/mutex        --> github.com/DeanThompson/ginpprof.MutexHandler.func1 (6 handlers)
INFO[0000] Listening and serving HTTPS on 192.168.0.2:8080  tty.go=261
```

## 环境设置

使用`systemctl`进行service部署showme的时候会报`TERM environment variable not set`，这个需要在service文件里面指定环境变量`TERM=xterm-256color`

```
root@8.8.8.8:/etc/systemd/system# cat showme.service 
[Unit]
Description=showme
After=syslog.target
After=network.target

[Service]
# Modify these two values and uncomment them if you have
# repos with lots of files and get an HTTP error 500 because
# of that
###
#LimitMEMLOCK=infinity
#LimitNOFILE=65535
Type=simple
User=root
Group=root
WorkingDirectory=/tls
ExecStart=/usr/bin/showme tty -P 9999 -w -a -t -f -m 10 -u $user -p $pwd -c /tls/server.crt -k /tls/server.key
# ExecReload=/bin/kill -s HUP $MAINPID
Restart=always
Environment=USER=root HOME=/root TERM=xterm-256color

[Install]
WantedBy=multi-user.target
```