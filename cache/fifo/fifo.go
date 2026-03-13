package fifo

import "container/list"

type fifonode struct {
	key   string
	value string
}

type FifoCache struct {
	Fifolist *list.List
	cache    map[string]*list.Element
	maxsize  int
}

func NewFifoCache(maxsize int) *FifoCache {
	return &FifoCache{
		Fifolist: list.New(),
		cache:    make(map[string]*list.Element),
		maxsize:  maxsize,
	}

}

func (fifo *FifoCache) Get(key string) (string, bool) {
	if ele, ok := fifo.cache[key]; ok {
		return ele.Value.(*fifonode).value, true
	}
	return "", false
}

func (fifo *FifoCache) Add(key string, value string) {
	if ele, ok := fifo.cache[key]; ok {
		fifo.Fifolist.MoveToFront(ele)
		ele.Value.(*fifonode).value = value
		return
	}
	fifo.cache[key] = fifo.Fifolist.PushFront(&fifonode{
		key:   key,
		value: value,
	})
	if len(fifo.cache) > fifo.maxsize {
		node := fifo.Fifolist.Back()
		delete(fifo.cache, node.Value.(*fifonode).key)
		fifo.Fifolist.Remove(node)
	}
}

func (fifo *FifoCache) Clear() {
	fifo.cache = make(map[string]*list.Element)
	fifo.Fifolist.Init()
}
