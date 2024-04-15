package log

import (
	"context"
	"github.com/satori/uuid"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/gorilla/mux"
)

func CreateLogMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), config.LoggerContextKey, logger.With(slog.String("ID", uuid.NewV4().String())))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
