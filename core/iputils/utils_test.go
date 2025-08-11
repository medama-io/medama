package iputils_test

import (
	"net/http"
	"testing"

	"github.com/medama-io/medama/iputils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIP(t *testing.T) {
	tests := []struct {
		name        string
		setupReq    func() *http.Request
		expectedIP  string
		expectError bool
	}{
		{
			name: "X-Forwarded-For single IP",
			setupReq: func() *http.Request {
				req := &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "10.0.0.1:12345",
				}
				req.Header.Set("X-Forwarded-For", "192.168.1.100")
				return req
			},
			expectedIP: "192.168.1.100",
		},
		{
			name: "X-Forwarded-For multiple IPs",
			setupReq: func() *http.Request {
				req := &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "10.0.0.1:12345",
				}
				req.Header.Set("X-Forwarded-For", "192.168.1.100, 10.0.0.2, 172.16.0.1")
				return req
			},
			expectedIP: "192.168.1.100",
		},
		{
			name: "X-Real-IP header",
			setupReq: func() *http.Request {
				req := &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "10.0.0.1:12345",
				}
				req.Header.Set("X-Real-IP", "192.168.1.200")
				return req
			},
			expectedIP: "192.168.1.200",
		},
		{
			name: "CF-Connecting-IP header",
			setupReq: func() *http.Request {
				req := &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "10.0.0.1:12345",
				}
				req.Header.Set("CF-Connecting-IP", "192.168.1.100")
				return req
			},
			expectedIP: "192.168.1.100",
		},
		{
			name: "RemoteAddr fallback",
			setupReq: func() *http.Request {
				return &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "192.168.1.50:12345",
				}
			},
			expectedIP: "192.168.1.50",
		},
		{
			name: "IPv6 address",
			setupReq: func() *http.Request {
				req := &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "[2001:db8::1]:12345",
				}
				return req
			},
			expectedIP: "2001:db8::1",
		},
		{
			name: "Invalid RemoteAddr",
			setupReq: func() *http.Request {
				return &http.Request{
					Header:     make(http.Header),
					RemoteAddr: "invalid-address",
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupReq()

			ip, err := iputils.GetIP(req)

			if tt.expectError {
				require.Error(t, err, "expected an error for test case: %s", tt.name)
				return
			}

			require.NoError(t, err, "unexpected error for test case: %s", tt.name)
			assert.Equal(t, tt.expectedIP, ip.String(), "expected IP to match for test case: %s", tt.name)
		})
	}
}
