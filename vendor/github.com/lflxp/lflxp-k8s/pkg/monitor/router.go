package monitor

import (
	"fmt"

	"github.com/lflxp/lflxp-k8s/pkg/monitor/services"
	"github.com/lflxp/lflxp-k8s/pkg/proxy"
	"github.com/lflxp/lflxp-k8s/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func RegisterMonitor(router *gin.Engine) {
	grafana := viper.GetString("proxy.grafana")
	if grafana == "" {
		grafana = "http://grafana.monitoring:3000"
	}

	monitorGroup := router.Group("")
	{
		monitorGroup.GET("/api/monitor/prometheus/*action", GetPrometheus)
		// http://192.168.99.115:8002/monitor/grafana/d/x6XeSRrVk/node-exporter-nodes?orgId=1&refresh=30s
		monitorGroup.Any("/monitor/grafana/*action", proxy.NewHttpProxyByGinCustom(fmt.Sprintf("%s/monitor/grafana", grafana), nil))
		monitorGroup.GET("/api/monitor/metrics/pod/:namespace", GetMetricsPod)
		monitorGroup.GET("/api/monitor/metrics/node", GetMetricsNode)
	}
}

func GetPrometheus(c *gin.Context) {
	query := c.Request.URL.RawQuery
	action := c.Param("action")

	resp, code, err := services.GetPrometheus(action, query)
	if err != nil {
		utils.SendErrorMessage(c, code, "prometheus api error", err.Error())
		return
	}
	utils.SendSuccessMessage(c, code, resp)
}

func GetMetricsPod(c *gin.Context) {
	namespace := c.Param("namespace")

	if namespace == "" {
		utils.SendErrorMessage(c, 400, "namespace is empty", "namespace is empty")
		return
	}

	data, err := services.GetPodMetrics(namespace)
	if err != nil {
		utils.SendErrorMessage(c, 400, "error", err.Error())
		return
	}

	utils.SendSuccessMessage(c, 200, data)
}

func GetMetricsNode(c *gin.Context) {
	data, err := services.GetNodeMetrics()
	if err != nil {
		utils.SendErrorMessage(c, 400, "error", err.Error())
		return
	}

	utils.SendSuccessMessage(c, 200, data)
}
