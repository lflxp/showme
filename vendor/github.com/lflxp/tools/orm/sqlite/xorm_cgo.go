//go:build cgo

package sqlite

// import (
// 	"path/filepath"
// 	"sync"

// 	_ "github.com/mattn/go-sqlite3"
// 	"k8s.io/client-go/util/homedir"
// 	"xorm.io/core"
// 	"xorm.io/xorm"
// )

// var (
// 	ormOnce sync.Once
// 	orm     *xorm.Engine
// )

// func NewOrm() *xorm.Engine {
// 	ormOnce.Do(func() {
// 		home := homedir.HomeDir()
// 		var err error
// 		orm, err = xorm.NewEngine("sqlite3", filepath.Join(home, ".lflxp-k8s.db"))
// 		if err != nil {
// 			panic(err)
// 		}

// 		orm.ShowSQL(false)
// 		// orm.Logger().SetLevel(log.LogLevel(core.LOG_WARNING))
// 		orm.SetMaxIdleConns(300)
// 		orm.SetMaxOpenConns(300)
// 		orm.SetMapper(core.SnakeMapper{})
// 		// logger.Errorf("viper------------- is %s ",viper.GetString("snakemapper"))
// 		// tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, viper.GetString("snakemapper"))
// 		// Engine.SetTableMapper(tbMapper)
// 		orm.SetColumnMapper(core.SameMapper{})

// 		// err = orm.Sync2(new(History))
// 		// if err != nil {
// 		// 	logger.Fatal(err.Error())
// 		// }
// 	})

// 	return orm
// }
