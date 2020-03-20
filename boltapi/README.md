bolt模块功能主要是快速提供一个可供CRUD的api接口，保证性能的同时可以抓取监控数据进行数据聚合和记录

![](https://github.com/lflxp/showme/blob/master/img/b1.png)

* web界面
* api接口
* 书序数据库
* 最好能本地监控+grafana展示
* vue+element-ui

# 使用

> showme api -s 

# Prefix界面

主要是提供`基础操作`功能，包括：

* Bucket查询、添加、删除
* Bucket的Key和Value展示
* Key的删除和修改功能
* Prefix前缀查询
* 分页功能
* 刷新功能

![](https://github.com/lflxp/showme/blob/master/img/b2.png)

# Range界面

与Prefix界面唯一的区别就是`搜索功能`，包括：

* 提供按照Key的时间Range进行查询
* 时间数据为：[min,max]的数组，eg：["20200313000000", "20200412235959"]
* 时间格式为：`yyyyMMddHHmmss`

# Orm

TODO：基于Value的时序字段检索

# Backup

提供http方式备份数据的功能

# Swagger

基于gin-swagger提供的api接口可视化界面

![](https://github.com/lflxp/showme/blob/master/img/b3.png)

# 造数据

> while true;do curl -X POST "http://127.0.0.1:8080/api/v1/key/add/monitor/$RANDOM/$RANDOM" -H "accept: application/json";done