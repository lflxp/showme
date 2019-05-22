package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"

	bolt "github.com/coreos/bbolt"
	log "github.com/sirupsen/logrus"
)

var (
	Db   *bolt.DB
	Mmap map[string]string
	err  error
)

const demo string = `[{
	"id": 0,
	"index": [0],
	"label": "1",
	"children": [{
		"id": 1,
		"index": [0, 0],
		"label": "1-1",
		"children": [{
			"id": 2,
			"index": [0, 0, 0],
			"label": "1-1-1",
			"children": null
		}]
	}, {
		"id": 3,
		"index": [0, 1],
		"label": "1-2",
		"children": [{
			"id": 4,
			"index": [0, 1, 0],
			"label": "1-2-1",
			"children": null
		}]
	}]
}, {
	"id": 9,
	"index": [1],
	"label": "2",
	"children": [{
		"id": 5,
		"index": [1, 0],
		"label": "2-1",
		"children": [{
			"id": 6,
			"index": [1, 0, 0],
			"label": "2-1-1",
			"children": null
		}]
	}, {
		"id": 7,
		"index": [1, 1],
		"label": "2-2",
		"children": [{
			"id": 8,
			"index": [1, 1, 0],
			"label": "2-2-1",
			"children": null
		}]
	}]
}]`

func InitDB(dbName, defaultDb string) {
	Mmap = map[string]string{}
	Db, err = bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	//init tables

	err = CreateBucket(defaultDb)
	if err != nil {
		fmt.Println(err.Error())
		// panic(err)
		log.WithFields(log.Fields{
			"bbolt.go": "CreateBucket",
		}).Error(err.Error())
	}

	//初始化user表
	us, _ := GetValueByBucketName(defaultDb, "test")
	if len(us) == 0 {
		beego.Critical("初始化测试数据")
		AddKeyValueByBucketName(defaultDb, "test", "test", demo)
	}
	GetAllTables(defaultDb)
	// log.Println(Mmap)
	log.Println("init db logging success")
}

func CreateBucket(db string) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(db))
		if err != nil {
			return fmt.Errorf("create bucket: %s ", err.Error())
		}
		//如果没有index表和index key则立
		return nil
	})
	AddTables(db, "test")
	return err
}

func DeleteBucket(db, tablename string) error {
	err := Db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(tablename))
		if err != nil {
			return fmt.Errorf("delete bucket: %s ", err.Error())
		}
		return nil
	})
	DeleteTables(db, tablename)
	return err
}

func AddTables(db, tablename string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Put([]byte(tablename), []byte(tablename))
		Mmap[string(tablename)] = string(tablename)
		return err
	})
}

func GetAllTables(db string) (map[string]string, error) {
	tmp := map[string]string{}
	err := Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp[string(k)] = string(v)
			Mmap[string(k)] = string(v)
		}
		return nil
	})
	return tmp, err
}

func GetAllByTables(name string) (map[string]string, error) {
	tmp := map[string]string{}
	err := Db.View(func(tx *bolt.Tx) error {
		if _, ok := Mmap[name]; !ok {
			return errors.New(fmt.Sprintf("%s not exist", name))
		}
		b := tx.Bucket([]byte(name))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tmp[string(k)] = string(v)
			Mmap[string(k)] = string(v)
		}
		return nil
	})
	return tmp, err
}

func DeleteTables(db, tablename string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(db))
		err := b.Delete([]byte(tablename))
		delete(Mmap, tablename)
		return err
	})
}

func AddKeyValueByBucketName(db, table, key, value string) error {
	// fmt.Println(Mmap)
	if _, ok := Mmap[table]; !ok {
		log.Debug(fmt.Printf("%s is not exist\n", table))
		CreateBucket(db)
	}
	// log.Println(Mmap)
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
}

func AddKeyValueByBucketNameAuto(table, key, value string) error {

	return Db.Update(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		var err error
		if _, ok := Mmap[table]; !ok {
			b, err = tx.CreateBucket([]byte(table))
			if err != nil {
				return err
			}
		} else {
			b = tx.Bucket([]byte(table))
			if err != nil {
				return err
			}
		}

		err = b.Put([]byte(key), []byte(value))
		return err
	})
}

func GetValueByBucketName(table, key string) ([]byte, error) {
	var value []byte
	var err error
	Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if _, ok := Mmap[table]; ok {
			value = b.Get([]byte(key))
		} else {
			err = errors.New(fmt.Sprintf("table %s not exist", table))
		}
		return err
	})
	return value, err
}

func DeleteKeyValueByBucket(table, key string) error {
	return Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		err := b.Delete([]byte(key))
		return err
	})
}
