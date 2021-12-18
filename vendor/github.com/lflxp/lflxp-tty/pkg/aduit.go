package pkg

import (
	"time"
)

// TODO
// cgo效率太低 后期转mmap或者bbolt
type Aduit struct {
	Id         int64     `json:"id" storm:"id,increment"`
	Remoteaddr string    `json:"remoteaddr"`
	Token      string    `json:"token" `
	Command    string    `json:"command"`
	Pid        int       `json:"pid"`
	Status     string    `json:"status"`
	Created    time.Time `json:"created"`
}

func (this *Aduit) Save() error {
	err := boltDB.Save(this)
	return err
}

func GetAduit(name string) ([]Aduit, error) {
	var rs []Aduit
	err := boltDB.All(&rs)
	return rs, err
}

// 记录谁来访问过
type Whos struct {
	Id         int64     `json:"id" storm:"id,increment"`
	Remoteaddr string    `json:"remoteaddr"`
	Path       string    `json:"path"`
	Created    time.Time `json:"created"`
}

func (this *Whos) Save() error {
	err := boltDB.Save(this)
	return err
}

func GetWhos(name string) ([]Whos, error) {
	var rs []Whos
	err := boltDB.All(&rs)
	return rs, err
}
