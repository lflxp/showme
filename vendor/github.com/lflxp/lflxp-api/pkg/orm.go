package pkg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TimeSeries struct {
	ID      int               `storm:"id,increment"`
	Name    string            `json:"name" storm:"index"`
	Data    map[string]string `json:"data"`
	Created string            `storm:"index"`
}

func (this *TimeSeries) Save() error {
	err := boltDB.Save(this)
	return err
}

func RegisterOrmAPI(router *gin.Engine) {
	apiGroup := router.Group("/api/v1/orm/")

	apiGroup.GET("/getall", GetOrmAllData)
	apiGroup.GET("/range/:min/:max", GetOrmRangeData)
	apiGroup.POST("/add", AddOrm)
}

// @Summary  查询TimeSeries接口
// @Description TimeSeries Range范围查询
// @Tags TimeSeries
// @Param min path string true "MIN"
// @Param max path string true "MAX"
// @Success 200 {string} string "success" //成功返回的数据结构， 最后是示例
// @Router /api/v1/orm/range/{min}/{max} [get]
func GetOrmRangeData(c *gin.Context) {
	min := c.Params.ByName("min")
	max := c.Params.ByName("max")
	if max == "" || min == "" {
		c.String(200, "min or max is none")
	}

	var datas []TimeSeries
	err := boltDB.Range("Created", min, max, &datas)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSONP(200, datas)
}

// @Summary  查询TimeSeries接口
// @Description TimeSeries GET
// @Tags TimeSeries
// @Success 200 {string} string "success" //成功返回的数据结构， 最后是示例
// @Router /api/v1/orm/getall [get]
func GetOrmAllData(c *gin.Context) {
	var datas []TimeSeries
	err := boltDB.All(&datas)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSONP(200, datas)
}

// @Summary  新增TimeSeries接口
// @Description TimeSeries Add
// @Tags TimeSeries
// @Param data body TimeSeries true "NAME"
// @Success 200 {object} TimeSeries TimeSeries{} //成功返回的数据结构， 最后是示例
// @Router /api/v1/orm/add [post]
func AddOrm(c *gin.Context) {
	var temp TimeSeries
	data, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(data, &temp)
	if err != nil {
		c.JSONP(http.StatusBadRequest, err.Error())
	} else {
		// temp.Created = time.Now().UnixNano()
		temp.Created = time.Now().Format("20060102150405")
		err := temp.Save()
		if err != nil {
			c.JSONP(http.StatusOK, err.Error())
		} else {
			c.String(http.StatusOK, "success add")
		}
	}
}
