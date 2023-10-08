package util

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

// This is a simple cache implementation that uses a map to store entries.
// Since we only expect to store a few entries (a couple sesion ids), this should be enough.

// Cache stores arbitrary data with expiration time.
type Cache struct {
	items sync.Map
	close chan struct{}
}

// An item represents arbitrary data with expiration time.
type item struct {
	data    interface{}
	expires int64
}

var (
	ErrCacheMiss   = errors.New("cache miss")
	ErrCacheExpire = errors.New("cache expired")
	ErrInvalidCast = errors.New("error casting cache item")
)

// NewCache creates a new cache that asynchronously cleans
// expired entries after the given time passes.
func NewCache(ctx context.Context, cleaningInterval time.Duration) *Cache {
	cache := &Cache{
		close: make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(cleaningInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				cache.items.Range(func(key, value interface{}) bool {
					item, ok := value.(item)
					if !ok {
						slog.ErrorContext(ctx, "error casting cache item", slog.String("error", ErrInvalidCast.Error()))
						return false
					}

					if item.expires > 0 && now > item.expires {
						cache.items.Delete(key)
					}

					return true
				})

			case <-cache.close:
				return
			}
		}
	}()

	return cache
}

// Get gets the value for the given key.
func (c *Cache) Get(ctx context.Context, key interface{}) (interface{}, error) {
	obj, exists := c.items.Load(key)

	if !exists {
		slog.DebugContext(ctx, "cache miss", slog.String("error", ErrCacheMiss.Error()))
		return nil, ErrCacheMiss
	}

	item, ok := obj.(item)
	if !ok {
		slog.DebugContext(ctx, "error casting cache item", slog.String("error", ErrInvalidCast.Error()))
		return nil, ErrInvalidCast
	}

	if item.expires > 0 && time.Now().UnixNano() > item.expires {
		slog.DebugContext(ctx, "cache expired", slog.String("error", ErrCacheExpire.Error()))
		return nil, ErrCacheExpire
	}

	return item.data, nil
}

// Has checks if the cache has the given key.
func (c *Cache) Has(ctx context.Context, key interface{}) (bool, error) {
	_, err := c.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) || errors.Is(err, ErrCacheExpire) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Set sets a value for the given key with an expiration duration.
// If the duration is 0 or less, it will be stored forever.
func (c *Cache) Set(key interface{}, value interface{}, duration time.Duration) {
	c.items.Store(key, item{
		data:    value,
		expires: time.Now().Add(duration).UnixNano(),
	})
}

// Range calls f sequentially for each key and value present in the cache.
// If f returns false, range stops the iteration.
func (c *Cache) Range(ctx context.Context, f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item, ok := value.(item)
		if !ok {
			slog.ErrorContext(ctx, "error casting cache item", slog.String("error", ErrInvalidCast.Error()))
			return false
		}

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	c.items.Range(fn)
}

// Delete deletes the key and its value from the cache.
func (c *Cache) Delete(key interface{}) {
	c.items.Delete(key)
}

// Close closes the cache and frees up resources.
func (c *Cache) Close() {
	c.close <- struct{}{}
	c.items = sync.Map{}
}
