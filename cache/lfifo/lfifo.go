package lfifo

import "container/list"

type lfifonode struct {
	key  string
	val  string
	freq int
}

type LfifoCache struct {
	maxSize  int
	maxfreq  int
	Lififoll *list.List
	data     map[string]*list.Element
	freqMap  map[string]int
}

func NewLfifoCache(maxSize, maxfreq int) *LfifoCache {
	return &LfifoCache{
		maxSize:  maxSize,
		data:     make(map[string]*list.Element),
		freqMap:  make(map[string]int),
		Lififoll: list.New(),
		maxfreq:  maxfreq,
	}
}

func (c *LfifoCache) Get(key string) (string, bool) {
	if ele, ok := c.data[key]; ok {
		node := ele.Value.(*lfifonode)
		c.freqMap[key]++
		if c.freqMap[key] > c.maxfreq {
			c.Lififoll.MoveToBack(ele)
		}
		return node.val, true
	}
	return "", false
}
func (c *LfifoCache) Add(key string, val string) {
	if c.maxSize <= 0 {
		return
	}
	if ele, ok := c.data[key]; ok {
		node := ele.Value.(*lfifonode)

		node.val = val
		c.freqMap[key]++
		if c.freqMap[key] > c.maxfreq {
			c.Lififoll.MoveToBack(ele)
		}
	}
	if c.Lififoll.Len() >= c.maxSize {
		c.removeLowFreqNode()
	}

	node := &lfifonode{
		key:  key,
		val:  val,
		freq: 0,
	}
	ele := c.Lififoll.PushFront(node)
	c.data[key] = ele
	c.freqMap[key] = 0
}
func (lfifo *LfifoCache) removeElement(element *list.Element) {
	node := element.Value.(*lfifonode)
	lfifo.Lififoll.Remove(element)
	delete(lfifo.data, node.key)
	delete(lfifo.freqMap, node.key)
}

func (lfifo *LfifoCache) removeLowFreqNode() {
	// 从头开始查找频率最低的节点
	var lowestFreqNode *list.Element
	lowestFreq := lfifo.maxfreq + 1

	// 遍历队列，找到频率最低的节点
	for e := lfifo.Lififoll.Front(); e != nil; e = e.Next() {
		node := e.Value.(*lfifonode)
		freq := lfifo.freqMap[node.key]
		if freq < lowestFreq {
			lowestFreq = freq
			lowestFreqNode = e
		}
		// 如果找到频率为0的节点，直接移除
		if freq == 0 {
			lfifo.removeElement(e)
			return
		}
	}

	// 移除频率最低的节点
	if lowestFreqNode != nil {
		lfifo.removeElement(lowestFreqNode)
	}
}
func (lfifo *LfifoCache) Remove(key string) bool {

	if element, ok := lfifo.data[key]; ok {
		lfifo.Lififoll.Remove(element)
		delete(lfifo.data, key)
		delete(lfifo.freqMap, key)
		return true
	}
	return false
}

// Len 返回缓存大小
func (lfifo *LfifoCache) Len() int {
	return lfifo.Lififoll.Len()
}

// Clear 清空缓存
func (lfifo *LfifoCache) Clear() {
	lfifo.Lififoll = list.New()
	lfifo.data = make(map[string]*list.Element)
	lfifo.freqMap = make(map[string]int)
}
