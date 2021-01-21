package pkg

import (
	"fmt"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
)

var boltDB *storm.DB

// New/Get boltDBs
func NewBolt() *storm.DB {
	if boltDB == nil {
		log.Info("初始化bolt数据库")
		homepath, err := Home()
		if err != nil {
			panic(err)
		}

		boltDB, err = storm.Open(fmt.Sprintf("%s/.showme.bolt", homepath))
		if err != nil {
			panic(err)
		}
	}
	return boltDB
}
