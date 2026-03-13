package arc

import "container/list"

type ARCCache struct {
	maxsize  int
	p        int
	T1       *list.List
	B1       *list.List
	T2       *list.List
	B2       *list.List
	data     map[string]*list.Element
	listType map[string]int //1 T1 ,2 B1 , 3 T2, 4 B2
}

type node struct {
	key string
	val string
}

func NewARCCache(maxsize int) *ARCCache {
	return &ARCCache{
		maxsize:  maxsize,
		data:     make(map[string]*list.Element),
		listType: make(map[string]int),
		p:        0,
		T1:       list.New(),
		B2:       list.New(),
		T2:       list.New(),
		B1:       list.New(),
	}
}

func (c *ARCCache) Get(key string) (string, bool) {
	if e, ok := c.data[key]; ok {
		node := e.Value.(*node)
		litype := c.listType[key]
		switch litype {
		case 1:
			c.T1.Remove(e)
			c.listType[key] = 3
			c.T2.MoveToFront(e)
		case 3:
			c.T2.MoveToFront(e)
		}
		return node.val, true
	}
	return "", false
}
