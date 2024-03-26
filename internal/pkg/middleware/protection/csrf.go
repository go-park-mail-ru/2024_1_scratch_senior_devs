package protection

import (
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/cookie"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/gorilla/mux"
	"github.com/satori/uuid"
	"log/slog"
	"net/http"
)

func genCsrfToken() string {
	return uuid.NewV4().String()
}

func checkCsrfToken(logger *slog.Logger, w http.ResponseWriter, r *http.Request) bool {
	headerToken := r.Header.Get("X-CSRF-Token")
	if headerToken == "" {
		log.LogHandlerError(logger, http.StatusForbidden, "empty X-CSRF-Token header")
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	csrfCookie, err := r.Cookie(cookie.CsrfCookie)
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

func SetCsrfToken(w http.ResponseWriter) {
	csrfToken := genCsrfToken()

	http.SetCookie(w, cookie.GenCsrfTokenCookie(csrfToken))
	w.Header().Set("X-CSRF-Token", csrfToken)
}

func CreateCsrfMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			csrfLogger := logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
				if !checkCsrfToken(csrfLogger, w, r) {
					return
				}
				SetCsrfToken(w)
			}

			csrfLogger.Info("success")
			next.ServeHTTP(w, r)
		})
	}
}
