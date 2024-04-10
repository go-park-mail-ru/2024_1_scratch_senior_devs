package recover

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
func TestCreateRecoverMiddleware(t *testing.T) {

	tests := []struct {
		name string
	}{
		{
			name: "Test Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()
			handler(res, req)

			mw := CreateRecoverMiddleware(testLogger)
			mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
			resp := res.Result()
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
