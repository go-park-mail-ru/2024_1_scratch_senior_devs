package cookie

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
)

func GenJwtTokenCookie(token string, expTime time.Time, cfg config.JwtConfig) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.JwtCookie,
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  expTime,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func DelJwtTokenCookie(cfg config.JwtConfig) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.JwtCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}
