// Package pokecache implements a thread-safe, in-memory cache with
// automatic entry expiration (reaping) based on a configurable interval.
package pokecache

import (
	"fmt"
	"sync"
	"time"
)

// Cache provides a thread-safe, in-memory storage for raw byte data
// with automatic expiration of entries.
type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

// NewCache creates a new Cache instance with a background reaping goroutine.
// The reaper will remove entries older than the specified interval.
func NewCache(interval time.Duration) *Cache {
	cache := &Cache{}
	cache.cache = make(map[string]cacheEntry)
	cache.interval = interval

	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			cache.reapLoop()
		}
	}()

	return cache
}

// Add inserts a new byte slice into the cache with the current timestamp.
func (c *Cache) Add(key string, val []byte) {
	fmt.Println("pokecache.Add", key)
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

// Get retrieves a byte slice from the cache. It returns the data and
// a boolean indicating if the key was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	fmt.Println("pokecache.Get", key)
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.cache[key]; ok {
		return entry.val, ok
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	fmt.Println("pokecache.reapLoop")
	c.mu.Lock()
	defer c.mu.Unlock()

	interval := c.interval
	for key, entry := range c.cache {
		if time.Since(entry.createdAt).Seconds() > float64(interval.Seconds()) {
			delete(c.cache, key)
		}
	}
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
