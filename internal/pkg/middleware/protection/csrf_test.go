package protection

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

var testLogger *slog.Logger
var testConfig *config.Config

func init() {
	testLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	testConfig = config.LoadConfig("../../config/config.yaml", testLogger)
}
func TestCsrfMiddleware(t *testing.T) {
	type args struct {
		method string
		token  string
		cookie *http.Cookie
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Test_Safe_Method",
			args: args{
				method: http.MethodGet,
				token:  "",
				cookie: &http.Cookie{
					Name:     testConfig.AuthHandler.Csrf.CsrfCookie,
					Secure:   true,
					Value:    "",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Test_Unsafe_Method_No_Token",
			args: args{
				method: http.MethodPost,
				token:  "",
				cookie: &http.Cookie{
					Name:     testConfig.AuthHandler.Csrf.CsrfCookie,
					Secure:   true,
					Value:    "",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Test_No_Cookie",
			args: args{
				method: http.MethodPost,
				token:  "",
				cookie: &http.Cookie{
					Name:     "",
					Secure:   true,
					Value:    "",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Test_Cookie_Token_Not_Match",
			args: args{
				method: http.MethodPost,
				token:  "9f6eae16-9e58-481d-b4a1-ec96e3f4cb93",
				cookie: &http.Cookie{
					Name:     testConfig.AuthHandler.Csrf.CsrfCookie,
					Secure:   true,
					Value:    "9f6eae16-9e58-481d-a5a6-ec96e3f4cb93",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				},
			},
			wantStatus: http.StatusForbidden,
		},
		{
			name: "Test_Success",
			args: args{
				method: http.MethodPost,
				token:  "9f6eae16-9e58-481d-b4a1-ec96e3f4cb93",
				cookie: &http.Cookie{
					Name:     testConfig.AuthHandler.Csrf.CsrfCookie,
					Secure:   true,
					Value:    "9f6eae16-9e58-481d-b4a1-ec96e3f4cb93",
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				},
			},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(tt.args.method, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()
			handler(res, req)
			req.Header.Add("X-Csrf-Token", tt.args.token)

			req.AddCookie(tt.args.cookie)
			mw := CreateCsrfMiddleware(testLogger, testConfig.AuthHandler.Csrf)
			mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
			resp := res.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

		})
	}
}
