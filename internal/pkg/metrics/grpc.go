package metrics

import "github.com/prometheus/client_golang/prometheus"

type GrpcMetrics struct {
	HitsTotal *prometheus.CounterVec
	name      string
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
	metr.name = name
	return &metr, nil
}
func (m *GrpcMetrics) IncreaseHits(path string) {
	m.HitsTotal.WithLabelValues(path, m.name).Inc()
}
