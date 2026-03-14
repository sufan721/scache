package lru

import (
	"container/list"
	"time"
)

type Cache struct {
	maxSize int
	ll      *list.List
	cache   map[string]*list.Element
}

type node struct {
	key    string
	value  string
	expire int64
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

		n := ele.Value.(*node)

		// TTL 检查
		if time.Now().UnixNano() > n.expire {
			c.ll.Remove(ele)
			delete(c.cache, key)
			return "", false
		}

		c.ll.MoveToFront(ele)
		return n.value, true
	}

	return "", false
}

func (c *Cache) Add(key string, value string, ttl time.Duration) {

	// 如果 key 已存在
	if ele, ok := c.cache[key]; ok {

		c.ll.MoveToFront(ele)

		n := ele.Value.(*node)
		n.value = value
		n.expire = time.Now().Add(ttl).UnixNano()

		return
	}

	// 新节点
	expire := time.Now().Add(ttl).UnixNano()

	ele := c.ll.PushFront(&node{
		key:    key,
		value:  value,
		expire: expire,
	})

	c.cache[key] = ele

	if c.maxSize > 0 && c.ll.Len() > c.maxSize {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {

	ele := c.ll.Back()

	if ele != nil {

		c.ll.Remove(ele)

		kv := ele.Value.(*node)

		delete(c.cache, kv.key)
	}
}
