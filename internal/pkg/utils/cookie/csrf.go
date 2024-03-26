package cookie

import (
	"net/http"
	"time"
)

const CsrfCookie = "YouNoteCSRF"

func GenCsrfTokenCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     CsrfCookie,
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
		Name:     CsrfCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
