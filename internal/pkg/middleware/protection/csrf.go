package protection

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/config"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/cookie"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
)

var unsafeMethods = []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

func genCsrfToken() string {
	return uuid.NewV4().String()
}

func checkCsrfToken(logger *slog.Logger, w http.ResponseWriter, r *http.Request, cfg config.CsrfConfig) bool {
	headerToken := r.Header.Get("X-CSRF-Token")
	if headerToken == "" {
		log.LogHandlerError(logger, http.StatusForbidden, "empty X-CSRF-Token header")
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	csrfCookie, err := r.Cookie(cfg.CsrfCookie)
	if err != nil {
		log.LogHandlerError(logger, http.StatusForbidden, "no csrf cookie: "+err.Error())
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	if csrfCookie.Value != headerToken {
		log.LogHandlerError(logger, http.StatusForbidden, "tokens in cookie and header are different")
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	return true
}

func SetCsrfToken(w http.ResponseWriter, cfg config.CsrfConfig) {
	csrfToken := genCsrfToken()

	http.SetCookie(w, cookie.GenCsrfTokenCookie(csrfToken, cfg))
	w.Header().Set("X-Csrf-Token", csrfToken)
}

func CreateCsrfMiddleware(cfg config.CsrfConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			csrfLogger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GFN()))

			if slices.Contains(unsafeMethods, r.Method) {
				if !checkCsrfToken(csrfLogger, w, r, cfg) {
					//	return
				}
				SetCsrfToken(w, cfg)
			}

			csrfLogger.Info("success")
			next.ServeHTTP(w, r)
		})
	}
}
