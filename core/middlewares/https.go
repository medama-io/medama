package middlewares

import (
	"net"
	"net/http"
	"strings"
)

const httpsProtocol = "https://"

func HTTPSRedirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sb strings.Builder

		// Pre-allocate the buffer size for the redirect URL.
		sb.Grow(len(httpsProtocol) + len(r.Host) + len(r.URL.RequestURI()))

		sb.WriteString("https://")

		// Remove the port from the host if it exists as we are redirecting to the default HTTPS port.
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}

		sb.WriteString(host)
		sb.WriteString(r.URL.RequestURI())

		// Close old HTTP connection.
		w.Header().Set("Connection", "close")

		http.Redirect(w, r, sb.String(), http.StatusMovedPermanently)
	}
}
