package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lflxp/showme/api/collector"
	"github.com/lflxp/showme/utils"
)

func AddSetup(router *gin.Engine) {
	router.GET("/setup", func(c *gin.Context) {
		c.JSON(http.StatusOK, collector.HardwareStatic)
	})
}

func AddCpu(router *gin.Engine) {
	router.GET("/cpu", func(c *gin.Context) {
		data, err := utils.CpuLoad()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			c.String(http.StatusOK, data)
		}
	})
}

func AddShell(router *gin.Engine) {
	router.POST("/shell", func(c *gin.Context) {
		var shell Shell
		data, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(data, &shell)
		// data, _ := ioutil.ReadAll(c.Request.Body)
		// log.Infof("cmd %s cmd2 %s query %s %s", cmd, cmd2, query, string(data))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			rs, err := shell.Exec()
			if err != nil {
				c.String(http.StatusOK, err.Error())
			}
			c.String(http.StatusOK, rs)
		}
	})
}

func AddSqlite(router *gin.Engine) {
	router.POST("/exec", func(c *gin.Context) {
		var db Db
		data, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(data, &db)
		// data, _ := ioutil.ReadAll(c.Request.Body)
		// log.Infof("cmd %s cmd2 %s query %s %s", cmd, cmd2, query, string(data))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			rs, err := db.Exec()
			if err != nil {
				c.String(http.StatusOK, err.Error())
			}
			c.JSON(http.StatusOK, rs)
		}
	})

	router.POST("/select", func(c *gin.Context) {
		var db Db
		data, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(data, &db)
		// data, _ := ioutil.ReadAll(c.Request.Body)
		// log.Infof("cmd %s cmd2 %s query %s %s", cmd, cmd2, query, string(data))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		} else {
			rs, err := db.Select()
			if err != nil {
				c.String(http.StatusOK, err.Error())
			}
			c.JSON(http.StatusOK, rs)
		}
	})
}
