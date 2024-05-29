package metrics

import "github.com/prometheus/client_golang/prometheus"

type WebsocketMetrics struct {
	Connections prometheus.Gauge
}

func NewWebsocketMetrics() (*WebsocketMetrics, error) {
	metr := WebsocketMetrics{}

	metr.Connections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "ws_users_total",
		},
	)
	if err := prometheus.Register(metr.Connections); err != nil {
		return &WebsocketMetrics{}, err
	}

	return &metr, nil
}

func (m *WebsocketMetrics) IncreaseConnections() {
	m.Connections.Inc()
}

func (m *WebsocketMetrics) DecreaseConnections() {
	m.Connections.Dec()
}
