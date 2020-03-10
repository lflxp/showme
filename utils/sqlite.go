package utils

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var Engine *xorm.Engine

func InitSqlite() {
	homepath, err := Home()
	if err != nil {
		panic(err)
	}

	Engine, err = xorm.NewEngine("sqlite3", fmt.Sprintf("%s/.showme.db", homepath))
	if err != nil {
		log.Error(err.Error())
	}
	log.Infof("初始化sqlite数据库 %s/.showme.db", homepath)
}
