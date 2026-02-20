package pokecache

import (
	"fmt"
	_ "fmt"
	"sync"
	"time"
)

type Cache struct {
	cache    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

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

func (c *Cache) Add(key string, val []byte) {
	fmt.Println("pokecache.Add")
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	fmt.Println("pokecache.Get")
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
