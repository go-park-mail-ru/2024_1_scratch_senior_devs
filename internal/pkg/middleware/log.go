package middleware

import (
	"context"
	"github.com/satori/uuid"
	"net/http"
)

type RequestIdKey string

const (
	RequestIdContextKey RequestIdKey = "request_id"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestIdContextKey, uuid.NewV4())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
