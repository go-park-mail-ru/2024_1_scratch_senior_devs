package cookie

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/middleware/jwt"
	"net/http"
	"time"
)

func GenTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     jwt.JwtCookie,
		Secure:   false,
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
		Path:     "/",
	}
}

func DelTokenCookie() *http.Cookie {
	return &http.Cookie{
		Name:   jwt.JwtCookie,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
}
