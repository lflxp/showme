static功能主要是快速启动一个http服务进行文件的传输，包括文件上传和下载，摆脱无工具可用的尴尬境地。

目前static新增了视频模式，通过过滤常用的视频文件格式在前端通过video标签进行直接播放，本地离线视频服务。

![](https://github.com/lflxp/showme/blob/master/img/httpstatic.png)

## 功能点

* Web界面操作
* 静态文件下载
* 批量文件上传
* 视频文件播放
* prometheus监控

## 使用

> showme static -v -p 9091 -c 20 -f /tmp

## 参数

```
通过本地http服务进行简单的文件传输和文件展示

Usage:
  showme static [flags]

Flags:
  -h, --help           help for static
  -c, --pagesize int   每页显示视频数 (default 20)
  -f, --path string    加载目录 (default "./")
  -p, --port string    服务端口 (default "9090")
  -t, --types string   过滤视频类型，多个用逗号隔开 (default ".avi,.wma,.rmvb,.rm,.mp4,.mov,.3gp,.mpeg,.mpg,.mpe,.m4v,.mkv,.flv,.vob,.wmv,.asf,.asx")
  -v, --video          是否切换为视频模式

Global Flags:
      --config string   config file (default is $HOME/.showme.yaml)
```

## 优化

* web页面进行全功能操作
* web页面进行文件下载
* web页面进行上传（无curl命令操作，方便快捷）
* web页面查看监控指标，可对接prometheus server监控
* web视频文件直接加载观看功能