package iputils_test

import (
	"net/netip"
	"testing"

	"github.com/medama-io/medama/iputils"
	"github.com/stretchr/testify/assert"
)

func TestLoadFromCommaSeparated(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ips      string
		expected []string
	}{
		{
			name:     "single IP",
			ips:      "192.168.1.1",
			expected: []string{"192.168.1.1"},
		},
		{
			name:     "multiple IPs",
			ips:      "192.168.1.1,10.0.0.1,172.16.0.1",
			expected: []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"},
		},
		{
			name:     "IPs with whitespace",
			ips:      "192.168.1.1, 10.0.0.1 , 172.16.0.1",
			expected: []string{"192.168.1.1", "10.0.0.1", "172.16.0.1"},
		},
		{
			name:     "empty string",
			ips:      "",
			expected: []string{},
		},
		{
			name:     "mixed valid and invalid",
			ips:      "192.168.1.1,invalid,10.0.0.1",
			expected: []string{"192.168.1.1", "10.0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			filter := iputils.NewIPFilter()

			for _, expectedIP := range tt.expected {
				ip := netip.MustParseAddr(expectedIP)

				res := filter.HasIP(ip)
				assert.False(t, res)
			}

			filter.LoadFromCommaSeparated(tt.ips)

			for _, expectedIP := range tt.expected {
				ip := netip.MustParseAddr(expectedIP)

				res := filter.HasIP(ip)
				assert.True(t, res)
			}
		})
	}
}

func TestAbusiveIPs(t *testing.T) {
	t.Parallel()

	filter := iputils.NewIPFilter()

	ip := netip.MustParseAddr("0.0.0.0")

	// Test that abusive IP exists.
	filter.SetBlockAbusiveIPs(true)
	assert.True(t, filter.HasIP(ip), "expected abusive IP to be blocked when enabled")

	// Now check it doesn't exist.
	filter.SetBlockAbusiveIPs(false)
	assert.False(t, filter.HasIP(ip), "expected abusive IP to not be blocked when disabled")
}

func TestTorExitNodes(t *testing.T) {
	t.Parallel()

	filter := iputils.NewIPFilter()

	ip := netip.MustParseAddr("171.25.193.25")

	filter.SetBlockTorExitNodes(true)
	assert.True(t, filter.HasIP(ip), "expected tor exit node IP to be blocked when enabled")

	// Now check it doesn't exist.
	filter.SetBlockTorExitNodes(false)
	assert.False(t, filter.HasIP(ip), "expected tor exit node IP to not be blocked when disabled")
}
