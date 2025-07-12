package iputils

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"
)

// GetIP extracts the client IP address from an HTTP request.
//
// It handles common cases like X-Forwarded-For, X-Real-IP headers, and the remote address.
// Returns a valid netip.Addr if a valid IP is found, otherwise returns an error.
func GetIP(r *http.Request) (netip.Addr, error) {
	// Check X-Forwarded-For header first (common for proxies).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, we want the first one (client).
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			addr, err := netip.ParseAddr(ip)
			if err == nil {
				return addr, nil
			}
		}
	}

	// Check X-Real-IP header (common for nginx).
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		addr, err := netip.ParseAddr(strings.TrimSpace(xrip))
		if err == nil {
			return addr, nil
		}
	}

	// Fall back to RemoteAddr if headers are not available.
	if r.RemoteAddr != "" {
		// RemoteAddr is typically in the format "IP:port"
		ipStr := r.RemoteAddr
		if idx := strings.LastIndex(ipStr, ":"); idx != -1 {
			ipStr = ipStr[:idx]
		}

		// Remove brackets from IPv6 addresses [2001:db8::1]:8080 -> 2001:db8::1
		ipStr = strings.TrimPrefix(ipStr, "[")
		ipStr = strings.TrimSuffix(ipStr, "]")

		addr, err := netip.ParseAddr(ipStr)
		if err == nil {
			return addr, nil
		}
	}

	// Return an invalid address if no valid ip was found.
	return netip.Addr{}, errors.New("no valid ip found in request")
}
