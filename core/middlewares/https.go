package middlewares

import (
	"net"
	"net/http"
	"strings"
)

// HTTPSRedirect returns a middleware that redirects HTTP requests to HTTPS.
func HTTPSRedirect(isSSL bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldRedirect(r, isSSL) {
				redirectToHTTPS(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// HTTPSRedirectFunc is a handler function that redirects HTTP requests to HTTPS.
func HTTPSRedirectFunc(isSSL bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if shouldRedirect(r, isSSL) {
			redirectToHTTPS(w, r)
		}
	}
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	// Remove the port from the host if it exists as we are redirecting to the default HTTPS port.
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		host = r.Host
	}

	url := "https://" + host + r.URL.RequestURI()

	// Close old HTTP connection.
	w.Header().Set("Connection", "close")

	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func shouldRedirect(r *http.Request, isSSL bool) bool {
	// Skip if the request is over localhost or already using HTTPS.
	if r.TLS != nil {
		return false
	}

	isLocalhost := r.Host == "localhost" || strings.HasPrefix(r.Host, "localhost:") || strings.HasPrefix(r.Host, "127.")

	if !isSSL && isLocalhost {
		return false
	}

	return true
}
