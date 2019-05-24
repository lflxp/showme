package mysql

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/lflxp/showme/utils"
)

func BeforeRun(in string) error {
	var mysql *basic
	if in != "mysql" {
		mysql = NewBasic()
		// parse input -u -p -P -H -db
		inputs := strings.Split(strings.TrimSpace(in), " ")
		for n, x := range inputs {
			if x == "-u" {
				if n == len(inputs)-1 {
					return errors.New("some args no given value")
				}
				Username = inputs[n+1]
			} else if x == "-P" {
				if n == len(inputs)-1 {
					return errors.New("some args no given value")
				}
				Port = inputs[n+1]
			} else if x == "-db" {
				if n == len(inputs)-1 {
					return errors.New("some args no given value")
				}
				Dbname = inputs[n+1]
			} else if x == "-p" {
				if n == len(inputs)-1 {
					return errors.New("some args no given value")
				}
				Password = inputs[n+1]
			} else if x == "-H" {
				if n == len(inputs)-1 {
					return errors.New("some args no given value")
				}
				Ip = inputs[n+1]
			}
		}

		if Username == "" {
			Username = "root"
		}
		if Port == "" {
			Port = "3306"
		}
		if Dbname == "" {
			Dbname = "mysql"
		}
		if Ip == "" {
			Ip = "127.0.0.1"
		}
		if !strings.Contains(in, "-p") {
			Password = utils.Prompt(fmt.Sprintf("Please input user: %s password", Username))
		}
		fmt.Printf("%s *** %s %s %s \n", Username, Port, Ip, Dbname)
	} else {
		return errors.New("nothing input")
	}
	if in == "mysql test GetHostAndIps" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetHostAndIps()
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Hostname %s\nIps %s", mysql.Hostname, mysql.Ips))
	} else if in == "mysql test GetShowDatabases" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowDatabases()
		if err != nil {
			return err
		}
		fmt.Printf("Dbs %s\n", mysql.Dbs)
	} else if in == "mysql test GetShowGlobalVariables" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowGlobalVariables()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowVariables" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowVariables()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowGlobalStatus" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowGlobalStatus()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowStatus" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowStatus()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowEngineInnodbStatus" {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowEngineInnodbStatus()
		if err != nil {
			return err
		}
	} else if strings.Contains(in, "mysql processlist") {
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
		err = mysql.GetShowProcesslist()
		if err != nil {
			return err
		}
	} else if in == "mysql status" {
		fmt.Println("mysql status todo")
	} else {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		// 获取退出信号
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		ok := true

		interval := 1
		num := 0

		// 主机信息
		for _, x := range utils.CollectEasy() {
			fmt.Println(x)
		}

		for {
			select {
			case s := <-c:
				fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
				ok = false
				break
			case <-t.C:
				tmp := NewBasic()
				err := tmp.InitMysqlConn()
				if err != nil {
					panic(err)
					return err
				}
				// defer tmp.CloseConn()

				tmp.GetHostAndIps()
				tmp.GetShowDatabases()
				tmp.GetShowGlobalStatus()
				tmp.GetShowGlobalVariables()
				tmp.GetShowStatus()
				tmp.GetShowVariables()

				if num == 0 {
					var tmptable string
					tmp_table_x := float64(tmp.Created_tmp_disk_tables) / float64(tmp.Created_tmp_tables) * 100
					if tmp_table_x < 10.0 {
						tmptable = utils.Colorize(floatToString(tmp_table_x, 2), "green", "", false, false)
					} else {
						tmptable = utils.Colorize(floatToString(tmp_table_x, 2), "red", "", false, true)
					}
					fmt.Printf("%s: %s \n", utils.Colorize("        DB        ", "white", "red", true, true), utils.Colorize(tmp.Dbs, "yellow", "", false, false))
					fmt.Printf("%s: %s \n", utils.Colorize("        Var       ", "white", "red", true, true), utils.Colorize("binlog_format", "purple", "", false, false)+"["+tmp.Var_binlog_format+"]"+utils.Colorize(" max_binlog_cache_size", "purple", "", false, false)+"["+tmp.Var_max_binlog_cache_size+"]"+utils.Colorize(" max_binlog_size", "purple", "", false, false)+"["+changeUntils(tmp.Var_max_binlog_size)+"]"+utils.Colorize(" sync_binlog", "purple", "", false, false)+"["+tmp.Var_sync_binlog+"]")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("max_connect_errors", "purple", "", false, false)+"["+tmp.Var_max_connect_errors+"]"+utils.Colorize(" max_connections", "purple", "", false, false)+"["+tmp.Var_max_connections+"]"+utils.Colorize(" max_user_connections", "purple", "", false, false)+"["+tmp.Var_max_user_connections+"]"+utils.Colorize(" max_used_connections", "purple", "", false, false)+"["+string(tmp.Max_used_connections)+"]")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("open_files_limit", "purple", "", false, false)+"["+tmp.Var_open_files_limit+"]"+utils.Colorize(" table_definition_cache", "purple", "", false, false)+"["+tmp.Var_table_definition_cache+"]"+utils.Colorize(" Aborted_connects", "purple", "", false, false)+"["+tmp.Aborted_connects+"]"+utils.Colorize(" Aborted_clients", "purple", "", false, false)+"["+tmp.Aborted_clients+"]")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("Binlog_cache_disk_use", "purple", "", false, false)+"["+tmp.Binlog_cache_disk_use+"]"+utils.Colorize(" Select_scan", "purple", "", false, false)+"["+string(tmp.Select_scan)+"]"+utils.Colorize(" Select_full_join", "purple", "", false, false)+"["+tmp.Select_full_join+"]"+utils.Colorize(" Slow_queries", "purple", "", false, false)+"["+string(tmp.Slow_queries)+"]")
					if tmp.Rpl_semi_sync_master_status != "" {
						fmt.Printf("%s  %s\n", "                  ", utils.Colorize("Rpl_semi_sync_master_status", "purple", "", false, false)+"["+tmp.Rpl_semi_sync_master_status+"]"+utils.Colorize(" Rpl_semi_sync_slave_status", "purple", "", false, false)+"["+tmp.Rpl_semi_sync_slave_status+"]"+utils.Colorize(" rpl_semi_sync_master_timeout", "purple", "", false, false)+"["+tmp.rpl_semi_sync_master_timeout+"]")
					}
					if tmp.Master_Host != "" {
						fmt.Printf("%s  %s\n", "                  ", utils.Colorize("Master_Host", "purple", "", false, false)+"["+tmp.Master_Host+"]"+utils.Colorize(" Master_User", "purple", "", false, false)+"["+tmp.Master_User+"]"+utils.Colorize(" Master_Port", "purple", "", false, false)+"["+tmp.Master_Port+"]"+utils.Colorize(" Master_Server_Id", "purple", "", false, false)+"["+tmp.Master_Server_Id+"]")
						io := ""
						sql := ""
						if tmp.Slave_IO_Running != "Yes" {
							io = utils.Colorize("No", "red", "", false, true)
						} else {
							io = utils.Colorize("Yes", "green", "", false, false)
						}
						if tmp.Slave_SQL_Running != "Yes" {
							sql = utils.Colorize("No", "red", "", false, true)
						} else {
							sql = utils.Colorize("Yes", "green", "", false, false)
						}
						fmt.Printf("%s  %s\n", "                  ", utils.Colorize("Slave_IO_Running", "purple", "", false, false)+"["+io+"]"+utils.Colorize(" Slave_SQL_Running", "purple", "", false, false)+"["+sql+"]\n")
					}
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("table_open_cache", "purple", "", false, false)+"["+tmp.Var_table_open_cache+"]"+utils.Colorize(" thread_cache_size", "purple", "", false, false)+"["+tmp.Var_thread_cache_size+"]"+utils.Colorize(" Opened_tables", "purple", "", false, false)+"["+tmp.Opened_tables+"]"+utils.Colorize(" Created_tmp_disk_tables_ratio", "purple", "", false, false)+"["+tmptable+"]")

					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_adaptive_flushing", "purple", "", false, false)+"["+tmp.Var_innodb_adaptive_flushing+"]"+utils.Colorize(" innodb_adaptive_hash_index", "purple", "", false, false)+"["+tmp.Var_innodb_adaptive_hash_index+"]"+utils.Colorize(" innodb_buffer_pool_size", "purple", "", false, false)+"["+changeUntils(tmp.Var_innodb_buffer_pool_size)+"]"+"")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_file_per_table", "purple", "", false, false)+"["+tmp.Var_innodb_file_per_table+"]"+utils.Colorize(" innodb_flush_log_at_trx_commit", "purple", "", false, false)+"["+tmp.Var_innodb_flush_log_at_trx_commit+"]"+utils.Colorize(" innodb_flush_method", "purple", "", false, false)+"["+tmp.Var_innodb_flush_method+"]"+"")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_io_capacity", "purple", "", false, false)+"["+tmp.Var_innodb_io_capacity+"]"+utils.Colorize(" innodb_lock_wait_timeout", "purple", "", false, false)+"["+tmp.Var_innodb_lock_wait_timeout+"]"+utils.Colorize(" innodb_log_buffer_size", "purple", "", false, false)+"["+changeUntils(tmp.Var_innodb_log_buffer_size)+"]"+"")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_log_file_size", "purple", "", false, false)+"["+changeUntils(tmp.Var_innodb_log_file_size)+"]"+utils.Colorize(" innodb_log_files_in_group", "purple", "", false, false)+"["+tmp.Var_innodb_log_files_in_group+"]"+utils.Colorize(" innodb_max_dirty_pages_pct", "purple", "", false, false)+"["+tmp.Var_innodb_max_dirty_pages_pct+"]")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_open_files", "purple", "", false, false)+"["+tmp.Var_innodb_open_files+"]"+utils.Colorize(" innodb_read_io_threads", "purple", "", false, false)+"["+tmp.Var_innodb_read_io_threads+"]"+utils.Colorize(" innodb_thread_concurrency", "purple", "", false, false)+"["+tmp.Var_innodb_thread_concurrency+"]"+"")
					fmt.Printf("%s  %s\n", "                  ", utils.Colorize("innodb_write_io_threads", "purple", "", false, false)+"["+tmp.Var_innodb_write_io_threads+"]"+"\n")
				}

				FilterTitle(in, num, interval)
				FilterValue(in, num, interval, tmp)
				Before = tmp
				tmp.CloseConn()
			}
			num++
			// 终止循环
			if !ok {
				break
			}
		}
	}
	return nil
}

// 组装标题
func FilterTitle(in string, count, interval int) {
	title := utils.GetTimeTitle()
	columns := utils.GetTimeColumns()

	if strings.Contains(in, "-lazy") {
		title += utils.GetComTitle()
		columns += utils.GetComColumns()
		title += utils.GetHitTitle()
		columns += utils.GetHitColumns()
	}
	if strings.Contains(in, "-com") {
		title += utils.GetComTitle()
		columns += utils.GetComColumns()
	}
	if strings.Contains(in, "-hit") {
		title += utils.GetHitTitle()
		columns += utils.GetHitColumns()
	}
	if strings.Contains(in, "-innodb") {
		title += utils.GetInnodbPagesTitle()
		columns += utils.GetInnodbPagesColumns()

		title += utils.GetInnodbDataTitle()
		columns += utils.GetInnodbDataColumns()

		title += utils.GetInnodbLogTitle()
		columns += utils.GetInnodbLogColumns()

		title += utils.GetInnodbStatusTitle()
		columns += utils.GetInnodbStatusColumns()
	}
	if strings.Contains(in, "-innodb_rows") {
		title += utils.GetInnodbRowsTitle()
		columns += utils.GetInnodbRowsColumns()
	}
	if strings.Contains(in, "-innodb_pages") {
		title += utils.GetInnodbPagesTitle()
		columns += utils.GetInnodbPagesColumns()
	}
	if strings.Contains(in, "-innodb_data") {
		title += utils.GetInnodbDataTitle()
		columns += utils.GetInnodbDataColumns()
	}
	if strings.Contains(in, "-innodb_log") {
		title += utils.GetInnodbLogTitle()
		columns += utils.GetInnodbLogColumns()
	}
	if strings.Contains(in, "-innodb_status") {
		title += utils.GetInnodbStatusTitle()
		columns += utils.GetInnodbStatusColumns()
	}
	if strings.Contains(in, "-T") {
		title += utils.GetThreadsTitle()
		columns += utils.GetThreadsColumns()
	}
	if strings.Contains(in, "-B") {
		title += utils.GetBytesTitle()
		columns += utils.GetBytesColumns()
	}
	if strings.Contains(in, "-semi") {
		title += utils.GetSemiTitle()
		columns += utils.GetSemiColumns()
	}
	if strings.Contains(in, "-slave") {
		title += utils.GetSlaveTitle()
		columns += utils.GetSlaveColumns()
	}

	if count%20 == 0 {
		fmt.Println(title)
		fmt.Println(columns)
	}
}

// if 顺序决定展示命令
func FilterValue(in string, num, interval int, mysql *basic) error {
	value, err := utils.TimeNow()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// -t,-l,-c,-s,-com,-hit
	if strings.Contains(in, "-lazy") {
		if num == 0 {
			value += utils.Colorize("    0     0     0      0     0", "green", "", false, false) + utils.Colorize("|", "green", "", false, false)
			value += utils.Colorize("100.00 100.00 100.00 100.00 100.00       0 100.00", "green", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateCom(interval)
			value += mysql.CreateHit(interval)
		}
	}
	if strings.Contains(in, "-com") {
		if num == 0 {
			value += utils.Colorize("    0     0     0      0     0", "green", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateCom(interval)
		}
	}

	if strings.Contains(in, "-hit") {
		if num == 0 {
			value += utils.Colorize("100.00 100.00 100.00 100.00 100.00       0 100.00", "green", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateHit(interval)
		}
	}

	if strings.Contains(in, "-innodb") {
		if num == 0 {
			value += utils.Colorize("      0      0      0     0", "yellow", "", false, false) + utils.Colorize("|", "green", "", false, false)
			value += utils.Colorize("     0      0      0       0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
			value += utils.Colorize("     0       0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
			value += utils.Colorize("    0      0      0     0     0     0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbPages(interval)
			value += mysql.CreateInnodbData(interval)
			value += mysql.CreateInnodbLog(interval)
			value += mysql.CreateInnodbStatus(interval)
		}
	}

	if strings.Contains(in, "-innodb_rows") {
		if num == 0 {
			value += utils.Colorize("    0     0     0      0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbRows(interval)
		}
	}

	if strings.Contains(in, "-innodb_pages") {
		if num == 0 {
			value += utils.Colorize("      0      0      0     0", "yellow", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbPages(interval)
		}
	}

	if strings.Contains(in, "-innodb_data") {
		if num == 0 {
			value += utils.Colorize("     0      0      0       0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbData(interval)
		}
	}

	if strings.Contains(in, "-innodb_log") {
		if num == 0 {
			value += utils.Colorize("     0       0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbLog(interval)
		}
	}

	if strings.Contains(in, "-innodb_status") {
		if num == 0 {
			value += utils.Colorize("    0      0      0     0     0     0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateInnodbStatus(interval)
		}
	}

	if strings.Contains(in, "-T") {
		if num == 0 {
			value += utils.Colorize("   0    0    0    0      0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateThreads(interval)
		}
	}

	if strings.Contains(in, "-B") {
		if num == 0 {
			value += utils.Colorize("      0      0", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateBytes(interval)
		}
	}

	if strings.Contains(in, "-semi") {
		if num == 0 {
			value += utils.Colorize("100ms 100ms 1000 1000  1000", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateSemi(interval)
		}
	}

	if strings.Contains(in, "-slave") {
		if num == 0 {
			value += utils.Colorize(" 1066312331  1066312331 6312331 6312331", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateSlave(interval)
		}
	}

	fmt.Println(value)
	return nil
}
