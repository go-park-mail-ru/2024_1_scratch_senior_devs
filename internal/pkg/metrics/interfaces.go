package metrics

//go:generate mockgen -source=interfaces.go -destination=mocks/mock.go

type DBMetrics interface {
	IncreaseErrors(queryName string)
	ObserveResponseTime(queryName string, observeTime float64)
}

type WSMetrics interface {
	IncreaseConnections()
	DecreaseConnections()
}
