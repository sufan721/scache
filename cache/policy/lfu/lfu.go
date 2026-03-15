package lfu

import (
	"container/list"
	"time"
)

type lfunode struct {
	key    string
	val    string
	freq   int
	expire int64
}

type LfuCache struct {
	maxSize int
	minfreq int
	cache   map[string]*list.Element
	freqMap map[int]*list.List
}

func NewLfuCache(maxSize int) *LfuCache {
	if maxSize <= 0 {
		panic("maxSize must be greater than zero")
	}
	return &LfuCache{
		maxSize: maxSize,
		minfreq: 0,
		cache:   make(map[string]*list.Element),
		freqMap: make(map[int]*list.List),
	}
}

func (L *LfuCache) Add(key string, val string, ttl time.Duration) {
	if L.maxSize <= 0 {
		return
	}
	if ele, ok := L.cache[key]; ok {
		node := ele.Value.(*lfunode)
		node.val = val
		node.expire = time.Now().Add(ttl).UnixNano()
		L.updateFreq(ele)
		return
	}
	if len(L.cache) >= L.maxSize {
		L.removeMinFreq()
	}
	node := &lfunode{
		key:    key,
		val:    val,
		freq:   1,
		expire: time.Now().Add(ttl).UnixNano(),
	}
	if L.freqMap[node.freq] == nil {
		L.freqMap[node.freq] = list.New()
	}
	ele := L.freqMap[node.freq].PushFront(node)
	L.cache[key] = ele
	L.minfreq = 1
}

func (L *LfuCache) Get(key string) (string, bool) {

	if ele, ok := L.cache[key]; ok {

		node := ele.Value.(*lfunode)

		// 懒删除
		if time.Now().UnixNano() > node.expire {
			L.freqMap[node.freq].Remove(ele)
			delete(L.cache, key)
			return "", false
		}

		L.updateFreq(ele)
		return node.val, true
	}
	return "", false
}

func (L *LfuCache) updateFreq(element *list.Element) {

	node := element.Value.(*lfunode)

	oldFreq := node.freq

	L.freqMap[oldFreq].Remove(element)

	if L.freqMap[oldFreq].Len() == 0 {

		delete(L.freqMap, oldFreq)

		if L.minfreq == oldFreq {
			L.minfreq++
		}
	}

	node.freq++

	newFreq := node.freq

	if L.freqMap[newFreq] == nil {
		L.freqMap[newFreq] = list.New()
	}

	newele := L.freqMap[newFreq].PushFront(node)

	L.cache[node.key] = newele
}

func (L *LfuCache) removeMinFreq() {

	list := L.freqMap[L.minfreq]

	if list == nil {
		return
	}

	back := list.Back()

	if back != nil {

		node := back.Value.(*lfunode)

		list.Remove(back)

		delete(L.cache, node.key)

		if list.Len() == 0 {
			delete(L.freqMap, L.minfreq)
		}
	}
}

func (L *LfuCache) Len() int {
	return len(L.cache)
}

func (L *LfuCache) Clear() {
	L.cache = make(map[string]*list.Element)
	L.freqMap = make(map[int]*list.List)
	L.minfreq = 0
}
