package metricsmw

import (
	"context"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	metrics metrics.GrpcMetrics
}

func NewGrpcMw(metrics metrics.GrpcMetrics) *GrpcMiddleware {
	return &GrpcMiddleware{
		metrics: metrics,
	}
}

func (m *GrpcMiddleware) ServerMetricsInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	h, err := handler(ctx, req)

	m.metrics.IncreaseHits(info.FullMethod)
	return h, err
}
