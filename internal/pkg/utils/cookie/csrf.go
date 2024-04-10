package cookie

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
)

func GenCsrfTokenCookie(token string, cfg config.CsrfConfig) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.CsrfCookie,
		Secure:   true,
		Value:    token,
		HttpOnly: true,
		Expires:  time.Now().Add(cfg.CSRFLifeTime).UTC(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}

func DelCsrfTokenCookie(cfg config.CsrfConfig) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.CsrfCookie,
		Secure:   true,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	}
}
