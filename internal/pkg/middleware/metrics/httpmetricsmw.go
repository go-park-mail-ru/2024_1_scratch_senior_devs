package metricsmw

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/metrics"
	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func CreateHttpMetricsMiddleware(metr *metrics.HttpMetrics) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := NewResponseWriter(w)
			next.ServeHTTP(rw, r)
			status := http.StatusOK
			route := mux.CurrentRoute(r)
			path, _ := route.GetPathTemplate()
			statusCode := rw.statusCode
			if statusCode != http.StatusOK && statusCode != http.StatusCreated && statusCode != http.StatusNoContent {
				metr.IncreaseErrors(path, strconv.Itoa(statusCode))
			}
			metr.IncreaseHits(path, strconv.Itoa(statusCode))
			metr.ObserveResponseTime(status, path, time.Since(start).Seconds())
		})
	}
}
