package middlewares

import (
	"net/http"
	"strings"
)

// XAPICommitMiddleware creates a middleware to apply the X-API-Commit header to all routes.
func XAPICommitMiddleware(commit string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Include all routes except for /event routes.
			// Only apply to GET requests.
			if r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/api/event") {
				w.Header().Set("X-Api-Commit", commit)
			}

			next.ServeHTTP(w, r)
		})
	}
}
