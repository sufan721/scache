package bloom

//布隆过滤器 -> 用位图来判断元素是否存在，解决 缓存击穿问题

import (
	"hash/fnv"
)


type Bloom struct {
	bitset []bool
	size   uint
	k      uint
}

func New(size uint, k uint) *Bloom {
	return &Bloom{
		bitset: make([]bool, size),
		size:   size,
		k:      k,
	}
}

func (b *Bloom) hash(data string, seed uint) uint {
	h := fnv.New32a()
	h.Write([]byte(data))
	sum := h.Sum32()
	return (uint(sum) + seed*seed) % b.size
}

func (b *Bloom) Add(key string) {
	for i := uint(0); i < b.k; i++ {
		index := b.hash(key, i)
		b.bitset[index] = true
	}
}

func (b *Bloom) Contains(key string) bool {
	for i := uint(0); i < b.k; i++ {
		index := b.hash(key, i)
		if !b.bitset[index] {
			return false
		}
	}
	return true
}
