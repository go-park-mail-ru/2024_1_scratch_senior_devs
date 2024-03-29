package log

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/satori/uuid"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), config.RequestIdContextKey, uuid.NewV4())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
