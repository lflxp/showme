package tty

import (
	"time"

	"github.com/lflxp/showme/utils"
)

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
		err = utils.Engine.Desc("created").Find(&data)
	} else {
		err = utils.Engine.Where("id = ? or remoteaddr = ? or  token = ?", name, name, name).Desc("created").Find(&data)
	}
	return data, err
}
