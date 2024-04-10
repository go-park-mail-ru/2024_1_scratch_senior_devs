package recover

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/log"
	"github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/utils/responses"
	"github.com/gorilla/mux"
)

func CreateRecoverMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recoverLogger := logger.With(slog.String("ID", log.GetRequestId(r.Context())), slog.String("func", log.GFN()))

			defer func() {
				if err := recover(); err != nil {

					print(debug.Stack())
					log.LogHandlerError(recoverLogger, http.StatusInternalServerError, fmt.Sprintf("%v", err))
					responses.WriteErrorMessage(w, http.StatusInternalServerError, errors.New("internal server error"))
					return
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
