package middleware

import (
	"net/http"
	"os"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_HOST_NAME"))
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
