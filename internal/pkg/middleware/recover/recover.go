package recover

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/response"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

func CreateRecoverMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recoverLogger := logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

			defer func() {
				if err := recover(); err != nil {
					log.LogHandlerError(recoverLogger, http.StatusInternalServerError, fmt.Sprintf("%v", err))
					response.WriteErrorMessage(w, http.StatusInternalServerError, "internal server error")
					return
				}

				recoverLogger.Info("success")
			}()

			next.ServeHTTP(w, r)
		})
	}
}
