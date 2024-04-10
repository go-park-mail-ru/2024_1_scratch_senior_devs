package protection

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsMiddleware(t *testing.T) {
	type args struct {
		next   http.Handler
		method string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_methodGet",
			args: args{
				next:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				method: http.MethodGet,
			},
		},
		{
			name: "Test_methodGet",
			args: args{
				next:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
				method: http.MethodOptions,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(tt.args.method, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()
			handler(res, req)
			CorsMiddleware(http.HandlerFunc(handler)).ServeHTTP(res, req)

		})
	}
}
