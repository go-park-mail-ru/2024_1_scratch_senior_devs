package cookie

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"net/http"
	"time"
)

func GenJwtTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     config.JwtCookie,
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func DelJwtTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     config.JwtCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}
