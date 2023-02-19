package music

import (
	"github.com/lflxp/lflxp-music/core/middlewares/template"
	"github.com/lflxp/tools/orm/sqlite"
)

func init() {
	template.Register(new(Musichistory))
}

type Musichistory struct {
	Id       int64   `xorm:"id" name:"id" json:"id"`
	Duration float64 `xorm:"duration" name:"duration" verbose_name:"duration字段测试" search:"true" json:"duration"`
	Album    string  `xorm:"album" name:"album" json:"album"`
	Image    string  `xorm:"image" name:"image" json:"image"`
	Name     string  `xorm:"name notnull unique" name:"name" json:"name"`
	Url      string  `xorm:"url" name:"url" json:"url"`
	Singer   string  `xorm:"singer" name:"singer" json:"singer"`
	User     string  `xorm:"user" name:"user" json:"user"`
}

func (m *Musichistory) Post() (int64, error) {
	return sqlite.NewOrm().Insert(m)
}

type PostData struct {
	Param struct {
		Data Musichistory `json:"data"`
	} `json:"param"`
}
