package middlewares

import (
	"net/http"
	"net/url"
)

// HTTPSRedirectFunc is a handler function that redirects HTTP requests to HTTPS.
func HTTPSRedirectFunc(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target := url.URL{
			Scheme:   "https",
			Host:     host,
			Path:     r.URL.Path,
			RawPath:  r.URL.RawPath,
			RawQuery: r.URL.RawQuery,
		}

		// Close old HTTP connection.
		w.Header().Set("Connection", "close")

		http.Redirect(w, r, target.String(), http.StatusMovedPermanently)
	}
}
