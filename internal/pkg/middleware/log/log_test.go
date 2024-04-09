package log

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(config.RequestIdContextKey)

		if requestID == nil {
			t.Errorf("requestID not set in context")
		}
	})

	middleware := LogMiddleware(handler)
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)
}
