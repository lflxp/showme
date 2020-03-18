package tty

import (
	"time"

	"github.com/lflxp/showme/utils"
)

// TODO
// cgo效率太低 后期转mmap或者bbolt
type Aduit struct {
	Id         int64     `json:"id"`
	Remoteaddr string    `json:"remoteaddr" xorm:"varchar(40) notnull index"`
	Token      string    `json:"token" xorm:"varchar(32)"`
	Command    string    `json:"command" xorm:"varchar(10)"`
	Pid        int       `json:"pid"`
	Status     string    `json:"status" xorm:"varchar(40)"`
	Created    time.Time `json:"created" xorm:"created"`
}

func AddAduit(data *Aduit) error {
	_, err := utils.Engine.Insert(data)
	return err
}

func GetAduit(name string) ([]Aduit, error) {
	var err error
	data := make([]Aduit, 0)
	if name == "" {
		err = utils.Engine.Desc("created").Limit(300).Find(&data)
	} else {
		err = utils.Engine.Where("id = ? or remoteaddr = ? or  token = ?", name, name, name).Desc("created").Limit(300).Find(&data)
	}
	return data, err
}

// 记录谁来访问过
type Whos struct {
	Id         int64     `json:"id"`
	Remoteaddr string    `json:"remoteaddr" xorm:"varchar(40)" notnull index`
	Path       string    `json:"path" xorm:"varchar(20)"`
	Created    time.Time `json:"created" xorm:"created"`
}

func AddWhos(data *Whos) error {
	_, err := utils.Engine.Insert(data)
	return err
}

func GetWhos(name string) ([]Whos, error) {
	var err error
	data := make([]Whos, 0)
	if name == "" {
		err = utils.Engine.Desc("created").Limit(300).Find(&data)
	} else {
		err = utils.Engine.Where("id = ? or remoteaddr = ? or  path = ?", name, name, name).Desc("created").Limit(300).Find(&data)
	}
	return data, err
}
