package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DatabaseMetrics struct {
	name   string
	Times  *prometheus.HistogramVec
	Errors *prometheus.CounterVec
}

func NewDatabaseMetrics(name string) (DatabaseMetrics, error) {
	metr := DatabaseMetrics{name: name}

	metr.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Number of total errors.",
		},
		[]string{"query", "db"},
	)
	if err := prometheus.Register(metr.Errors); err != nil {
		return DatabaseMetrics{}, err
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_total_times",
		},
		[]string{"query", "db"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return DatabaseMetrics{}, err
	}

	return metr, nil
}

func (m *DatabaseMetrics) IncreaseErrors(queryName string) {
	m.Errors.WithLabelValues(queryName, m.name).Inc()
}

func (m *DatabaseMetrics) ObserveResponseTime(queryName string, observeTime float64) {
	m.Times.WithLabelValues(queryName, m.name).Observe(observeTime)
}
