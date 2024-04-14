package log

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

func CreateLogMiddleware(logger *slog.Logger, cfg config.JwtConfig) mux.MiddlewareFunc {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger = logger.With(slog.String("ID", uuid.NewV4().String()))
			ctx := context.WithValue(r.Context(), config.LoggerContextKey, logger)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
