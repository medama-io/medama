package middlewares

import (
	"net"
	"net/http"
)

// HTTPSRedirectFunc is a handler function that redirects HTTP requests to HTTPS.
func HTTPSRedirectFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Remove the port from the host if it exists as we are redirecting to the default HTTPS port.
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host // Probably has no port.
		}

		url := "https://" + host + r.URL.RequestURI()

		// Close old HTTP connection.
		w.Header().Set("Connection", "close")

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
