package lru

import "container/list"

type Cache struct {
	maxSize int
	ll      *list.List
	cache   map[string]*list.Element
}

type entry struct {
	key   string
	value string
}

func New(maxSize int) *Cache {
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
		return ele.Value.(entry).value, true
	}
	return "", false
}

func (c *Cache) Add(key string, value string) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry).value = value
	}
	c.cache[key] = c.ll.PushFront(entry{key, value})
	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(entry)
		delete(c.cache, kv.key)
	}
}
