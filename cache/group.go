package cache

import (
	"fmt"
	"gocache/cache/bloom"
	"gocache/cache/stats"
	"gocache/peer"
	"sync"
	"golang.org/x/sync/singleflight"
)

type Getter interface {
	Get(key string) (string, error)
}

type GetterFunc func(key string) (string, error)

func (f GetterFunc) Get(key string) (string, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	maincache *Cache
	hotcache  *Cache
	peer      peer.PeerPicker
	loader    *singleflight.Group
	bloom     *bloom.Bloom
	stats     *stats.Stats
}

var (
	RW     sync.RWMutex
	Groups = make(map[string]*Group)
)

func NewGroup(name string, size int, getter Getter) *Group {
	m := &Group{
		name:      name,
		getter:    getter,
		maincache: NewCache(size),
		hotcache:  NewCache(size / 8),
		loader:    &singleflight.Group{},
		bloom:     bloom.New(1<<20, 3),
		stats:     &stats.Stats{},
	}
	RW.Lock()
	Groups[name] = m
	defer RW.Unlock()
	return m
}


func GetGroup(name string) *Group {
	RW.RLock()
	defer RW.RUnlock()
	return Groups[name]
}


func (g *Group) Get(key string) (string, error) {
	g.stats.RecordRequest()
	if key == "" {
		return "", fmt.Errorf("key is required")
	}

	// 布隆过滤器
	if g.bloom != nil && !g.bloom.Contains(key) {
		return "", fmt.Errorf("key not exist")
	}
	if v, ok := g.hotcache.Get(key); ok {
		return v, nil
	}
	if v, ok := g.maincache.Get(key); ok {
		return v, nil
	}
	g.stats.RecordMiss()
	return g.load(key)
}


func (g *Group) Bloom() *bloom.Bloom {
	return g.bloom
}


func (g *Group) load(key string) (string, error) {

	v, err, _ := g.loader.Do(key, func() (interface{}, error) {
		if g.peer != nil {
			if peer, ok := g.peer.PickPeer(key); ok {
				value, err := peer.Get(g.name, key)
				if err == nil {
					g.hotcache.Add(key, value)
					return value, nil
				}
				fmt.Println("peer err:", err)
			}
		}

		return g.getLocally(key)
	})

	if err != nil {
		return "", err
	}

	return v.(string), nil
}

func (g *Group) getLocally(key string) (string, error) {
	g.stats.RecordDBLoad()
	value, err := g.getter.Get(key)
	if err != nil {
		return "", err
	}
	g.maincache.Add(key, value)
	if g.bloom != nil {
		g.bloom.Add(key)
	}
	return value, nil
}

func (g *Group) RegisterPeer(peer peer.PeerPicker) {
	if g.peer != nil {
		panic("RegisterPeers called more than once")
	}
	g.peer = peer
}
func (g *Group) Stats() stats.Stats {
	return stats.Stats{
		Requests: g.stats.Requests,
		Hits:     g.stats.Hits,
		Misses:   g.stats.Misses,
		DBLoads:  g.stats.DBLoads,
	}
}
