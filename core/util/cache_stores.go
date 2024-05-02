package util

type CacheStore struct {
	cache map[string]struct{}
}

// Simple maps used for existence checks, e.g. validating
// all incoming hostnames against a list of existing hostnames.
func NewCacheStore() CacheStore {
	return CacheStore{
		cache: make(map[string]struct{}),
	}
}

// Add a key to the cache.
func (c *CacheStore) Add(key string) {
	c.cache[key] = struct{}{}
}

// Add all keys to the cache.
func (c *CacheStore) AddAll(keys []string) {
	for _, key := range keys {
		c.cache[key] = struct{}{}
	}
}

// Check if a key exists in the cache.
func (c *CacheStore) Has(key string) bool {
	_, ok := c.cache[key]
	return ok
}

// Remove a key from the cache.
func (c *CacheStore) Remove(key string) {
	delete(c.cache, key)
}

// Clear the cache.
func (c *CacheStore) Clear() {
	c.cache = make(map[string]struct{})
}
