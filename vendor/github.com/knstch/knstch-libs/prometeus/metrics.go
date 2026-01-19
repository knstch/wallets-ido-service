package metrics

import (
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount    *kitprom.Counter
	RequestDuration *kitprom.Histogram
)

func InitBasicMetrics() {
	fieldKeys := []string{"path", "code"}

	RequestCount = kitprom.NewCounterFrom(stdprom.CounterOpts{
		Subsystem: "http",
		Name:      "request_count",
		Help:      "Number of requests",
	}, fieldKeys)

	RequestDuration = kitprom.NewHistogramFrom(stdprom.HistogramOpts{
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "Requests duration",
		Buckets:   stdprom.DefBuckets,
	}, fieldKeys)
}
