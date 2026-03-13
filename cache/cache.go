package cache

import (
	"gocache/cache/lfu"
	"gocache/cache/lru"
	"sync"
)

type Policytype int

const (
	LRU Policytype = iota
	Lfu
	ARC
	FIFO
)

type Cache struct {
	mu          sync.RWMutex
	policy      Policy
	maxSize     int
	policy_type Policytype
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		maxSize: maxSize,
	}
}

func (c *Cache) Add(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.policy == nil {
		switch c.policy_type {
		case LRU:
			c.policy = lru.NewLruCache(c.maxSize)
		case Lfu:
			c.policy = lfu.NewLfuCache(c.maxSize)
		case FIFO:
			c.policy = fifo.New(c.maxSize)
		}
	}
	c.policy.Add(key, value)
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.policy == nil {
		return "", false
	}
	return c.policy.Get(key)
}
