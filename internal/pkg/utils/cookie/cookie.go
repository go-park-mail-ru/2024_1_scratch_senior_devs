package cookie

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/jwt"
	"net/http"
	"time"
)

func GenTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     jwt.JwtCookie,
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func DelTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:     jwt.JwtCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}
