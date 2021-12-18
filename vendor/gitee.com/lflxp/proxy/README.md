# 介绍

学习、使用、封装golang的反向代理，实现HTTP和TCP的反向代理功能

# 功能

## 自定义反向代理

`func NewHttpProxyByGinCustom(target string, filter map[string]string) func(c *gin.Context)`生成gin框架适应的7层反向代理。

参数说明:

1. target 需要反向代理的服务器，可以是域名也可以是ip+端口，需要写全路径：https://www.baidu.com 或者 http://127.0.0.1:9906 ,`其中http:// 或 https:// Schema不能少`
2. filter 是需要鉴权的Header，可以填nil不进行校验，校验必须满足才能通过

## TCP反向代理

`func NewTCPProxy(from, to string) error` 原地启动tcp服务

参数说明：

1. from 本地启动服务，例如: 127.0.0.1:6000 或者 0.0.0.0:6000 `不需要http://`
2. to 需要反向代理的后端服务器，例如： 127.0.0.1:3306

# 运行

> make

```
root in proxy on  master [✘!?] via 🐹 v1.16.6 took 37s 
❯ make 
go run examples/main.go
2021-10-28T23:05:31+08:00 | INFO  | ***tcp-proxy started.**** backend:127.0.0.1:9998 bind:0.0.0.0:9999
2021-10-28T23:05:36+08:00 | INFO  | ***client connected.**** conn:172.21.64.1:58387
2021-10-28T23:05:36+08:00 | INFO  | ***client connected.**** conn:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | INFO  | ***backend connected.**** backend:127.0.0.1:34048 conn:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | INFO  | ***backend connected.**** backend:127.0.0.1:34046 conn:172.21.64.1:58387
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:172.21.64.1:58388 recv:549 send:549 to:127.0.0.1:9998
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:3090 send:3090 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:172.21.64.1:58388 recv:462 send:462 to:127.0.0.1:9998
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
2021-10-28T23:05:36+08:00 | DEBUG | ***Proxyying**** from:127.0.0.1:9998 recv:4096 send:4096 to:172.21.64.1:58388
```

# 资料

https://blog.csdn.net/mengxinghuiku/article/details/65448600