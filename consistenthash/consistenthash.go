package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {

	m := &Map{
		replicas: replicas,
		hashMap:  make(map[int]string),
	}

	if fn != nil {
		m.hash = fn
	} else {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

func (m *Map) Add(nodes ...string) {

	for _, node := range nodes {

		for i := 0; i < m.replicas; i++ {

			hash := int(m.hash([]byte(strconv.Itoa(i) + node)))

			m.keys = append(m.keys, hash)

			m.hashMap[hash] = node
		}
	}

	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {

	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {

		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
