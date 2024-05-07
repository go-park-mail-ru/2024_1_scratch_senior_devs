package log

import (
	"testing"
)

func TestLogMiddleware(t *testing.T) {
	//handler := func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = r.Context().Value("logger").(*slog.Logger)
	//	w.WriteHeader(http.StatusOK)
	//}
	//
	//req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
	//res := httptest.NewRecorder()
	//handler(res, req)
	//
	//mw := CreateLogMiddleware(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	//mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
	//
	//assert.Equal(t, http.StatusOK, res.Code)
}
