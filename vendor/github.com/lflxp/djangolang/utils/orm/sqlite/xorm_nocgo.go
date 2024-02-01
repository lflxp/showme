package sqlite

import (
	"log/slog"
	"path/filepath"
	"sync"

	"k8s.io/client-go/util/homedir"
	_ "modernc.org/sqlite"
	"xorm.io/core"
	"xorm.io/xorm"
)

var (
	ormOnce sync.Once
	orm     *xorm.Engine
	dbname  string
	// 数据库驱动类型 默认sqlite
	driverName string
	// 数据库连接地址 root:123@/test?charset=utf8
	dataResourceName string
)

func NewOrm() *xorm.Engine {
	if dbname == "" || driverName == "" {
		slog.Debug("Dbname or DriverName is empty")
		home := homedir.HomeDir()
		dbname = ".djangolang.db"
		dataResourceName = filepath.Join(home, dbname)
		driverName = "sqlite"
	}

	ormOnce.Do(func() {

		var err error
		orm, err = xorm.NewEngine(driverName, dataResourceName)
		if err != nil {
			panic(err)
		}

		orm.ShowSQL(false)
		// orm.Logger().SetLevel(log.LogLevel(core.LOG_WARNING))
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
