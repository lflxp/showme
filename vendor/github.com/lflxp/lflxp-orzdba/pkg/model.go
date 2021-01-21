package pkg

import (
	"database/sql"
)

var (
	Username string
	Password string
	Ip       string
	Port     string
	Dbname   string
	Before   *basic
)

func NewBasic() *basic {
	return &basic{}
}

type basic struct {
	// mysql conn
	mysqlConn *sql.DB
	//basic info
	Hostname                           string
	Ips                                string
	Dbs                                string
	Var_binlog_format                  string
	Var_max_binlog_cache_size          string
	Var_max_binlog_size                int
	Var_max_connect_errors             string
	Var_max_connections                string
	Var_max_user_connections           string
	Var_open_files_limit               string
	Var_sync_binlog                    string
	Var_table_definition_cache         string
	Var_table_open_cache               string
	Var_thread_cache_size              string
	Var_innodb_adaptive_flushing       string
	Var_innodb_adaptive_hash_index     string
	Var_innodb_buffer_pool_size        int
	Var_innodb_file_per_table          string
	Var_innodb_flush_log_at_trx_commit string
	Var_innodb_flush_method            string
	Var_innodb_io_capacity             string
	Var_innodb_lock_wait_timeout       string
	Var_innodb_log_buffer_size         int
	Var_innodb_log_file_size           int
	Var_innodb_log_files_in_group      string
	Var_innodb_max_dirty_pages_pct     string
	Var_innodb_open_files              string
	Var_innodb_read_io_threads         string
	Var_innodb_thread_concurrency      string
	Var_innodb_write_io_threads        string

	//mysql -e "show global status" 不用\G
	//mysql -com
	Com_select   int
	Com_insert   int
	Com_update   int
	Com_delete   int
	Com_commit   int
	Com_rollback int
	//mysql -hit
	//while true;do s1=`mysql -e 'show global status'|grep -w -E 'Innodb_buffer_pool_read_requests|Innodb_buffer_pool_reads'|xargs echo|awk '{print $2}'`;s2=`mysql -e 'show global status'|grep -w -E 'Innodb_buffer_pool_read_requests|Innodb_buffer_pool_reads'|xargs echo|awk '{print $4}'`;sleep 1;ss1=`mysql -e 'show global status'|grep -w -E 'Innodb_buffer_pool_read_requests|Innodb_buffer_pool_reads'|xargs echo|awk '{print $2}'`;ss2=`mysql -e 'show global status'|grep -w -E 'Innodb_buffer_pool_read_requests|Innodb_buffer_pool_reads'|xargs echo|awk '{print $4}'`;rs1=$(($ss1-$s1+1));rs2=$(($ss2-$s2));rs3=$((1000000*($rs1-$rs2)/$rs1));echo $rs1,$rs2,$rs3;done
	// (Innodb_buffer_pool_read_requests - Innodb_buffer_pool_reads) / Innodb_buffer_pool_read_requests * 100%,每秒的计算
	Innodb_buffer_pool_read_requests int
	Innodb_buffer_pool_reads         int
	//mysql -innodb_rows
	Innodb_rows_inserted int
	Innodb_rows_updated  int
	Innodb_rows_deleted  int
	Innodb_rows_read     int
	//mysql -innodb_pages
	Innodb_buffer_pool_pages_data    int
	Innodb_buffer_pool_pages_free    int
	Innodb_buffer_pool_pages_dirty   int
	Innodb_buffer_pool_pages_flushed int
	//mysql --innodb_data
	Innodb_data_reads   int
	Innodb_data_writes  int
	Innodb_data_read    int
	Innodb_data_written int
	//mysql --innodb_log
	Innodb_os_log_fsyncs  int
	Innodb_os_log_written int
	//mysql --threads
	Threads_running   int
	Threads_connected int
	Threads_created   int
	Threads_cached    int
	//mysql --bytes
	Bytes_received int
	Bytes_sent     int
	//mysql --innodb_status show engine innodb status
	//log unflushed = Log sequence number - Log flushed up to
	//uncheckpointed bytes = Log sequence number - Last checkpoint at
	//mysql -e "show engine innodb status\G"|grep -n -E -A4 -B1 "^TRANSACTIONS|LOG|ROW OPERATIONS"
	//mysql -e "show engine innodb status\G"|grep -E "Last checkpoint|read view|queries inside|queue"
	Log_sequence    int
	Log_flushed     int
	History_list    int
	Last_checkpoint int
	Read_view       int
	Query_inside    int
	Query_queue     int
	//addition
	//show status
	Max_used_connections  int
	Aborted_connects      string
	Aborted_clients       string
	Select_full_join      string
	Binlog_cache_disk_use string
	Binlog_cache_use      string
	Opened_tables         string
	//Thread_cache_hits = (1 - Threads_created / connections ) * 100%
	Connections             int
	Qcache_hits             int
	Handler_read_first      int
	Handler_read_key        int
	Handler_read_next       int
	Handler_read_prev       int
	Handler_read_rnd        int
	Handler_read_rnd_next   int
	Handler_rollback        int
	Created_tmp_tables      int
	Created_tmp_disk_tables int
	Slow_queries            int
	Key_read_requests       int
	Key_reads               int
	Key_write_requests      int
	Key_writes              int
	Select_scan             int
	//半同步
	Rpl_semi_sync_master_net_avg_wait_time int
	Rpl_semi_sync_master_no_times          int
	Rpl_semi_sync_master_no_tx             int
	Rpl_semi_sync_master_status            string
	Rpl_semi_sync_master_tx_avg_wait_time  int
	Rpl_semi_sync_master_wait_sessions     int
	Rpl_semi_sync_master_yes_tx            int
	Rpl_semi_sync_slave_status             string
	rpl_semi_sync_master_timeout           string
	//Slave状态监控
	Master_Host           string
	Master_User           string
	Master_Port           string
	Slave_IO_Running      string
	Slave_SQL_Running     string
	Master_Server_Id      string
	Seconds_Behind_Master int
	Read_Master_Log_Pos   int
	Exec_Master_Log_Pos   int
}
