package utils

import (
	"fmt"
	"log/slog"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var Engine *xorm.Engine

func InitSqlite() {
	homepath, err := Home()
	if err != nil {
		panic(err)
	}

	Engine, err = xorm.NewEngine("sqlite3", fmt.Sprintf("%s/.showme.db", homepath))
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(fmt.Sprintf("初始化sqlite数据库 %s/.showme.db", homepath))
}
