package protection

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestJwtMiddleware(t *testing.T) {
	testConfig := config.JwtConfig{
		JwtCookie: "YouNoteJWT",
	}
	jwt, err := GenJwtToken(models.User{}, time.Minute)
	if err != nil {
		testLogger.Error(err.Error())
	}

	type args struct {
		header string
		cookie *http.Cookie
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Test_NoHeader",
			args: args{header: "", cookie: &http.Cookie{
				Name:     "",
				Secure:   true,
				Value:    "",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_WrongHeaderFormat",
			args: args{header: "adfsgdsfadsdfs", cookie: &http.Cookie{
				Name:     "",
				Secure:   true,
				Value:    "",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_NoCookie",
			args: args{header: "Bearer qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     "",
				Secure:   true,
				Value:    "",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Cookie_And_Header_NotEqual",
			args: args{header: "Bearer qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     testConfig.JwtCookie,
				Secure:   true,
				Value:    "qdl",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Invalid_Token",
			args: args{header: "Bearer qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     testConfig.JwtCookie,
				Secure:   true,
				Value:    "qdlwfjjvfoepwk",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Success",
			args: args{header: "Bearer " + jwt,
				cookie: &http.Cookie{
					Name:     testConfig.JwtCookie,
					Secure:   true,
					Value:    jwt,
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				}},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()

			req.Header.Add("Authorization", tt.args.header)
			req.AddCookie(tt.args.cookie)
			handler(res, req)

			mw := CreateJwtMiddleware(testConfig)
			mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
			resp := res.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestJwtWebsocketMiddleware(t *testing.T) {
	testConfig := config.JwtConfig{
		JwtCookie: "YouNoteJWT",
	}
	jwt, err := GenJwtToken(models.User{}, time.Minute)
	if err != nil {
		testLogger.Error(err.Error())
	}

	type args struct {
		header string
		cookie *http.Cookie
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Test_NoHeader",
			args: args{header: "", cookie: &http.Cookie{
				Name:     "",
				Secure:   true,
				Value:    "",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_NoCookie",
			args: args{header: "qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     "",
				Secure:   true,
				Value:    "",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Cookie_And_Header_NotEqual",
			args: args{header: "qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     testConfig.JwtCookie,
				Secure:   true,
				Value:    "qdl",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Invalid_Token",
			args: args{header: "qdlwfjjvfoepwk", cookie: &http.Cookie{
				Name:     testConfig.JwtCookie,
				Secure:   true,
				Value:    "qdlwfjjvfoepwk",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Minute).UTC(),
				Path:     "/",
				SameSite: http.SameSiteStrictMode,
			}},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Test_Success",
			args: args{header: jwt,
				cookie: &http.Cookie{
					Name:     testConfig.JwtCookie,
					Secure:   true,
					Value:    jwt,
					HttpOnly: true,
					Expires:  time.Now().Add(time.Minute).UTC(),
					Path:     "/",
					SameSite: http.SameSiteStrictMode,
				}},
			wantStatus: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {}
			req := httptest.NewRequest(http.MethodGet, "http://www.your-domain.com/", nil)
			res := httptest.NewRecorder()

			req.Header.Add("Sec-WebSocket-Protocol", tt.args.header)
			req.AddCookie(tt.args.cookie)
			handler(res, req)

			mw := CreateJwtWebsocketMiddleware(testConfig)
			mw(http.HandlerFunc(handler)).ServeHTTP(res, req)
			resp := res.Result()
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestReadAndCloseBody(t *testing.T) {

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

			handler(res, req)

			ReadAndCloseBody(http.HandlerFunc(handler)).ServeHTTP(res, req)

			_, err := req.Body.Read(nil)
			assert.NotNil(t, err)

		})
	}
}
