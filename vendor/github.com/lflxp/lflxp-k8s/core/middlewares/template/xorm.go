package template

import (
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

var (
	ormOnce sync.Once
	orm     *xorm.Engine
)

func NewOrm() *xorm.Engine {
	ormOnce.Do(func() {
		var err error
		orm, err = xorm.NewEngine("sqlite3", "./demo.db")
		if err != nil {
			panic(err)
		}

		orm.ShowSQL(true)
		orm.Logger().SetLevel(log.LogLevel(core.LOG_DEBUG))
		orm.SetMaxIdleConns(300)
		orm.SetMaxOpenConns(300)
		orm.SetMapper(core.SnakeMapper{})
		// logger.Errorf("viper------------- is %s ",viper.GetString("snakemapper"))
		// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, viper.GetString("snakemapper"))
		// Engine.SetTableMapper(tbMapper)
		orm.SetColumnMapper(core.SameMapper{})

		// err = orm.Sync2(new(History))
		// if err != nil {
		// 	logger.Fatal(err.Error())
		// }
	})

	return orm
}
