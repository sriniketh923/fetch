package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	ServiceName = "SSS V2"
	// CollectorURL = "0.0.0.0:4317"
	// Insecure     = true
)

const (
	GoodRequestCode = "OK"
	BadRequestCode  = "FAILED"
)

type Instrumentation struct {
	ApiLatency *prometheus.HistogramVec
	ApiCalls   *prometheus.CounterVec
}

func NewInstrumentation() Instrumentation {
	p := Instrumentation{}

	p.ApiLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "response_time_ms",
		Help:    "response time of SSS API",
		Buckets: []float64{20, 50, 75, 100, 150, 200, 300, 500, 800, 1000}, //ms
	},
		[]string{"name", "url"},
	)

	p.ApiCalls = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls",
			Help: "API calls made to external urls",
		},
		[]string{"name", "url", "status_code"},
	)

	return p
}
