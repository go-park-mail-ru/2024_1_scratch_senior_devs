package metricsmw

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHttpMetricsMiddleware(t *testing.T) {
	//тут стоит вообще-то замокать интефейс метрик и проверять, что совершаются вызовы
	//соответствующих функций, когда надо
	metr, _ := metrics.NewHttpMetrics("test")
	tests := []struct {
		name         string
		expectedCode int
		handler      func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name:         "Test_Success",
			expectedCode: 200,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			name:         "Test_Fail",
			expectedCode: 400,
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			router := mux.NewRouter()
			router.Use(CreateHttpMetricsMiddleware(metr, logger))
			router.HandleFunc("/test", tt.handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
