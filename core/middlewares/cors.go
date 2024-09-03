package middlewares

import (
	"net/http"
	"path"
	"strings"

	"github.com/rs/cors"
)

// CORSAllowedOriginsMiddleware creates a middleware to apply CORS headers based on the allowed origins.
// Typically this won't need a custom list of allowed origins as the client will be served from the same domain.
// But it is useful for development and external dashboards as we need to pass credentials from different domains.
func CORSAllowedOriginsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	// Create a CORS handler with custom options for the allowed origins
	customCORS := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		// CORS blocks the custom headers by default, so we need to allow them explicitly
		ExposedHeaders: []string{"x-api-commit"},
	})

	// Create a default CORS handler
	defaultCORS := cors.Default()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uPath := path.Clean(r.URL.Path)

			if allowedOrigins != nil && strings.HasPrefix(uPath, "/api") && !strings.HasPrefix(uPath, "/api/event") {
				// Apply modified CORS headers for API routes.
				customCORS.Handler(next).ServeHTTP(w, r)
			} else {
				// Apply default CORS headers
				defaultCORS.Handler(next).ServeHTTP(w, r)
			}
		})
	}
}
