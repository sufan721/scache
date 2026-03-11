package cache

import (
	"fmt"
	"gocache/peer"
	"sync"
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
	peer      peer.PeerPicker
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
	if key == "" {
		return "", fmt.Errorf("key is required")
	}
	if v, ok := g.maincache.Get(key); ok {
		fmt.Println("cache hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (string, error) {
	if g.peer != nil {
		pee, ok := g.peer.PickPeer(key)
		if ok {
			value, err := pee.Get(key)
			if err == nil {
				return value, nil
			}
			fmt.Println("peer error：", err)

		}
	}
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (string, error) {
	value, err := g.getter.Get(key)
	if err != nil {
		return "", err
	}
	g.maincache.Add(key, value)
	return value, nil
}

func (g *Group) RegisterPeer(peer peer.PeerPicker) {
	if peer != nil {
		panic("RegisterPeers called more than once")
	}
	g.peer = peer
}
