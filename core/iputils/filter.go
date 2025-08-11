package iputils

import (
	"bufio"
	_ "embed"
	"net/netip"
	"strings"
	"sync"

	"github.com/gaissmai/bart"

	"github.com/medama-io/medama/util/logger"
)

//go:embed presets/abusive_ips.txt
var abusiveIPsData string

//go:embed presets/tor_exit_nodes.txt
var torExitNodesData string

type IPFilter struct {
	mu sync.RWMutex

	manualIPs    map[netip.Addr]struct{}
	torExitNodes map[netip.Addr]struct{}

	// Use special table to store CIDR ranges instead of individual IPs.
	abusiveIPs bart.Lite

	blockAbusiveIPs   bool
	blockTorExitNodes bool
}

func NewIPFilter() *IPFilter {
	f := &IPFilter{
		manualIPs:    make(map[netip.Addr]struct{}),
		torExitNodes: make(map[netip.Addr]struct{}),

		abusiveIPs: bart.Lite{},
	}

	f.loadPresets()

	return f
}

func (f *IPFilter) loadPresets() {
	f.mu.Lock()
	defer f.mu.Unlock()

	l := logger.Get()
	count := 0

	// Load abusive IPs. Firehol only provides CIDR ranges, so we need to parse them differently.
	scanner := bufio.NewScanner(strings.NewReader(abusiveIPsData))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if prefix, err := netip.ParsePrefix(line); err == nil {
				// Add the CIDR prefix to the Lite table
				f.abusiveIPs.Insert(prefix)
				count++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		l.Error().Err(err).Msg("error reading abusive IPs data")
	}

	l.Debug().Int("count", count).Msg("loaded abusive ip prefixes from presets")
	count = 0

	// Load Tor exit nodes.
	scanner = bufio.NewScanner(strings.NewReader(torExitNodesData))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			if addr, err := netip.ParseAddr(line); err == nil {
				f.torExitNodes[addr] = struct{}{}
				count++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		l.Error().Err(err).Msg("error reading tor exit nodes data")
	}

	l.Debug().Int("count", count).Msg("loaded tor exit nodes from presets")
}

func (f *IPFilter) LoadFromCommaSeparated(ips string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.manualIPs = make(map[netip.Addr]struct{})
	if ips == "" {
		return
	}

	for _, ipStr := range strings.Split(ips, ",") {
		if addr, err := netip.ParseAddr(strings.TrimSpace(ipStr)); err == nil {
			f.manualIPs[addr] = struct{}{}
		}
	}
}

func (f *IPFilter) SetBlockAbusiveIPs(block bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blockAbusiveIPs = block
}

func (f *IPFilter) SetBlockTorExitNodes(block bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.blockTorExitNodes = block
}

func (f *IPFilter) HasIP(ip netip.Addr) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check manual IPs
	if _, exists := f.manualIPs[ip]; exists {
		return true
	}

	// Check abusive IPs if enabled
	if f.blockAbusiveIPs {
		if f.abusiveIPs.Contains(ip) {
			return true
		}
	}

	// Check Tor exit nodes if enabled
	if f.blockTorExitNodes {
		if _, exists := f.torExitNodes[ip]; exists {
			return true
		}
	}

	return false
}
