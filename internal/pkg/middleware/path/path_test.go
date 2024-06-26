package path

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMiddleware(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
	res := httptest.NewRecorder()

	handler(res, req)
	PathMiddleware(http.HandlerFunc(handler)).ServeHTTP(res, req)
	assert.Equal(t, req.URL.String(), "http://www.your-domain.com")
}
