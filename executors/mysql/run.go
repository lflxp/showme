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
		Username = "root"
		Password = "Nw8jfr3zDZ"
		Ip = "10.128.142.132"
		Port = "3306"
		Dbname = "user"
		err := mysql.InitMysqlConn()
		if err != nil {
			return err
		}
		defer mysql.CloseConn()
	} else {
		return errors.New("nothing input")
	}
	if in == "mysql test GetHostAndIps" {
		err := mysql.GetHostAndIps()
		if err != nil {
			return err
		}
		fmt.Println(fmt.Sprintf("Hostname %s\nIps %s", mysql.Hostname, mysql.Ips))
	} else if in == "mysql test GetShowDatabases" {
		err := mysql.GetShowDatabases()
		if err != nil {
			return err
		}
		fmt.Printf("Dbs %s\n", mysql.Dbs)
	} else if in == "mysql test GetShowGlobalVariables" {
		err := mysql.GetShowGlobalVariables()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowVariables" {
		err := mysql.GetShowVariables()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowGlobalStatus" {
		err := mysql.GetShowGlobalStatus()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowStatus" {
		err := mysql.GetShowStatus()
		if err != nil {
			return err
		}
	} else if in == "mysql test GetShowEngineInnodbStatus" {
		err := mysql.GetShowEngineInnodbStatus()
		if err != nil {
			return err
		}
	} else if in == "mysql processlist" {
		err := mysql.GetShowProcesslist()
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

		interval := 20
		num := 0

		// 主机信息
		mysql.GetHostAndIps()
		mysql.GetShowDatabases()
		mysql.GetShowGlobalStatus()
		mysql.GetShowGlobalVariables()
		mysql.GetShowStatus()
		mysql.GetShowVariables()

		// print net info
		// xo := utils.MonitorNet{}
		// xo.Get()

		// err := utils.GetHostInfo()
		// if err != nil {
		// 	fmt.Println(err.Error())
		// 	return
		// }

		FilterTitle(in, num, interval)

		for {
			select {
			case s := <-c:
				fmt.Printf("\n\033[1;4;31m%s:罒灬罒:小伙子走了哟！\033[0m\n", s)
				ok = false
				break
			case <-t.C:
				FilterTitle(in, num, interval)
				FilterValue(in, num, interval, mysql)
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
	} else {
		if strings.Contains(in, "-com") {
			title += utils.GetComTitle()
			columns += utils.GetComColumns()
		}
		if strings.Contains(in, "-hit") {
			title += utils.GetHitTitle()
			columns += utils.GetHitColumns()
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
		if strings.Contains(in, "-threads") {
			title += utils.GetThreadsTitle()
			columns += utils.GetThreadsColumns()
		}
		if strings.Contains(in, "-bytes") {
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
	}

	if count%interval == 0 {
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
			value += utils.Colorize("0 0 0  0 0", "", "", false, false) + utils.Colorize("|", "dgreen", "", false, false)
			value += utils.Colorize("100.00 100.00 100.00 100.00 100.00   0 100.00|", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
		} else {
			value += mysql.CreateCom(interval, Before)
			value += mysql.CreateHit(interval)
		}
	} else {
		if strings.Contains(in, "-com") {
			if num == 0 {
				value += utils.Colorize("0 0 0  0 0", "", "", false, false) + utils.Colorize("|", "dgreen", "", false, false)
			} else {
				value += mysql.CreateCom(interval, Before)
			}
		}

		if strings.Contains(in, "-hit") {
			if num == 0 {
				value += utils.Colorize("100.00 100.00 100.00 100.00 100.00   0 100.00|", "", "", false, false) + utils.Colorize("|", "green", "", false, false)
			} else {
				value += mysql.CreateHit(interval)
			}
		}
	}

	fmt.Println(value)
	Before = mysql
	return nil
}
