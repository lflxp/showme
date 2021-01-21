package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/host"
)

func (this *basic) GetHostAndIps() error {
	n, err := host.Info()
	if err != nil {
		return err
	}
	this.Hostname = n.Hostname
	this.Ips = strings.Join(GetIps(), ",")
	return nil
}

func (this *basic) InitMysqlConn() error {
	conn, err := MysqlConn(Username, Password, Ip, Port, Dbname)
	if err != nil {
		return err
	}
	this.mysqlConn = conn
	return nil
}

func (this *basic) CloseConn() error {
	err := this.mysqlConn.Close()
	return err
}

// show databases
func (this *basic) GetShowDatabases() error {
	var db string
	dbs := []string{}
	// todo show slave status

	rows, err := this.mysqlConn.Query("show databases")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&db)
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}

	this.Dbs = strings.Join(dbs, ",")
	return nil
}

// show global varibales
func (this *basic) GetShowGlobalVariables() error {
	var (
		variable_name string
		value         string
	)
	showGlobalVariables := map[string]string{}
	rows, err := this.mysqlConn.Query("show global variables")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&variable_name, &value)
		if err != nil {
			return err
		}
		// fmt.Println(variable_name, value)
		showGlobalVariables[variable_name] = value
	}

	if v, ok := showGlobalVariables["binlog_format"]; ok {
		this.Var_binlog_format = v
	}
	if v, ok := showGlobalVariables["max_binlog_cache_size"]; ok {
		this.Var_max_binlog_cache_size = v
	}
	if v, ok := showGlobalVariables["max_binlog_size"]; ok {
		this.Var_max_binlog_size, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalVariables["max_connect_errors"]; ok {
		this.Var_max_connect_errors = v
	}
	if v, ok := showGlobalVariables["max_connections"]; ok {
		this.Var_max_connections = v
	}
	if v, ok := showGlobalVariables["max_user_connections"]; ok {
		this.Var_max_user_connections = v
	}
	if v, ok := showGlobalVariables["open_files_limit"]; ok {
		this.Var_open_files_limit = v
	}
	if v, ok := showGlobalVariables["sync_binlog"]; ok {
		this.Var_sync_binlog = v
	}
	if v, ok := showGlobalVariables["table_definition_cache"]; ok {
		this.Var_table_definition_cache = v
	}
	if v, ok := showGlobalVariables["table_open_cache"]; ok {
		this.Var_table_open_cache = v
	}
	if v, ok := showGlobalVariables["thread_cache_size"]; ok {
		this.Var_thread_cache_size = v
	}
	if v, ok := showGlobalVariables["innodb_adaptive_flushing"]; ok {
		this.Var_innodb_adaptive_flushing = v
	}
	if v, ok := showGlobalVariables["innodb_adaptive_hash_index"]; ok {
		this.Var_innodb_adaptive_hash_index = v
	}
	if v, ok := showGlobalVariables["innodb_buffer_pool_size"]; ok {
		this.Var_innodb_buffer_pool_size, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalVariables["innodb_flush_log_at_trx_commit"]; ok {
		this.Var_innodb_flush_log_at_trx_commit = v
	}
	if v, ok := showGlobalVariables["innodb_flush_method"]; ok {
		this.Var_innodb_flush_method = v
	}
	if v, ok := showGlobalVariables["innodb_io_capacity"]; ok {
		this.Var_innodb_io_capacity = v
	}
	if v, ok := showGlobalVariables["innodb_lock_wait_timeout"]; ok {
		this.Var_innodb_lock_wait_timeout = v
	}
	if v, ok := showGlobalVariables["innodb_log_buffer_size"]; ok {
		this.Var_innodb_log_buffer_size, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalVariables["innodb_log_file_size"]; ok {
		this.Var_innodb_log_file_size, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalVariables["innodb_log_files_in_group"]; ok {
		this.Var_innodb_log_files_in_group = v
	}
	if v, ok := showGlobalVariables["innodb_max_dirty_pages_pct"]; ok {
		this.Var_innodb_max_dirty_pages_pct = v
	}
	if v, ok := showGlobalVariables["innodb_open_files"]; ok {
		this.Var_innodb_open_files = v
	}
	if v, ok := showGlobalVariables["innodb_read_io_threads"]; ok {
		this.Var_innodb_read_io_threads = v
	}
	if v, ok := showGlobalVariables["innodb_thread_concurrency"]; ok {
		this.Var_innodb_thread_concurrency = v
	}
	if v, ok := showGlobalVariables["innodb_write_io_threads"]; ok {
		this.Var_innodb_write_io_threads = v
	}

	return nil
}

// show variables
func (this *basic) GetShowVariables() error {
	var (
		variable_name string
		value         string
	)
	showVariables := map[string]string{}
	rows, err := this.mysqlConn.Query("show variables")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&variable_name, &value)
		if err != nil {
			return err
		}
		// fmt.Println(variable_name, value)
		showVariables[variable_name] = value
	}
	return nil
}

// show global status
func (this *basic) GetShowGlobalStatus() error {
	var (
		variable_name string
		value         string
	)
	showGlobalStatus := map[string]string{}
	rows, err := this.mysqlConn.Query("show global status")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&variable_name, &value)
		if err != nil {
			return err
		}
		// fmt.Println(variable_name, value)
		showGlobalStatus[variable_name] = value
	}
	//mysql -com
	if v, ok := showGlobalStatus["Com_select"]; ok {
		this.Com_select, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Com_insert"]; ok {
		this.Com_insert, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Com_delete"]; ok {
		this.Com_delete, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Com_update"]; ok {
		this.Com_update, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Com_commit"]; ok {
		this.Com_commit, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Com_rollback"]; ok {
		this.Com_rollback, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql -hit
	if v, ok := showGlobalStatus["Innodb_buffer_pool_read_requests"]; ok {
		this.Innodb_buffer_pool_read_requests, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_buffer_pool_reads"]; ok {
		this.Innodb_buffer_pool_reads, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql -innodb_rows
	if v, ok := showGlobalStatus["Innodb_rows_inserted"]; ok {
		this.Innodb_rows_inserted, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_rows_updated"]; ok {
		this.Innodb_rows_updated, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_rows_deleted"]; ok {
		this.Innodb_rows_deleted, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_rows_read"]; ok {
		this.Innodb_rows_read, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql -innodb_pages
	if v, ok := showGlobalStatus["Innodb_buffer_pool_pages_data"]; ok {
		this.Innodb_buffer_pool_pages_data, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_buffer_pool_pages_free"]; ok {
		this.Innodb_buffer_pool_pages_free, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_buffer_pool_pages_dirty"]; ok {
		this.Innodb_buffer_pool_pages_dirty, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_buffer_pool_pages_flushed"]; ok {
		this.Innodb_buffer_pool_pages_flushed, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql --innodb_data
	if v, ok := showGlobalStatus["Innodb_data_reads"]; ok {
		this.Innodb_data_reads, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_data_writes"]; ok {
		this.Innodb_data_writes, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_data_read"]; ok {
		this.Innodb_data_read, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_data_written"]; ok {
		this.Innodb_data_written, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql --innodb_log
	if v, ok := showGlobalStatus["Innodb_os_log_fsyncs"]; ok {
		this.Innodb_os_log_fsyncs, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Innodb_os_log_written"]; ok {
		this.Innodb_os_log_written, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql --threads
	if v, ok := showGlobalStatus["Threads_running"]; ok {
		this.Threads_running, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Threads_connected"]; ok {
		this.Threads_connected, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Threads_created"]; ok {
		this.Threads_created, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Threads_cached"]; ok {
		this.Threads_cached, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	//mysql --bytes
	if v, ok := showGlobalStatus["Bytes_received"]; ok {
		this.Bytes_received, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showGlobalStatus["Bytes_sent"]; ok {
		this.Bytes_sent, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// show status
func (this *basic) GetShowStatus() error {
	var (
		variable_name string
		value         string
	)
	showStatus := map[string]string{}
	rows, err := this.mysqlConn.Query("show status")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&variable_name, &value)
		if err != nil {
			return err
		}
		// fmt.Println(variable_name, value)
		showStatus[variable_name] = value
	}

	//show status
	if v, ok := showStatus["Max_used_connections"]; ok {
		this.Max_used_connections, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Aborted_connects"]; ok {
		this.Aborted_connects = v
	}
	if v, ok := showStatus["Aborted_clients"]; ok {
		this.Aborted_clients = v
	}
	if v, ok := showStatus["Select_full_join"]; ok {
		this.Select_full_join = v
	}
	if v, ok := showStatus["Binlog_cache_disk_use"]; ok {
		this.Binlog_cache_disk_use = v
	}
	if v, ok := showStatus["Binlog_cache_use"]; ok {
		this.Binlog_cache_use = v
	}
	if v, ok := showStatus["Opened_tables"]; ok {
		this.Opened_tables = v
	}
	//Thread_cache_hits = (1 - Threads_created / connections ) * 100%
	if v, ok := showStatus["Connections"]; ok {
		this.Connections, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Qcache_hits"]; ok {
		this.Qcache_hits, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_first"]; ok {
		this.Handler_read_first, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_key"]; ok {
		this.Handler_read_key, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_next"]; ok {
		this.Handler_read_next, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_prev"]; ok {
		this.Handler_read_prev, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_rnd"]; ok {
		this.Handler_read_rnd, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_read_rnd_next"]; ok {
		this.Handler_read_rnd_next, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Handler_rollback"]; ok {
		this.Handler_rollback, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Created_tmp_tables"]; ok {
		this.Created_tmp_tables, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Created_tmp_disk_tables"]; ok {
		this.Created_tmp_disk_tables, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Slow_queries"]; ok {
		this.Slow_queries, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Key_read_requests"]; ok {
		this.Key_read_requests, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Key_reads"]; ok {
		this.Key_reads, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Key_write_requests"]; ok {
		this.Key_write_requests, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Key_writes"]; ok {
		this.Key_writes, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showStatus["Select_scan"]; ok {
		this.Select_scan, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// todo show processlist
func (this *basic) GetShowProcesslist() error {
	var (
		id, user, host, db, command, time, state, info string
	)
	// showSlaveStatus := map[string]string{}
	rows, err := this.mysqlConn.Query("show processlist")
	if err != nil {
		return err
	}

	num := 0
	for rows.Next() {
		// rs, err := rows.Columns()
		// if err != nil {
		// 	return err
		// }
		// fmt.Println("Columnes", rs)
		if num == 0 {
			fmt.Println("Id User Host db Command Time State Info")
			num++
		}
		err = rows.Scan(&id, &user, &host, &db, &command, &time, &state, &info)
		if err != nil {
			return err
		}
		fmt.Printf("%s %s %s %s %s %s %s %s\n", id, user, host, db, command, time, state, info)
		// showSlaveStatus[variable_name] = value
	}

	return nil
}

// todo show slave status
func (this *basic) GetShowSlaveStatus() error {
	var (
		variable_name string
		value         string
	)
	showSlaveStatus := map[string]string{}
	rows, err := this.mysqlConn.Query("show slave status")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&variable_name, &value)
		if err != nil {
			return err
		}
		fmt.Println(variable_name, value)
		showSlaveStatus[variable_name] = value
	}
	//半同步
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_net_avg_wait_time"]; ok {
		this.Rpl_semi_sync_master_net_avg_wait_time, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_no_times"]; ok {
		this.Rpl_semi_sync_master_no_times, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_no_tx"]; ok {
		this.Rpl_semi_sync_master_no_tx, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_status"]; ok {
		this.Rpl_semi_sync_master_status = v
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_tx_avg_wait_time"]; ok {
		this.Rpl_semi_sync_master_tx_avg_wait_time, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_wait_sessions"]; ok {
		this.Rpl_semi_sync_master_wait_sessions, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_master_yes_tx"]; ok {
		this.Rpl_semi_sync_master_yes_tx, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Rpl_semi_sync_slave_status"]; ok {
		this.Rpl_semi_sync_slave_status = v
	}
	if v, ok := showSlaveStatus["rpl_semi_sync_master_timeout"]; ok {
		this.rpl_semi_sync_master_timeout = v
	}
	//Slave状态监控

	if v, ok := showSlaveStatus["Master_Host"]; ok {
		this.Master_Host = v
	}
	if v, ok := showSlaveStatus["Master_User"]; ok {
		this.Master_User = v
	}
	if v, ok := showSlaveStatus["Master_Port"]; ok {
		this.Master_Port = v
	}
	if v, ok := showSlaveStatus["Slave_IO_Running"]; ok {
		this.Slave_IO_Running = v
	}
	if v, ok := showSlaveStatus["Slave_SQL_Running"]; ok {
		this.Slave_SQL_Running = v
	}
	if v, ok := showSlaveStatus["Master_Server_Id"]; ok {
		this.Master_Server_Id = v
	}
	if v, ok := showSlaveStatus["Seconds_Behind_Master"]; ok {
		this.Seconds_Behind_Master, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Read_Master_Log_Pos"]; ok {
		this.Read_Master_Log_Pos, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if v, ok := showSlaveStatus["Exec_Master_Log_Pos"]; ok {
		this.Exec_Master_Log_Pos, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// todo mmp value to struct
// show engine innodb status
func (this *basic) GetShowEngineInnodbStatus() error {
	var (
		types  string
		name   string
		status string
	)
	rows, err := this.mysqlConn.Query("show engine innodb status")
	if err != nil {
		return err
	}

	for rows.Next() {
		err = rows.Scan(&types, &name, &status)
		if err != nil {
			return err
		}
		fmt.Println(types, name, status)
	}

	for _, x := range strings.Split(status, "\n") {
		if strings.Contains(x, "Log sequence number") {
			tmp := strings.Split(strings.TrimSpace(x), " ")
			this.Log_sequence, err = strconv.Atoi(tmp[len(tmp)-1])
			if err != nil {
				return err
			}
		}
		if strings.Contains(x, "Log flushed up to") {
			tmp := strings.Split(strings.TrimSpace(x), " ")
			this.Log_flushed, err = strconv.Atoi(tmp[len(tmp)-1])
			if err != nil {
				return err
			}
		}
		if strings.Contains(x, "History list length") {
			tmp := strings.Split(strings.TrimSpace(x), " ")
			this.History_list, err = strconv.Atoi(tmp[len(tmp)-1])
			if err != nil {
				return err
			}
		}
		if strings.Contains(x, "Last checkpoint at") {
			tmp := strings.Split(strings.TrimSpace(x), " ")
			this.Last_checkpoint, err = strconv.Atoi(tmp[len(tmp)-1])
			if err != nil {
				return err
			}
		}
		if strings.Contains(x, "read views open") {
			tmp := strings.Split(strings.TrimSpace(x), " ")
			this.Read_view, err = strconv.Atoi(tmp[0])
			if err != nil {
				return err
			}
		}
		if strings.Contains(x, "queries inside") {
			tmp := strings.Split(strings.TrimSpace(x), ",")
			this.Query_inside, err = strconv.Atoi(strings.Split(strings.TrimSpace(tmp[0]), " ")[0])
			if err != nil {
				return err
			}

			this.Query_queue, err = strconv.Atoi(strings.Split(strings.TrimSpace(tmp[1]), " ")[0])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (this *basic) CreateCom(interval int) string {
	data_detail := ""
	insert_diff := (this.Com_insert - Before.Com_insert) / interval
	update_diff := (this.Com_update - Before.Com_update) / interval
	delete_diff := (this.Com_delete - Before.Com_delete) / interval
	select_diff := (this.Com_select - Before.Com_select) / interval
	commit_diff := (this.Com_commit - Before.Com_commit) / interval
	rollback_diff := (this.Com_rollback - Before.Com_rollback) / interval
	tps := rollback_diff + commit_diff

	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(insert_diff)))+strconv.Itoa(insert_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(update_diff)))+strconv.Itoa(update_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(delete_diff)))+strconv.Itoa(delete_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(select_diff)))+strconv.Itoa(select_diff), "yellow", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(tps)))+strconv.Itoa(tps), "yellow", "", false, false)
	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateHit(interval int) string {
	data_detail := ""
	read_request := (this.Innodb_buffer_pool_read_requests - Before.Innodb_buffer_pool_read_requests) / interval
	read := (this.Innodb_buffer_pool_reads - Before.Innodb_buffer_pool_reads) / interval
	key_read := (this.Key_reads - Before.Key_reads) / interval
	key_write := (this.Key_writes - Before.Key_writes) / interval
	key_read_request := (this.Key_read_requests - Before.Key_read_requests) / interval
	key_write_request := (this.Key_write_requests - Before.Key_write_requests) / interval
	//innodb hit
	hrr := (this.Handler_read_rnd - Before.Handler_read_rnd) / interval
	hrrn := (this.Handler_read_rnd_next - Before.Handler_read_rnd_next) / interval
	hrf := (this.Handler_read_first - Before.Handler_read_first) / interval
	hrk := (this.Handler_read_key - Before.Handler_read_key) / interval
	hrn := (this.Handler_read_next - Before.Handler_read_next) / interval
	hrp := (this.Handler_read_prev - Before.Handler_read_prev) / interval
	//key buffer read hit
	key_read_hit := (float64(key_read_request-key_read) + 0.0001) / (float64(key_read_request) + 0.0001) * 100
	key_write_hit := (float64(key_write_request-key_write) + 0.0001) / (float64(key_write_request) + 0.0001) * 100
	index_total_hit := (100 - (100 * (float64(this.Handler_read_rnd+this.Handler_read_rnd_next) + 0.0001) / (0.0001 + float64(this.Handler_read_first+this.Handler_read_key+this.Handler_read_next+this.Handler_read_prev+this.Handler_read_rnd+this.Handler_read_rnd_next))))
	index_current_hit := 100.00
	if hrr+hrrn != 0 {
		index_current_hit = (100 - (100 * (float64(hrr+hrrn) + 0.0001) / (0.0001 + float64(hrf+hrk+hrn+hrp+hrr+hrrn))))
	}
	query_hits_s := (this.Qcache_hits - Before.Qcache_hits) / interval
	com_select_s := (this.Com_select - Before.Com_select) / interval

	query_hit := (float64(query_hits_s) + 0.0001) / (float64(query_hits_s+com_select_s) + 0.0001) * 100
	//innodb_hit
	innodb_hit := ((float64(read_request-read) + 0.0001) / (float64(read_request) + 0.0001)) * 100

	data_detail += hit(6, key_read_hit)
	data_detail += hit(7, key_write_hit)
	data_detail += hit(7, index_current_hit)
	data_detail += hit(7, index_total_hit)
	data_detail += hit(7, query_hit)
	// lor = read_request
	data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(read_request)))+strconv.Itoa(read_request), "", "", false, false)

	data_detail += hit(7, innodb_hit)

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateInnodbRows(interval int) string {
	data_detail := ""
	innodb_rows_inserted_diff := (this.Innodb_rows_inserted - Before.Innodb_rows_inserted) / interval
	innodb_rows_updated_diff := (this.Innodb_rows_updated - Before.Innodb_rows_updated) / interval
	innodb_rows_deleted_diff := (this.Innodb_rows_deleted - Before.Innodb_rows_deleted) / interval
	innodb_rows_read_diff := (this.Innodb_rows_read - Before.Innodb_rows_read) / interval

	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(innodb_rows_inserted_diff)))+strconv.Itoa(innodb_rows_inserted_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(innodb_rows_updated_diff)))+strconv.Itoa(innodb_rows_updated_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(innodb_rows_deleted_diff)))+strconv.Itoa(innodb_rows_deleted_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(innodb_rows_read_diff)))+strconv.Itoa(innodb_rows_read_diff), "", "", false, false)

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateInnodbPages(interval int) string {
	data_detail := ""
	flush := (this.Innodb_buffer_pool_pages_flushed - Before.Innodb_buffer_pool_pages_flushed) / interval

	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(this.Innodb_buffer_pool_pages_data)))+strconv.Itoa(Before.Innodb_buffer_pool_pages_data), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(this.Innodb_buffer_pool_pages_free)))+strconv.Itoa(Before.Innodb_buffer_pool_pages_free), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(this.Innodb_buffer_pool_pages_dirty)))+strconv.Itoa(Before.Innodb_buffer_pool_pages_dirty), "yellow", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(flush)))+strconv.Itoa(flush), "yellow", "", false, false)

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateInnodbData(interval int) string {
	data_detail := ""
	innodb_data_reads_diff := (this.Innodb_data_reads - Before.Innodb_data_reads) / interval
	innodb_data_writes_diff := (this.Innodb_data_writes - Before.Innodb_data_writes) / interval
	innodb_data_read_diff := (this.Innodb_data_read - Before.Innodb_data_read) / interval
	innodb_data_written_diff := (this.Innodb_data_written - Before.Innodb_data_written) / interval

	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(innodb_data_reads_diff)))+strconv.Itoa(innodb_data_reads_diff), "", "", false, false)
	data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(innodb_data_writes_diff)))+strconv.Itoa(innodb_data_writes_diff), "", "", false, false)

	if innodb_data_read_diff/1024/1024 > 9 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(innodb_data_read_diff)/1024/1024, 1)))+floatToString(float64(innodb_data_read_diff)/1024/1024, 1)+"m", "red", "", false, true)
	} else if innodb_data_read_diff/1024/1024 <= 9 && innodb_data_read_diff/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(innodb_data_read_diff)/1024/1024, 1)))+floatToString(float64(innodb_data_read_diff)/1024/1024, 1)+"m", "", "", false, false)
	} else if innodb_data_read_diff/1024 >= 1 && innodb_data_read_diff/1024/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(innodb_data_read_diff/1024)))+strconv.Itoa(innodb_data_read_diff/1024)+"k", "", "", false, false)
	} else if innodb_data_read_diff/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(innodb_data_read_diff)))+strconv.Itoa(innodb_data_read_diff), "", "", false, false)
	}

	if innodb_data_written_diff/1024/1024 > 9 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(float64(innodb_data_written_diff)/1024/1024, 1)))+floatToString(float64(innodb_data_written_diff)/1024/1024, 1)+"m", "red", "", false, true)
	} else if innodb_data_written_diff/1024/1024 <= 9 && innodb_data_written_diff/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(float64(innodb_data_written_diff)/1024/1024, 1)))+floatToString(float64(innodb_data_written_diff)/1024/1024, 1)+"m", "", "", false, false)
	} else if innodb_data_written_diff/1024 >= 1 && innodb_data_written_diff/1024/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(innodb_data_written_diff/1024)))+strconv.Itoa(innodb_data_written_diff/1024)+"k", "", "", false, false)
	} else if innodb_data_written_diff/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(innodb_data_written_diff)))+strconv.Itoa(innodb_data_written_diff), "", "", false, false)
	}

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateInnodbLog(interval int) string {
	data_detail := ""
	innodb_os_log_fsyncs_diff := (this.Innodb_os_log_fsyncs - Before.Innodb_os_log_fsyncs) / interval
	innodb_os_log_written_diff := (this.Innodb_os_log_written - Before.Innodb_os_log_written) / interval

	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(innodb_os_log_fsyncs_diff)))+strconv.Itoa(innodb_os_log_fsyncs_diff), "", "", false, false)

	if innodb_os_log_written_diff/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(float64(innodb_os_log_written_diff)/1024/1024, 1)))+floatToString(float64(innodb_os_log_written_diff)/1024/1024, 1)+"m", "red", "", false, true)
	} else if innodb_os_log_written_diff/1024/1024 < 1 && innodb_os_log_written_diff/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(int(float64(innodb_os_log_written_diff)/1024/1024+0.5))))+strconv.Itoa(int(float64(innodb_os_log_written_diff)/1024/1024+0.5))+"k", "yellow", "", false, false)
	} else if innodb_os_log_written_diff/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(innodb_os_log_written_diff)))+strconv.Itoa(innodb_os_log_written_diff), "", "", false, false)
	}

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateInnodbStatus(interval int) string {
	data_detail := ""
	unflushed_log := this.Log_sequence - this.Log_flushed
	uncheckpointed_bytes := this.Log_sequence - this.Last_checkpoint
	//History_list
	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.History_list)))+strconv.Itoa(this.History_list), "", "", false, false)
	//unflushed_log
	if unflushed_log/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(unflushed_log)/1024/1024+0.5, 1)))+floatToString(float64(unflushed_log)/1024/1024+0.5, 1)+"m", "yellow", "", false, false)
	} else if unflushed_log/1024/1024 < 1 && unflushed_log/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(float64(unflushed_log)/1024+0.5))))+strconv.Itoa(int(float64(unflushed_log)/1024+0.5))+"k", "yellow", "", false, false)
	} else if unflushed_log/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(unflushed_log)))+strconv.Itoa(unflushed_log), "yellow", "", false, false)
	}

	//uncheckpointed_bytes
	if uncheckpointed_bytes/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(uncheckpointed_bytes)/1024/1024+0.5, 1)))+floatToString(float64(uncheckpointed_bytes)/1024/1024+0.5, 1)+"m", "yellow", "", false, false)
	} else if uncheckpointed_bytes/1024/1024 < 1 && uncheckpointed_bytes/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(float64(uncheckpointed_bytes)/1024+0.5))))+strconv.Itoa(int(float64(uncheckpointed_bytes)/1024+0.5))+"k", "yellow", "", false, false)
	} else if uncheckpointed_bytes/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(uncheckpointed_bytes)))+strconv.Itoa(uncheckpointed_bytes), "yellow", "", false, false)
	}

	//Read_views
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(this.Read_view)))+strconv.Itoa(this.Read_view), "", "", false, false)
	//inside
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(this.Query_inside)))+strconv.Itoa(this.Query_inside), "", "", false, false)
	//queue
	data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(this.Query_queue)))+strconv.Itoa(this.Query_queue), "", "", false, false)

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateThreads(interval int) string {
	data_detail := ""
	connections_dirr := (this.Connections - Before.Connections) / interval

	threads_created_diff := (this.Threads_created - Before.Threads_created) / interval

	thread_cache_hit := (1 - float64(threads_created_diff)/float64(connections_dirr)) * 100

	data_detail += Colorize(strings.Repeat(" ", 4-len(strconv.Itoa(this.Threads_running)))+strconv.Itoa(this.Threads_running), "", "", false, false)

	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.Threads_connected)))+strconv.Itoa(this.Threads_connected), "", "", false, false)

	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(threads_created_diff)))+strconv.Itoa(threads_created_diff), "", "", false, false)

	data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.Threads_cached)))+strconv.Itoa(this.Threads_cached), "", "", false, false)
	if thread_cache_hit > 99.0 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(thread_cache_hit, 2)))+floatToString(thread_cache_hit, 2), "green", "", false, false)
	} else if thread_cache_hit <= 99.0 && thread_cache_hit > 90.0 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(thread_cache_hit, 2)))+floatToString(thread_cache_hit, 2), "yellow", "", false, false)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 7-len(floatToString(thread_cache_hit, 2)))+floatToString(thread_cache_hit, 2), "red", "", false, false)
	}

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateBytes(interval int) string {
	data_detail := ""
	bytes_received_diff := (this.Bytes_received - Before.Bytes_received) / interval
	bytes_sent_diff := (this.Bytes_sent - Before.Bytes_sent) / interval

	if bytes_received_diff/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(bytes_received_diff)/1024/1024+0.5, 1)))+floatToString(float64(bytes_received_diff)/1024/1024+0.5, 1)+"m", "red", "", false, true)
	} else if bytes_received_diff/1024/1024 < 1 && bytes_received_diff/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(float64(bytes_received_diff)/1024+0.5))))+strconv.Itoa(int(float64(bytes_received_diff)/1024+0.5))+"k", "", "", false, false)
	} else if bytes_received_diff/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(bytes_received_diff)))+strconv.Itoa(bytes_received_diff), "", "", false, false)
	}

	if bytes_sent_diff/1024/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(floatToString(float64(bytes_sent_diff)/1024/1024+0.5, 1)))+floatToString(float64(bytes_sent_diff)/1024/1024+0.5, 1)+"m", "red", "", false, true)
	} else if bytes_sent_diff/1024/1024 < 1 && bytes_sent_diff/1024 >= 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(int(float64(bytes_sent_diff)/1024+0.5))))+strconv.Itoa(int(float64(bytes_sent_diff)/1024+0.5))+"k", "", "", false, false)
	} else if bytes_sent_diff/1024 < 1 {
		data_detail += Colorize(strings.Repeat(" ", 7-len(strconv.Itoa(bytes_sent_diff)))+strconv.Itoa(bytes_sent_diff), "", "", false, false)
	}

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateSemi(interval int) string {
	data_detail := ""
	if this.Rpl_semi_sync_master_net_avg_wait_time < 1000 {
		data_detail += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time)))+strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time)+"us", "", "", false, false)
	} else if this.Rpl_semi_sync_master_net_avg_wait_time >= 1000 && this.Rpl_semi_sync_master_net_avg_wait_time/1000/1000 <= 1 {
		data_detail += Colorize(strings.Repeat(" ", 3-len(strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time/1000)))+strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time/1000)+"ms", "", "", false, false)
	} else if this.Rpl_semi_sync_master_net_avg_wait_time/1000/1000 > 1 {
		data_detail += Colorize(strings.Repeat(" ", 4-len(strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time/1000/1000)))+strconv.Itoa(this.Rpl_semi_sync_master_net_avg_wait_time/1000/1000)+"s", "red", "", false, true)
	}

	if this.Rpl_semi_sync_master_tx_avg_wait_time < 1000 {
		data_detail += Colorize(strings.Repeat(" ", 4-len(strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time)))+strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time)+"us", "", "", false, false)
	} else if this.Rpl_semi_sync_master_tx_avg_wait_time > 1000 && this.Rpl_semi_sync_master_tx_avg_wait_time/1000/1000 <= 1 {
		data_detail += Colorize(strings.Repeat(" ", 4-len(strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time/1000)))+strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time/1000)+"ms", "", "", false, false)
	} else if this.Rpl_semi_sync_master_tx_avg_wait_time/1000/1000 > 1 {
		data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time/1000/1000)))+strconv.Itoa(this.Rpl_semi_sync_master_tx_avg_wait_time/1000/1000)+"s", "red", "", false, true)
	}

	if this.Rpl_semi_sync_master_no_tx > 1 {
		data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.Rpl_semi_sync_master_no_tx)))+strconv.Itoa(this.Rpl_semi_sync_master_no_tx), "red", "", false, true)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 5-len(strconv.Itoa(this.Rpl_semi_sync_master_no_tx)))+strconv.Itoa(this.Rpl_semi_sync_master_no_tx), "", "", false, true)
	}

	data_detail += Colorize(strings.Repeat(" ", 5-len(changeUntils(this.Rpl_semi_sync_master_yes_tx)))+changeUntils(this.Rpl_semi_sync_master_yes_tx), "", "", false, true)

	if this.Rpl_semi_sync_master_no_times > 1 {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(this.Rpl_semi_sync_master_no_times)))+strconv.Itoa(this.Rpl_semi_sync_master_no_times), "red", "", false, true)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 6-len(strconv.Itoa(this.Rpl_semi_sync_master_no_times)))+strconv.Itoa(this.Rpl_semi_sync_master_no_times), "", "", false, true)
	}
	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func (this *basic) CreateSlave(interval int) string {
	data_detail := ""
	checkNum := this.Read_Master_Log_Pos - this.Exec_Master_Log_Pos

	data_detail += Colorize(strings.Repeat(" ", 11-len(strconv.Itoa(this.Read_Master_Log_Pos)))+strconv.Itoa(this.Read_Master_Log_Pos), "", "", false, false)

	data_detail += Colorize(strings.Repeat(" ", 12-len(strconv.Itoa(this.Exec_Master_Log_Pos)))+strconv.Itoa(this.Exec_Master_Log_Pos), "", "", false, false)

	data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(checkNum)))+strconv.Itoa(checkNum), "", "", false, false)

	if this.Seconds_Behind_Master > 300 {
		data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(this.Seconds_Behind_Master)))+strconv.Itoa(this.Seconds_Behind_Master), "red", "", false, false)
	} else {
		data_detail += Colorize(strings.Repeat(" ", 8-len(strconv.Itoa(this.Seconds_Behind_Master)))+strconv.Itoa(this.Seconds_Behind_Master), "green", "", false, false)
	}

	data_detail += Colorize("|", "green", "", false, false)
	return data_detail
}

func floatToString(x float64, f int) string {
	rs := strconv.FormatFloat(x, 'f', f, 64)
	return rs
}

func perSecond_Float(before float64, after float64, time string) (string, bool) {
	var result interface{}
	var ok bool
	var rs string
	//转换时间为float64
	seconds, err := strconv.ParseFloat(time, 64)
	if err != nil {
		fmt.Println(err)
	}
	result = (after - before) / seconds
	// tmp = fmt.Sprintf("%s", reflect.TypeOf(result))
	switch result.(type) {
	case int:
		// fmt.Println("int")
		rs = strconv.Itoa(result.(int))
	case int64:
		fmt.Println("int64")
		rs = strconv.FormatInt(result.(int64), 64)
	// case float32:
	//  fmt.Println("float32")
	//  rs = strconv.FormatFloat(result.(float32), 'f', 4, 32)
	case float64:
		fmt.Println("float64")
		rs = strconv.FormatFloat(result.(float64), 'f', 4, 64)
	default:
		panic("not fount number type in perSecond_Float")
	}

	if result.(int) > 0 {
		ok = true
	} else {
		ok = false
	}

	return rs, ok
}

// 单位转换
func changeUntils(in int) string {
	var result string
	if in/1024 < 1 {
		tmp := strconv.Itoa(in)
		result = tmp
	} else if in/1024 >= 1 && in/1024/1024 < 1 {
		tmp := strconv.Itoa(in / 1024)
		result = tmp + "k"
	} else if in/1024/1024 >= 1 && in/1024/1024/1024 < 1 {
		tmp := strconv.Itoa(in / 1024 / 1024)
		result = tmp + "m"
	} else if in/1024/1024/1024 >= 1 && in/1024/1024/1024/1024 < 1 {
		tmp := strconv.Itoa(in / 1024 / 1024 / 1024)
		result = tmp + "g"
	} else if in/1024/1024/1024/1024 >= 1 {
		tmp := strconv.Itoa(in / 1024 / 1024 / 1024 / 1024)
		result = tmp + "pg"
	}
	return result
}

// 阈值颜色
func hit(num int, in float64) string {
	var result string
	if in > 99.0 {
		result = Colorize(strings.Repeat(" ", num-len(floatToString(in, 2)))+floatToString(in, 2), "green", "", false, false)
	} else if in > 90.0 && in <= 99.0 {
		result = Colorize(strings.Repeat(" ", num-len(floatToString(in, 2)))+floatToString(in, 2), "yellow", "", false, false)
	} else if in < 0.01 {
		result = Colorize(strings.Repeat(" ", num-len("100.00"))+"100.00", "green", "", false, false)
	} else {
		result = Colorize(strings.Repeat(" ", num-len(floatToString(in, 2)))+floatToString(in, 2), "red", "", false, true)
	}
	return result
}
