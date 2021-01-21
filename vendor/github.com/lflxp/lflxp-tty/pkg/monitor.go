package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cpuTemp    prometheus.Gauge       // 数组
	hdFailures *prometheus.CounterVec // 统计
	connects   prometheus.Gauge
	wsCounts   *prometheus.CounterVec
	cmdCounts  *prometheus.CounterVec
)

func init() {
	connects = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ws_connect",
		Help: "current websocket connections",
	})

	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	wsCounts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_count",
			Help: "Send Websocket Total Numbers",
		},
		[]string{"send"},
	)
	cmdCounts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "receive_count",
			Help: "Receive Websocket Total Numbers",
		},
		[]string{"receive"},
	)
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
	prometheus.MustRegister(connects)
	prometheus.MustRegister(wsCounts)
	prometheus.MustRegister(cmdCounts)
}

func doMetrics() {
	cpuTemp.Set(65.3)
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
}
