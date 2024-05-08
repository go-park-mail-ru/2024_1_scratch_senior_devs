package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type DatabaseMetrics struct {
	db      string
	service string
	Times   *prometheus.HistogramVec
	Errors  *prometheus.CounterVec
}

func NewDatabaseMetrics(name string, service string) (DatabaseMetrics, error) {
	metr := DatabaseMetrics{db: name, service: service}

	metr.Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name + "_db_errors_total",
			Help: "Number of total errors.",
		},
		[]string{"query", "db", "service"},
	)
	if err := prometheus.Register(metr.Errors); err != nil {
		return DatabaseMetrics{}, err
	}

	metr.Times = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name + "_db_total_times",
		},
		[]string{"query", "db", "service"},
	)
	if err := prometheus.Register(metr.Times); err != nil {
		return DatabaseMetrics{}, err
	}

	return metr, nil
}

func (m *DatabaseMetrics) IncreaseErrors(queryName string) {
	m.Errors.WithLabelValues(queryName, m.db, m.service).Inc()
}

func (m *DatabaseMetrics) ObserveResponseTime(queryName string, observeTime float64) {
	m.Times.WithLabelValues(queryName, m.db, m.service).Observe(observeTime)
}
