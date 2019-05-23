package utils

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "./showme.db")
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("初始化sqlite数据库showme.db")
}
