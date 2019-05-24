package api

import (
	"database/sql"
	"errors"

	"github.com/lflxp/showme/utils"
)

type Shell struct {
	Cmd   string `json:"cmd"`
	Token string `json:"token"`
}

func (this *Shell) Exec() (string, error) {
	if this.VaildToken() {
		rs, err := utils.ExecCommand(this.Cmd)
		return string(rs), err
	}
	return "", errors.New("token is failed")
}

func (this *Shell) VaildToken() bool {
	var rs bool
	if this.Token == "lxp" {
		rs = true
	} else {
		rs = false
	}
	return rs
}

type Db struct {
	Sql   string   `json:"sql"`
	Args  []string `json:"args"`
	Token string   `json:"token"`
}

func (this *Db) Exec() (sql.Result, error) {
	var (
		rs  sql.Result
		err error
	)
	if this.VaildToken() {
		rs, err = utils.Engine.Exec(this.Sql)
	}

	return rs, err
}

func (this *Db) Select() ([]map[string]string, error) {
	var (
		rs  []map[string]string
		err error
	)
	if this.VaildToken() {
		rs, err = utils.Engine.QueryString(this.Sql)
	}
	return rs, err
}

func (this *Db) VaildToken() bool {
	var rs bool
	if this.Token == "lxp" {
		rs = true
	} else {
		rs = false
	}
	return rs
}
