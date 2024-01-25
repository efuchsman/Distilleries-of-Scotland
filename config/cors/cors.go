package cors

import (
	"net/http"
	"strings"
)

// Middleware for Cross-Origin Resource Sharing (CORS)
// This middleware sets the appropriate CORS headers to allow safe cross-origin requests.
// By default, it allows only HTTP GET requests.
func SetCORSHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		origin := r.Header.Get("Origin")
		allowed := []string{"GET"}

		if origin != "" {
			headers.Set("Access-Control-Allow-Origin", "*")

			if strings.HasPrefix(origin, "http") {
				headers.Set("Access-Control-Allow-Origin", origin)
			}
			headers.Set("Access-Control-Allow-Headers", "*")
			headers.Set("Access-Control-Allow-Methods", strings.Join(allowed, ", "))
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
