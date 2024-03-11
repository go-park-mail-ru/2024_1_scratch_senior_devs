package utils

import (
	"log/slog"
	"strconv"
)

func LogHandlerInfo(logger *slog.Logger, statusCode int, msg string) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))
	logger.Info(msg)
}

func LogHandlerError(logger *slog.Logger, statusCode int, msg string) {
	logger = logger.With(slog.String("status", strconv.Itoa(statusCode)))
	logger.Error(msg)
}
