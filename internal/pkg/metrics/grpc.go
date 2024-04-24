package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type GrpcMetrics struct {
	HitsTotal *prometheus.CounterVec
	name      string
	Times     *prometheus.HistogramVec
	Errors    *prometheus.CounterVec
}

func NewGrpcMetrics(name string) (*GrpcMetrics, error) {
	var metr GrpcMetrics
	metr.HitsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_hits_total",
			Help: "Number of total hits.",
		},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metr.HitsTotal); err != nil {
		return nil, err
	}
	metr.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_errors_total",
			Help: "Number of total errors.",
		},
		[]string{"path", "service"},
	)
	if err := prometheus.Register(metr.Errors); err != nil {
		return nil, err
	}
	metr.name = name
	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_times",
		},
		[]string{"status", "path"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return nil, err
	}
	return &metr, nil
}
func (m *GrpcMetrics) IncreaseHits(path string) {
	m.HitsTotal.WithLabelValues(path, m.name).Inc()
}
func (m *GrpcMetrics) IncreaseErrors(path string) {
	m.Errors.WithLabelValues(path, m.name).Inc()
}
func (metr *GrpcMetrics) ObserveResponseTime(status int, path string, observeTime float64) {
	metr.Times.WithLabelValues(strconv.Itoa(status), path).Observe(observeTime)
}
