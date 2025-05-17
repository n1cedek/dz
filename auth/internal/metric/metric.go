package metric

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	nameSpace = "my_space"
	appName   = "my_app"
)

type Metrics struct {
	requestCounter        prometheus.Counter
	responseCounter       *prometheus.CounterVec
	histogramResponseTime *prometheus.HistogramVec
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: nameSpace,
				Subsystem: "grpc",
				Name:      appName + "_request_total",
				Help:      "Количество запросов к серверу",
			},
		),
		responseCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: nameSpace,
				Subsystem: "grpc",
				Name:      appName + "_response_total",
				Help:      "Количество ответов от сервера",
			},
			[]string{"status", "method"},
		),
		histogramResponseTime: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: nameSpace,
				Subsystem: "grpc",
				Name:      appName + "_histogram_response_time_seconds",
				Help:      "Время ответа от сервера",
				Buckets:   prometheus.ExponentialBuckets(0.0001, 2, 16),
			}, []string{"status"},
		),
	}
	prometheus.MustRegister(metrics.requestCounter)
	prometheus.MustRegister(metrics.responseCounter)
	prometheus.MustRegister(metrics.histogramResponseTime)
	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
func IncResponseCounter(status, method string) {
	metrics.responseCounter.WithLabelValues(status, method).Inc()
}
func IncHistogramCounter(status string, time float64) {
	metrics.histogramResponseTime.WithLabelValues(status).Observe(time)
}
