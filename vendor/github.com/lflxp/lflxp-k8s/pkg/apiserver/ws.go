package apiserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/lflxp-k8s/pkg/apiserver/model"
	"github.com/lflxp/lflxp-k8s/pkg/apiserver/ws/logs"
	"github.com/lflxp/lflxp-k8s/pkg/apiserver/ws/ssh"
	"github.com/lflxp/tools/httpclient"
)

const (
	WS_GET_LOGS      string = "/ws/logs/backend/:namespace/:pod/:container"
	WS_GET_LOGS_HTML string = "/ws/logs/html/:namespace/:pod/:container"
	WS_GET_SSH_HTML  string = "/ws/ssh/html/:namespace/:pod/:container"
	WS_GET_SSH       string = "/ws/ssh/backend/:namespace/:pod/:container"
)

func RegisterApiserverWS(router *gin.Engine) {
	wsGroup := router.Group("/")
	{
		wsGroup.GET(WS_GET_LOGS, ws_get_logs)
		wsGroup.GET(WS_GET_LOGS_HTML, ws_get_logs_html)
		wsGroup.GET(WS_GET_SSH_HTML, ws_get_ssh_html)
		wsGroup.GET(WS_GET_SSH, ws_get_ssh)
	}
}

func ws_get_logs_html(c *gin.Context) {
	// var data model.CoreV1
	// if err := c.BindJSON(&data); err != nil {
	// 	httpclient.SendErrorMessage(c, http.StatusBadRequest, "not found", err.Error())
	// 	return
	// }

	data := &model.CoreV1{}
	data.Namespace = c.Param("namespace")
	data.ContainerName = c.Param("container")
	data.PodName = c.Param("pod")

	c.HTML(200, "ws/logs.html", gin.H{
		"Host":      c.Request.Host,
		"Namespace": data.Namespace,
		"Container": data.ContainerName,
		"Pod":       data.PodName,
	})
}

func ws_get_logs(c *gin.Context) {
	data := &model.CoreV1{
		TailLines: 1000,
		Previous:  false,
		Fellow:    true,
	}
	data.Namespace = c.Param("namespace")
	data.ContainerName = c.Param("container")
	data.PodName = c.Param("pod")

	log.Debugf("data is %v", data)

	if data.Namespace == "" || data.ContainerName == "" || data.PodName == "" {
		httpclient.SendErrorMessage(c, http.StatusBadRequest, "args not allowed", "namespace or container name or pod name required")
		return
	}

	logs.GetLogs(data, c)
}

func ws_get_ssh_html(c *gin.Context) {
	// var data model.CoreV1
	// if err := c.BindJSON(&data); err != nil {
	// 	httpclient.SendErrorMessage(c, http.StatusBadRequest, "not found", err.Error())
	// 	return
	// }

	data := &model.CoreV1{}
	data.Namespace = c.Param("namespace")
	data.ContainerName = c.Param("container")
	data.PodName = c.Param("pod")

	c.HTML(200, "ws/ssh.html", gin.H{
		"Host":      c.Request.Host,
		"Namespace": data.Namespace,
		"Container": data.ContainerName,
		"Pod":       data.PodName,
	})
}

func ws_get_ssh(c *gin.Context) {
	// var data model.CoreV1
	// if err := c.BindJSON(&data); err != nil {
	// 	httpclient.SendErrorMessage(c, http.StatusBadRequest, "not found", err.Error())
	// 	return
	// }

	data := &model.CoreV1{}
	data.Namespace = c.Param("namespace")
	data.ContainerName = c.Param("container")
	data.PodName = c.Param("pod")

	ssh.WsHandler(data, c)
}
