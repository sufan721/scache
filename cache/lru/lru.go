package lru

import "container/list"

type Cache struct {
	maxSize int
	ll      *list.List
	cache   map[string]*list.Element
}

type node struct {
	key   string
	value string
}

func NewLruCache(maxSize int) *Cache {
	if maxSize <= 0 {
		panic("maxSize must be greater than zero")
	}
	return &Cache{
		maxSize: maxSize,
		ll:      list.New(),
		cache:   make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(node).value, true
	}
	return "", false
}

func (c *Cache) Add(key string, value string) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*node).value = value
	}
	c.cache[key] = c.ll.PushFront(node{key, value})
	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(node)
		delete(c.cache, kv.key)
	}
}
