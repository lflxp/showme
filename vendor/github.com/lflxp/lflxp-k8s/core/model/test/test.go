package test

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lflxp/lflxp-k8s/core/middlewares/template"

	"github.com/lflxp/tools/orm/sqlite"
)

func init() {
	template.Register(new(Demotest))
	// log.Debug("注册Demo test")
}

type Demotest struct {
	Id         int64  `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Country    string `json:"country" xorm:"varchar(255) not null"`
	Zoom       string `json:"zoom" xorm:"varchar(255) not null"`
	Company    string `json:"company" xorm:"varchar(255) not null"`
	Items      string `json:"items" xorm:"varchar(255) not null"`
	Production string `json:"production" xorm:"varchar(255) not null"`
	Count      string `json:"count" xorm:"varchar(255) not null"`
	Serial     string `json:"serial" xorm:"varchar(255) not null"`
	Extend     string `json:"extend" xorm:"varchar(255) not null"`
}

func (d *Demotest) GetByString(key string) (string, error) {

	// sqlite3 获取
	data, has, err := getByUUID(key)
	if err != nil {
		return "", err
	}

	// 不存在
	if !has {
		return "", errors.New(fmt.Sprintf("uuid %s not exists", key))
	}

	d = data

	// 放入缓存
	value, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(value, d)
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func getByUUID(uuid string) (*Demotest, bool, error) {
	data := new(Demotest)
	has, err := sqlite.NewOrm().Where("uuid = ?", uuid).Get(data)
	return data, has, err
}

func Add(data *Demotest) (int64, error) {
	affected, err := sqlite.NewOrm().Insert(data)
	return affected, err
}

func Del(id string) (int64, error) {
	data := new(Demotest)
	affected, err := sqlite.NewOrm().ID(id).Delete(data)
	return affected, err
}

func Update(id string, data *Demotest) (int64, error) {
	affected, err := sqlite.NewOrm().Table(new(Demotest)).ID(id).Update(data)
	return affected, err
}
