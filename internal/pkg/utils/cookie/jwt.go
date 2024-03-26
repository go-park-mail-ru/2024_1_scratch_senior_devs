package cookie

import (
	"net/http"
	"time"
)

const JwtCookie = "YouNoteJWT"

func GenJwtTokenCookie(token string, expTime time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     JwtCookie,
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
		Name:     JwtCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}
