package middlewares

import (
	"net/http"
	"path"
	"strings"

	"github.com/medama-io/medama/util/logger"
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
	})

	// Create a default CORS handler
	defaultCORS := cors.Default()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uPath := path.Clean(r.URL.Path)
			log := logger.Get()

			if allowedOrigins != nil && strings.HasPrefix(uPath, "/api") && !strings.HasPrefix(uPath, "/api/event") {
				// Apply modified CORS headers for API routes.
				log.Debug().Str("allowed_origins", strings.Join(allowedOrigins, ",")).Str("path", uPath).Msg("Applying custom CORS")
				customCORS.Handler(next).ServeHTTP(w, r)
			} else {
				// Apply default CORS headers
				defaultCORS.Handler(next).ServeHTTP(w, r)
			}
		})
	}
}
