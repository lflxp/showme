scan是扫描模块，提供：

* ip扫描
    * 192.168.0.1
* ip段扫描
    * 192.168.1-10.55
    * 192.168.1.1-255
* 端口扫描
    * 1-65535
* 扫描结果展示

![](https://github.com/lflxp/showme/blob/master/img/scan1.png)
![](https://github.com/lflxp/showme/blob/master/img/scan.png)

## 使用

scan是showme终端GUI显示，提供动态参数提示功能。

> showme

```bash
=> host@127.0.0.1 # showme
>>> scan 192.168.50.1-255
            192.168.50.1-255  192网段  
            192.168.40.1-255  192网段  
            192.168.1.1-255   192网段  
```

## 参数

`IP段参数说明`

* 支持单IP
* 支持【-】指定范围IP，包括A、B、C、D段

`GUI操作说明`

- Tab: Next View
- Enter: Select IP/Commit Input
- F5: Input New Scan IP or Port range
- ↑ ↓: Move View
- ^c: Exit
- F1: Help
- Space: search result with ip view and port view