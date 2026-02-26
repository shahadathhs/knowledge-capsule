package middleware

import (
	"net/http"
	"strings"
)

// CORS returns a middleware that adds CORS headers based on allowed origins.
// If allowedOrigins is nil or empty, no CORS headers are added (restrictive).
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	originSet := make(map[string]bool)
	for _, o := range allowedOrigins {
		if o != "" {
			originSet[o] = true
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && originSet[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// IsOriginAllowed checks if the origin is in the allowed list.
// Used by WebSocket upgrader CheckOrigin.
func IsOriginAllowed(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return true
	}
	for _, o := range allowedOrigins {
		if o != "" && (o == "*" || strings.EqualFold(o, origin)) {
			return true
		}
	}
	return false
}
