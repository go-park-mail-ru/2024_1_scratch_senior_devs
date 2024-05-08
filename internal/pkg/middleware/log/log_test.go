package log

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/stretchr/testify/assert"
)

func TestCreateLogMiddleware(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "OnlyTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))

			mw := CreateLogMiddleware(logger)
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()
			handler(res, req)

			mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
			got := log.GetLoggerFromContext(req.Context())
			assert.NotEqual(t, logger, got)
		})
	}
}
