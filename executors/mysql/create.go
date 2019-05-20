package mysql

import (
	"fmt"
	"strings"

	"github.com/lflxp/showme/utils"
	"github.com/shirou/gopsutil/host"
)

func (this *basic) GetHostAndIps() error {
	n, err := host.Info()
	if err != nil {
		return err
	}
	this.Hostname = n.Hostname
	this.Ips = strings.Join(utils.GetIps(), ",")
	return nil
}

func (this *basic) GetMysqlConn() error {
	conn, err := utils.MysqlConn(username, password, ip, port, dbname)
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
		fmt.Println(variable_name, value)
		showGlobalVariables[variable_name] = value
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
		fmt.Println(variable_name, value)
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
		fmt.Println(variable_name, value)
		showGlobalStatus[variable_name] = value
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
		fmt.Println(variable_name, value)
		showStatus[variable_name] = value
	}
	return nil
}

// todo show processlist
// todo show slave status
// todo mmp value to struct
