package middlewares

import (
	"net/http"
	"net/netip"
	"sync/atomic"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/medama-io/medama/iputils"
	"github.com/medama-io/medama/util/logger"
)

const (
	// Total number of unique IP prefixes to track.
	cacheSize         = 65536
	defaultLimit      = 100
	defaultWindow     = 1 * time.Minute
	ipv4DefaultPrefix = 24
	ipv6DefaultPrefix = 48
)

type RateLimiter struct {
	// This thread-safe cache automatically handles LRU eviction and time-based expiration.
	// The key is the IP prefix, and the value is an atomic integer for the request count.
	visitors   *expirable.LRU[netip.Prefix, *atomic.Int64]
	limit      int64
	window     time.Duration
	ipv4Prefix int
	ipv6Prefix int
}

// NewRateLimiter creates a new RateLimiter with a single, highly optimized cache.
func NewRateLimiter(next http.Handler) http.Handler {
	rl := &RateLimiter{
		visitors:   expirable.NewLRU[netip.Prefix, *atomic.Int64](cacheSize, nil, defaultWindow),
		limit:      defaultLimit,
		window:     defaultWindow,
		ipv4Prefix: ipv4DefaultPrefix,
		ipv6Prefix: ipv6DefaultPrefix,
	}

	return rl.Middleware(next)
}

// Middleware is the HTTP middleware handler for rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.Get()

		ip, err := iputils.GetIP(r)
		if err != nil {
			log.Warn().Err(err).Msg("rate limiter: could not get client address")
			http.Error(w, "could not determine client address", http.StatusBadRequest)
			return
		}

		bits := rl.ipv6Prefix
		if ip.Is4() {
			bits = rl.ipv4Prefix
		}

		prefix, err := ip.Prefix(bits)
		if err != nil {
			log.Warn().Err(err).Msg("rate limiter: could not get prefix for client address")
			http.Error(w, "invalid client address", http.StatusBadRequest)
			return
		}
		prefix = prefix.Masked()

		// Get the current counter for this prefix, if it exists.
		counter, ok := rl.visitors.Get(prefix)
		if !ok {
			// If it doesn't exist, create a new atomic counter and add it to the cache.
			counter = &atomic.Int64{}
			rl.visitors.Add(prefix, counter)
		}

		// Atomically increment the counter and check if it exceeds the limit.
		if counter.Add(1) > rl.limit {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
