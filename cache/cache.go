package cache

import (
	policy2 "gocache/cache/policy"
	"gocache/cache/policy/lfu"
	"gocache/cache/policy/lru"
	"math/rand"
	"sync"
	"time"
)

type Policytype int

const (
	LRU Policytype = iota
	Lfu
)

const baseTTL = time.Minute * 10

func Gettime() time.Duration {
	return baseTTL + time.Duration(rand.Intn(5))*time.Second
}


type Cache struct {
	mu          sync.RWMutex
	policy      policy2.Policy
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
		}
	}
	c.policy.Add(key, value, Gettime())
}


func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.policy == nil {
		return "", false
	}
	return c.policy.Get(key)
}
