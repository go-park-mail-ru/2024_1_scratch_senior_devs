package log

import (
	"context"
	"log/slog"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/satori/uuid"
	"google.golang.org/grpc"
)

type LogMiddleware struct {
	logger *slog.Logger
}

func NewGrpcLogMw(logger *slog.Logger) *LogMiddleware {
	return &LogMiddleware{
		logger: logger,
	}
}
func (m *LogMiddleware) ServerLogsInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	ctx = context.WithValue(ctx, config.LoggerContextKey, m.logger.With(slog.String("ID", uuid.NewV4().String())))
	h, err := handler(ctx, req)

	return h, err
}
