package referrer

import (
	"net/netip"
	"net/url"
	"strings"
)

type result struct {
	Host   string
	Group  string
	IsSpam bool
}

// Parse parses the referrer string and removes any query parameters, IP addresses and self-referencing URLs.
//
// It also checks the referrer against a list of known spammy hosts and returns a result struct containing the
// referrer host, its group name, and a spam flag.
func (r *Parser) Parse(referrer string, hostname string) (result, error) {
	referrerHost := ""
	referrerGroup := ""

	if referrer != "" {
		referrer, err := url.Parse(referrer)
		if err != nil {
			return result{}, err
		}

		referrerHost = referrer.Hostname()

		switch {
		// Remove any self-referencing hostnames.
		case referrerHost == hostname:
			referrerHost = ""

		// Filter out IP addresses from referrer.
		case isIPAddress(referrerHost):
			referrerHost = ""

		// If the referrer is a spammy host, set the spam flag to drop the request.
		case r.parser.IsSpam(referrerHost):
			return result{
				Host:   referrerHost,
				IsSpam: true,
			}, nil

		// Else, parse the referrer host to get the optional group name.
		default:
			referrerGroup = r.parser.Parse(referrerHost)
		}
	}

	return result{
		Host:   referrerHost,
		Group:  referrerGroup,
		IsSpam: false,
	}, nil
}

// isIPAddress checks if a string is a valid IP address (either IPv4 or IPv6).
func isIPAddress(host string) bool {
	// Fast path to check if the host _might_ be an IP address
	// by checking if it starts with a digit or contains a colon.
	if len(host) == 0 || (host[0] < '0' || host[0] > '9') && strings.IndexByte(host, ':') == -1 {
		return false
	}

	_, err := netip.ParseAddr(host)
	return err == nil
}
