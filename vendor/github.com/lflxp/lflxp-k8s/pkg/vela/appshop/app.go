package appshop

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/guonaihong/gout"
	"github.com/lflxp/tools/httpclient"
	"github.com/spf13/viper"
)

const (
	REPO_LIST  = "/shop/list"
	REPO_PARAM = "/shop/param"
)

func RegisterShop(router *gin.Engine) {
	shopGroup := router.Group("/api")
	{
		shopGroup.GET(REPO_LIST, repo_list)
		shopGroup.GET(REPO_PARAM, repo_param)
		shopGroup.POST(REPO_PARAM, repo_param_post)
		shopGroup.DELETE(REPO_PARAM, repo_param_delete)
		shopGroup.GET("/shop/test", repo_test)
	}
}

func repo_list(c *gin.Context) {
	token := c.Request.Header.Get("token")
	host := viper.GetString("auth.url")
	url := fmt.Sprintf("%s/api/v1/addons", host)
	code := 0
	body := ""
	err := httpclient.NewGoutClient().
		GET(url).
		// Debug(true).
		SetHeader(gout.H{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}).
		BindBody(&body).
		Code(&code).
		Do()

	if err != nil {
		httpclient.SendErrorMessage(c, code, "http request failed", err.Error())
		return
	}

	if code != 200 {
		slog.Error(body)
		httpclient.SendErrorMessage(c, code, fmt.Sprintf("%d", code), body)
		return
	}

	var data interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		slog.Error(body)
		httpclient.SendErrorMessage(c, 200, "json parse error", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, data)
}

func repo_test(c *gin.Context) {
	c.Redirect(302, "http://localhost:8080/d2admin/")
}

func repo_param(c *gin.Context) {
	token := c.Request.Header.Get("token")
	host := viper.GetString("auth.url")
	param, ok := c.GetQuery("path")
	if !ok {
		httpclient.SendErrorMessage(c, 400, "path not found", "path not found")
		return
	}
	url := fmt.Sprintf("%s/api/v1/%s", host, param)
	code := 0
	body := ""
	err := httpclient.NewGoutClient().
		GET(url).
		// Debug(true).
		SetHeader(gout.H{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}).
		BindBody(&body).
		Code(&code).
		Do()

	if err != nil {
		httpclient.SendErrorMessage(c, code, "http request failed", err.Error())
		return
	}

	if code != 200 {
		slog.Error(body)
		httpclient.SendErrorMessage(c, code, fmt.Sprintf("%d", code), body)
		return
	}

	var data interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		slog.Error(body)
		httpclient.SendErrorMessage(c, 200, "json parse error", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, data)
}

func repo_param_post(c *gin.Context) {
	var postData interface{}
	err := c.BindJSON(&postData)
	if err != nil {
		slog.Error(err.Error())
		httpclient.SendErrorMessage(c, 500, "bind json error", err.Error())
		return
	}
	token := c.Request.Header.Get("token")
	host := viper.GetString("auth.url")
	param, ok := c.GetQuery("path")
	if !ok {
		httpclient.SendErrorMessage(c, 400, "path not found", "path not found")
		return
	}
	url := fmt.Sprintf("%s/api/v1/%s", host, param)
	code := 0
	body := ""
	err = httpclient.NewGoutClient().
		POST(url).
		SetJSON(postData).
		Debug(true).
		SetHeader(gout.H{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}).
		BindBody(&body).
		Code(&code).
		Do()

	if err != nil {
		httpclient.SendErrorMessage(c, code, "http request failed", err.Error())
		return
	}

	if code != 200 {
		slog.Error(body)
		httpclient.SendErrorMessage(c, code, fmt.Sprintf("%d", code), body)
		return
	}

	var data interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		slog.Error(body)
		httpclient.SendErrorMessage(c, 200, "json parse error", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, data)
}

func repo_param_delete(c *gin.Context) {
	// 	var postData interface{}
	// 	err := c.BindJSON(&postData)
	// 	if err != nil {
	// 		slog.Error(err.Error())
	// 		httpclient.SendErrorMessage(c, 500, "bind json error", err.Error())
	// 		return
	// 	}
	token := c.Request.Header.Get("token")
	host := viper.GetString("auth.url")
	param, ok := c.GetQuery("path")
	if !ok {
		httpclient.SendErrorMessage(c, 400, "path not found", "path not found")
		return
	}
	url := fmt.Sprintf("%s/api/v1/%s", host, param)
	code := 0
	body := ""
	err := httpclient.NewGoutClient().
		DELETE(url).
		Debug(true).
		SetHeader(gout.H{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
		}).
		BindBody(&body).
		Code(&code).
		Do()

	if err != nil {
		httpclient.SendErrorMessage(c, code, "http request failed", err.Error())
		return
	}

	if code != 200 {
		slog.Error(body)
		httpclient.SendErrorMessage(c, code, fmt.Sprintf("%d", code), body)
		return
	}

	var data interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		slog.Error(body)
		httpclient.SendErrorMessage(c, 200, "json parse error", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, data)
}
