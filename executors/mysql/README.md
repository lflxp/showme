mysql Orzdba go版本性能分析工具

![](https://github.com/lflxp/showme/blob/master/img/mysql.png)
![](https://github.com/lflxp/showme/blob/master/img/mysql2.png)

## 特色

* 本地/远程连接
* Databases
* Variables
    * binlog_format
    * open_files_limit
    * ...
* dashboard
* test
    * GetHostAndIps 
    * GetShowDatabases   
    * GetShowGlobalVariables  
    * GetShowVariables           
    * GetShowStatus        
    * GetShowGlobalStatus     
* processlist
* status
    * innodb信息
    * CRUD信息
    * Net IO
    * Slave
    * Threads
    * Semi

## 使用

mysql是showme终端GUI显示，提供动态参数提示功能。

```bash
➜  tls git:(master) ✗ showme
>>> mysql status -H 10.135.26.4 -p system -innodb -com -lazy
```

## 参数

* -H Mysql连接主机，默认127.0.0.1 (default "127.0.0.1")
* -P Mysql连接端口,默认3306 (default "3306")
* -L Print to Logfile. (default "none")
* -u mysql用户名
* -p mysql密码
* -db Mysql 指定databases,默认：mysql
* -S mysql socket连接文件地址 (default "/tmp/mysql.sock")
* -mysql Print MySQLInfo (include -t,-com,-hit,-T,-B).
* -v 详细信息 【informational】
* -lazy `Info  (include -t,-l,-c,-s,-com,-hit)`
    * -------QPS----------TPS------- ----KeyBuffer------Index----Qcache---Innodb---(%)
* -com `MySQL Status(Com_select,Com_insert,Com_update,Com_delete)`
    * -------QPS----------TPS-------
* -hit `Innodb Hit%.`
    * ----KeyBuffer------Index----Qcache---Innodb---(%)
* -innodb `InnodbInfo(include -t,-innodb_pages,-innodb_data,-innodb_log,..`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query---
* -innodb_rows `Innodb Rows Status(Innodb_rows_inserted/updated/deleted/read)`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query--- ---innodb rows status---
* -innodb_pages `Innodb Buffer Pool Pages Status(Innodb_buffer_pool_pages_data`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query--- ---innodb bp pages status--
* -innodb_data `Innodb Data Status(Innodb_data_reads/writes/read/written)`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query--- -----innodb data status-----
* -innodb_log `Innodb Log  Status(Innodb_os_log_fsyncs/written)`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query--- --innodb log--
* -innodb_status `Innodb Status from Command: 'Show Engine Innodb Status'`
    * ---innodb bp pages status-- -----innodb data status----- --innodb log--   his --log(byte)--  read ---query---   his --log(byte)--  read ---query---
* -T `Threads Status(Threads_running,Threads_connected,Threads_crea..`
    * ----------threads---------
* -B `Bytes received from/send to MySQL(Bytes_received,Bytes_sent)`
    * -----bytes----
* -semi `半同步监控`
    * ---avg_wait--tx_times--semi
* -slave `SLAVE INFO`
    * ---------------SlaveStatus-------------
