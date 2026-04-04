package middleware

import (
	"net/http"
	"strings"
)

// CORS returns middleware that adds Cross-Origin Resource Sharing headers.
// If corsOrigin is "*", all origins are allowed (development mode).
// Multiple origins can be separated by commas.
func CORS(corsOrigin string) func(http.Handler) http.Handler {
	allowAll := corsOrigin == "*"
	var allowedSet map[string]bool
	if !allowAll && corsOrigin != "" {
		allowedSet = make(map[string]bool)
		for _, o := range strings.Split(corsOrigin, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedSet[o] = true
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if origin != "" && allowedSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
