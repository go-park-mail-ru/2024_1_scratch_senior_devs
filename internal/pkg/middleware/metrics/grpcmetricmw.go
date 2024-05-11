package metricsmw

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	metrics *metrics.GrpcMetrics
}

func NewGrpcMw(metrics *metrics.GrpcMetrics) *GrpcMiddleware {
	return &GrpcMiddleware{
		metrics: metrics,
	}
}

func mapStatusCodes(Err string) int { //TODO: add more errors
	switch Err {
	case "user not found":
		return http.StatusUnauthorized
	case "wrong password":
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError

	}
}
func (m *GrpcMiddleware) ServerMetricsInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	h, err := handler(ctx, req)
	status := http.StatusOK
	if err != nil {
		m.metrics.IncreaseErrors(info.FullMethod)
		status = mapStatusCodes(err.Error())

	}
	m.metrics.IncreaseHits(info.FullMethod)
	m.metrics.ObserveResponseTime(status, info.FullMethod, time.Since(start).Seconds())
	return h, err
}
