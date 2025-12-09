package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	datastore  map[string]cacheEntry
	mu         sync.RWMutex
	expiration time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(expiration time.Duration) *Cache {
	c := &Cache{
		datastore:  make(map[string]cacheEntry),
		expiration: expiration,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.expiration)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}

}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, entry := range c.datastore {
		if now.Sub(entry.createdAt) > c.expiration {
			delete(c.datastore, key)
		}
	}

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.datastore[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.datastore[key]
	if !ok {
		return nil, false
	}

	if time.Since(entry.createdAt) > c.expiration {
		return nil, false
	}

	return entry.val, true
}
