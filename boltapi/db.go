package boltapi

import (
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/lflxp/showme/utils"
	log "github.com/sirupsen/logrus"
)

var boltDB *storm.DB

func init() {
	log.Info("初始化bolt数据库")
	homepath, err := utils.Home()
	if err != nil {
		panic(err)
	}

	boltDB, err = storm.Open(fmt.Sprintf("%s/.showme.bolt", homepath))
	if err != nil {
		panic(err)
	}
}
