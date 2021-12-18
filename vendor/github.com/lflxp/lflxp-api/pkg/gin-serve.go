package pkg

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	bolt "go.etcd.io/bbolt"
)

func RegisterAPI(router *gin.Engine) {
	apiGroup := router.Group("/api/v1")

	apiGroup.GET("/db/backup", BackUP)                        // 备份
	apiGroup.GET("/bucket/get", GetBucket)                    // 获取所有bucket
	apiGroup.POST("/bucket/add/:name", CreateBucket)          // 创建bucket
	apiGroup.DELETE("/bucket/delete/:name", DeleteBucket)     // 删除bucket
	apiGroup.DELETE("/key/delete/:bucket/:key", DeleteKey)    // 删除指定bucket下的key
	apiGroup.POST("/key/add/:bucket/:key/:value", AddKey)     // 新增指定bucket下的key
	apiGroup.POST("/key/batch/:bucket/:key/:value", BatchKey) // 批量新增指定bucket下的key
	apiGroup.GET("/key/get/:bucket/:key", GetKey)             // 查询指定bucket下的key
	apiGroup.GET("/key/prefix/:bucket/:key", GetPreifxKey)    // 按前缀搜索指定bucket数据
	apiGroup.GET("/key/range/:bucket/:min/:max", GetRangeKey) // 按范围搜索指定bucket数据
	apiGroup.GET("/key/get/:bucket", GetAllKey)               // 查询指定bucket下所有key
}

// @Summary  BackUP数据库备份
// @Description curl http://localhost/api/v1/db/backup > my.db
// @Tags DB
// @Success 200 {string} string "success" //成功返回的数据结构， 最后是示例
// @Router /api/v1/db/backup [get]
func BackUP(c *gin.Context) {
	err := boltDB.Bolt.View(func(tx *bolt.Tx) error {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		c.Writer.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(c.Writer)
		return err
	})
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary  查询Bucket接口
// @Description Bucket GET
// @Tags Bucket
// @Success 200 {string} string "success" //成功返回的数据结构， 最后是示例
// @Router /api/v1/bucket/get [get]
func GetBucket(c *gin.Context) {
	result := []string{}
	boltDB.Bolt.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			b := []string{string(name)}
			result = append(result, b...)
			return nil
		})
	})
	c.JSONP(200, result)
}

// @Summary  新增Bucket接口
// @Description Bucket Add
// @Tags Bucket
// @Param name path string true "NAME"
// @Success 200 {string} string "success" //成功返回的数据结构， 最后是示例
// @Router /api/v1/bucket/add/{name} [post]
func CreateBucket(c *gin.Context) {
	name := c.Params.ByName("name")
	if name == "" {
		c.String(200, "no bucket name")
		return
	}
	boltDB.Bolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	c.String(200, "success")
}

// @Summary  删除Bucket接口
// @Description DELETE Bucket
// @Tags Bucket
// @Param name path string true "NAME"
// @Success 200 {string} string "success"
// @Router /api/v1/bucket/delete/{name} [delete]
func DeleteBucket(c *gin.Context) {
	name := c.Params.ByName("name")
	if name == "" {
		c.String(200, "no bucket name")
		return
	}
	boltDB.Bolt.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(name))
		if err != nil {
			c.String(200, "error no such bucket")
			return fmt.Errorf("bucket: %s", err)
		}
		return nil
	})
	c.String(200, "success")
}

// @Summary  删除Key接口
// @Description 删除指定bucket下的key
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param key path string true "KEY"
// @Success 200 {string} string "success"
// @Router /api/v1/key/delete/{bucket}/{key} [delete]
func DeleteKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	key := c.Params.ByName("key")
	if key == "" || bucket == "" {
		c.String(200, "bucket or key is none")
		return
	}
	boltDB.Bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			c.String(200, "error no such bucket")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Delete([]byte(key))
		if err != nil {
			c.String(200, "error Deleting KV")
			return fmt.Errorf("delete kv: %s", err)
		}
		return nil
	})
	c.String(200, "success")
}

// @Summary  批量添加Key接口
// @Description 批量添加指定bucket下的key
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param key path string true "KEY"
// @Param value path string true "VALUE"
// @Success 200 {string} string "success"
// @Router /api/v1/key/batch/{bucket}/{key}/{value} [post]
func BatchKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	key := c.Params.ByName("key")
	value := c.Params.ByName("value")
	if key == "" || bucket == "" {
		c.String(200, "bucket or key is none")
		return
	}
	boltDB.Bolt.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			c.String(200, "error no such bucket")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			c.String(200, "error writing KV")
			return fmt.Errorf("create kv: %s", err)
		}
		return nil
	})
	c.String(200, "success")
}

// @Summary  添加Key接口
// @Description 添加指定bucket下的key
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param key path string true "KEY"
// @Param value path string true "VALUE"
// @Success 200 {string} string "success"
// @Router /api/v1/key/add/{bucket}/{key}/{value} [post]
func AddKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	key := c.Params.ByName("key")
	value := c.Params.ByName("value")
	if key == "" || bucket == "" {
		c.String(200, "bucket or key is none")
		return
	}
	boltDB.Bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			c.String(200, "error no such bucket")
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(key), []byte(value))
		if err != nil {
			c.String(200, "error writing KV")
			return fmt.Errorf("create kv: %s", err)
		}
		return nil
	})
	c.String(200, "success")
}

// @Summary  查询Key接口
// @Description 查询指定bucket下的key
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param key path string true "KEY"
// @Success 200 {string} string "success"
// @Router /api/v1/key/get/{bucket}/{key} [get]
func GetKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	key := c.Params.ByName("key")
	if key == "" || bucket == "" {
		c.String(200, "bucket or key is none")
		return
	}
	var (
		result []byte
		status string
	)

	boltDB.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			result = b.Get([]byte(key))
			status = "success"
		} else {
			status = "error opening bucket| does it exist? | n"
		}
		return nil
	})
	c.JSONP(200, gin.H{
		key:      string(result),
		"status": status,
	})
}

type Result struct {
	Name   string
	Data   map[string]string
	Count  int
	Status string
}

// @Summary  查询Bucket All Key接口
// @Description 查询指定bucket下所有key
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Success 200 {string} string "success"
// @Router /api/v1/key/get/{bucket} [get]
func GetAllKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	if bucket == "" {
		c.String(200, "bucket  is none")
		return
	}
	result := Result{
		Name: bucket,
		Data: map[string]string{},
	}

	count := 0

	boltDB.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				result.Data[string(k)] = string(v)
				if count > 10000 {
					break
				}
				count++
			}
			result.Count = count
			result.Status = "success"
		} else {
			result.Status = "no such bucket available"
		}
		return nil
	})
	c.JSONP(200, result)
}

// @Summary  Search Prefix Key接口
// @Description 按前缀查询指定bucket数据
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param key path string true "KEY"
// @Success 200 {string} string "success"
// @Router /api/v1/key/prefix/{bucket}/{key} [get]
func GetPreifxKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	key := c.Params.ByName("key")
	if bucket == "" || key == "" {
		c.String(200, "bucket or key is none")
		return
	}
	result := Result{
		Name: bucket,
		Data: map[string]string{},
	}

	count := 0

	boltDB.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			prefix := []byte(key)
			c := b.Cursor()
			for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
				result.Data[string(k)] = string(v)
				if count > 10000 {
					break
				}
				count++
			}
			result.Count = count
			result.Status = "success"
		} else {
			result.Status = "no such bucket available"
		}
		return nil
	})
	c.JSONP(200, result)
}

// @Summary  Search Range接口
// @Description 范围搜索指定bucket数据
// @Tags Key
// @Param bucket path string true "BUCKET"
// @Param min path string true "MIN"
// @Param max path string true "MAX"
// @Success 200 {string} string "success"
// @Router /api/v1/key/range/{bucket}/{min}/{max} [get]
func GetRangeKey(c *gin.Context) {
	bucket := c.Params.ByName("bucket")
	min := c.Params.ByName("min")
	max := c.Params.ByName("max")
	if bucket == "" || max == "" || min == "" {
		c.String(200, "bucket or min or max is none")
	}
	result := Result{
		Name: bucket,
		Data: map[string]string{},
	}

	count := 0

	boltDB.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b != nil {
			min := []byte(min)
			max := []byte(max)
			c := b.Cursor()
			for k, v := c.Seek(min); bytes.Compare(k, max) <= 0; k, v = c.Next() {
				result.Data[string(k)] = string(v)
				// if count > 10000 {
				// 	break
				// }
				count++
			}
			result.Count = count
			result.Status = "success"
		} else {
			result.Status = "no such bucket available"
		}
		return nil
	})
	c.JSONP(200, result)
}
