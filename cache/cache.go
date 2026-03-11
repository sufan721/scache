package cache

import (
	"gocache/cache/lru"
	"sync"
)

type Cache struct {
	mu      sync.RWMutex
	data    *lru.Cache
	maxSize int
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
	}
}

func (c *Cache) Add(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil {
		c.data = lru.New(c.maxSize)
	}
	c.data.Add(key, value)
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.data == nil {
		return "", false
	}
	return c.data.Get(key)
}
