package iputils

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"
)

var ErrInvalidIP = errors.New("no valid ip found in request")

// GetIP extracts the client IP address from an HTTP request.
// It handles common cases like X-Forwarded-For, X-Real-IP headers, and the remote address.
//
// Returns a valid netip.Addr if a valid IP is found, otherwise returns an error.
func GetIP(r *http.Request) (netip.Addr, error) {
	// Check X-Forwarded-For header first (common for proxies).
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, we want the first one (client).
		if comma := strings.IndexByte(xff, ','); comma >= 0 {
			xff = xff[:comma]
		}

		if xff != "" {
			if addr, err := netip.ParseAddr(xff); err == nil {
				return addr, nil
			}
		}
	}

	// Check CF-Connecting-IP header (Cloudflare)
	if cfip := r.Header.Get("Cf-Connecting-Ip"); cfip != "" {
		if addr, err := netip.ParseAddr(cfip); err == nil {
			return addr, nil
		}
	}

	// Check X-Real-IP header (common for nginx).
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if addr, err := netip.ParseAddr(xri); err == nil {
			return addr, nil
		}
	}

	// Check X-Client-IP (common for Apache).
	if xci := r.Header.Get("X-Client-Ip"); xci != "" {
		if addr, err := netip.ParseAddr(xci); err == nil {
			return addr, nil
		}
	}

	// Check Fastly-Client-IP header (Fastly CDN).
	if fci := r.Header.Get("Fastly-Client-Ip"); fci != "" {
		if addr, err := netip.ParseAddr(fci); err == nil {
			return addr, nil
		}
	}

	// Fall back to RemoteAddr.
	if addrPort, err := netip.ParseAddrPort(r.RemoteAddr); err == nil {
		return addrPort.Addr(), nil
	}

	// If ParseAddrPort fails, RemoteAddr might be just an IP without a port.
	if addr, err := netip.ParseAddr(r.RemoteAddr); err == nil {
		return addr, nil
	}

	return netip.Addr{}, ErrInvalidIP
}

// GetAddrList extracts a list of netip.Addr from a comma-separated string.
func GetAddrList(ips string) ([]netip.Addr, error) {
	if ips == "" {
		return []netip.Addr{}, nil
	}

	//nolint:prealloc // We don't know the number in advance.
	var addrList []netip.Addr

	for _, ipStr := range strings.Split(ips, ",") {
		ipStr = strings.TrimSpace(ipStr)

		if ipStr == "" {
			continue
		}

		addr, err := netip.ParseAddr(ipStr)
		if err != nil {
			return nil, err
		}

		addrList = append(addrList, addr)
	}

	return addrList, nil
}

// GetAddrListString converts a slice of netip.Addr to a comma-separated string.
func GetAddrListString(addrs []netip.Addr) string {
	if len(addrs) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, addr := range addrs {
		if i > 0 {
			sb.WriteString(",")
		}

		sb.WriteString(addr.String())
	}

	return sb.String()
}
