package log

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestLogMiddleware(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "Only_Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()
			req = req.WithContext(context.Background())
			handler(res, req)

			LogMiddleware(http.HandlerFunc(handler)).ServeHTTP(res, req)
			assert.Nil(t, req.Context().Value(config.RequestIdContextKey))

		})
	}
}
