package protection

import (
	"net/http"
	"os"
)

const CSP = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST,PUT,DELETE,GET")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization,X-Csrf-Token")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_HOST_NAME"))
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Content-Security-Policy", CSP)
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
