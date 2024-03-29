package cookie

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"net/http"
	"time"
)

func GenCsrfTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     config.CsrfCookie,
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour).UTC(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}

func DelCsrfTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     config.CsrfCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
