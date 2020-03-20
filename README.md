# Introduction 

 Rapid diagnostic system status tool (performance monitoring, network scanning, mysql performance monitoring, kubectl status)

> 快速诊断系统状态工具（性能监控、网络扫描、mysql性能监控）

# Install

```bash
go get -u github.com/lflxp/showme
make install
showme -h
```

# Module

* [【DONE】tty 【WEB TERMINIAL】](./tty/README.md)

![](https://github.com/lflxp/showme/blob/master/img/tty.png)

* [【DONE】gopacket 网络流量分析](./executors/gopacket/README.md)

![](https://github.com/lflxp/showme/blob/master/img/gopacket.png)

* [【DONE】monitor 监控展示](./executors/monitor/README.md)

![monitor.png](https://github.com/lflxp/showme/blob/master/img/monitor.png)

* [【DONE】mysql 数据库监控](./executors/mysql/README.md)

![](https://github.com/lflxp/showme/blob/master/img/mysql2.png)

* [【DONE】scan 扫描工具](./executors/scan/README.md)

![](https://github.com/lflxp/showme/blob/master/img/scan.png)

* [【DONE】static 文件传输](./httpstatic/README.md)

![](https://github.com/lflxp/showme/blob/master/img/httpstatic.png)

* [【TO BE FIX】k8s 管理工具](./executors/kubectl/README.md)

![s1.png](https://github.com/lflxp/showme/blob/master/img/s1.png)

* [bolt 快速RESTFUL API](./boltapi/README.md)

![](https://github.com/lflxp/showme/blob/master/img/b1.png)

* [【TODO】proxy 代理工具](#PROXY)
  * http正向代理
  * http 反向代理
  * mysql tcp代理（负载均衡、读写分离、分布式调度）
  * socket5 代理
  * ss fq代理
  * ss server

* [【TODO】playbook 任务编排工具](#PLAYBOOK)

* [【TODO】nmap 高级扫描工具](#NMAP)

## NMAP

基于优秀的nmap工具进行封装，采用`gin`+`api`+`restful`+`remote`的方式进行远程调用。

## PLAYBOOK

基于Ansible-playbook开发的Go原型工具，`功能特点`有：

> showme playbook

* RPC远程操作
* Yaml Template
* Go Template
* Plugin Register
* Mini CMDB Required

## PROXY

基于GOLANG的各种代理工具，处于测试阶段。

`Usage`

```bash
➜  showme git:(master) ✗ showme proxy -h  
* http正向代理
* http 反向代理
* mysql tcp代理（负载均衡、读写分离、分布式调度）
* socket5 代理
* ss fq代理

Usage:
  showme proxy [command]

Available Commands:
  http        http正向代理
  httpreverse http反向代理
  mysql       mysql proxy
  socket5     socket5 http代理服务器
  ss          shadowsocks
```  

# Technology Stack

1. go-prompt
2. gocui/tcell/tview/ncurses/goncurses
3. 提示选项分为两种： 一、命令参数 dashboard status 二、配置参数 dashboard --status
4. github.com/jroimartin/gocui
5. github.com/gdamore/tcell
6. vue 【Element-UI】
7. websocket
8. gin web api
9. yaml

# 新增操作

0. suggests/commands.go 添加首字符命令添加提示
1. completers/options.go 添加含【-】的参数
2. completers/common.go -> FirstCommandFunc 添加命令提示 添加基于首字符的二级字符命令提示
3. executors 添加目录实现命令gocui展示
4. executors/executors.go 添加command对应的执行命令

# 多线程改造

https://blog.csdn.net/lengyuezuixue/article/details/79664409

# TODO

- 结合GuiLite进行美化
- tty 添加install自动部署systemctl服务的功能
- 对接Mini CMDB

# k8s resource list

```
alertmanagers.monitoring.coreos.com                           endpoints                                                     nodes.metrics.k8s.io                                          replicasets.extensions
apiservices.apiregistration.k8s.io                            etcdclusters.etcd.database.coreos.com                         persistentvolumeclaims                                        replicationcontrollers
certificatesigningrequests.certificates.k8s.io                events                                                        persistentvolumes                                             resourcequotas
clusterrolebindings.rbac.authorization.k8s.io                 events.events.k8s.io                                          poddisruptionbudgets.policy                                   rolebindings.rbac.authorization.k8s.io
clusterroles.rbac.authorization.k8s.io                        horizontalpodautoscalers.autoscaling                          pods                                                          roles.rbac.authorization.k8s.io
componentstatuses                                             ingresses.extensions                                          pods.metrics.k8s.io                                           secrets
configmaps                                                    jobs.batch                                                    podsecuritypolicies.extensions                                serviceaccounts
controllerrevisions.apps                                      leases.coordination.k8s.io                                    podsecuritypolicies.policy                                    servicemonitors.monitoring.coreos.com
cronjobs.batch                                                limitranges                                                   podtemplates                                                  services
customresourcedefinitions.apiextensions.k8s.io                mutatingwebhookconfigurations.admissionregistration.k8s.io    priorityclasses.scheduling.k8s.io                             statefulsets.apps
daemonsets.apps                                               namespaces                                                    prometheuses.monitoring.coreos.com                            storageclasses.storage.k8s.io
daemonsets.extensions                                         networkpolicies.extensions                                    prometheusrules.monitoring.coreos.com                         validatingwebhookconfigurations.admissionregistration.k8s.io
deployments.apps                                              networkpolicies.networking.k8s.io                             redisfailovers.storage.spotahome.com                          volumeattachments.storage.k8s.io
deployments.extensions                                        nodes                                                         replicasets.apps                                              
```
