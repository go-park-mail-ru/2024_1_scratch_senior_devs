package jwt

import (
	"io"
	"net/http"
)

func ReadAndCloseBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			_, _ = io.ReadAll(r.Body)
			defer r.Body.Close()
		}()

		next.ServeHTTP(w, r)
	})
}
