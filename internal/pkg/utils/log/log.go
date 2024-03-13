package log

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/log"
	"github.com/satori/uuid"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
)

func GFN() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	values := strings.Split(frame.Function, "/")

	return values[len(values)-1]
}

func GetRequestId(ctx context.Context) string {
	requestID, _ := ctx.Value(log.RequestIdContextKey).(uuid.UUID)
	return requestID.String()
}

func LogHandlerInfo(logger *slog.Logger, statusCode int, msg string) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))
	logger.Info(msg)
}

func LogHandlerError(logger *slog.Logger, statusCode int, msg string) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))
	logger.Error(msg)
}