package recover

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			RecoverMiddleware(http.HandlerFunc(handler)).ServeHTTP(res, req)
			resp := res.Result()
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
