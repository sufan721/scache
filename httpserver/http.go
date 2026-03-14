package httpserver

import (
	"fmt"
	"gocache/cache"
	"gocache/consistenthash"
	"gocache/peer"
	"log"
	"net/http"
	"strings"
	"sync"
)

const basePath = "/_gocache/"

type HTTPPool struct {
	self        string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*peer.HttpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{self: self}
}
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(50, nil)

	p.peers.Add(peers...)

	p.httpGetters = make(map[string]*peer.HttpGetter, len(peers))

	for _, peers := range peers {

		p.httpGetters[peers] = &peer.HttpGetter{
			BaseUrl: peers + basePath,
		}

	}
}
func (p *HTTPPool) PickPeer(key string) (peer.PeerGetter, bool) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {

		log.Println("Pick peer", peer)

		return p.httpGetters[peer], true
	}

	return nil, false
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, basePath) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	log.Println("request：", r.URL.Path)
	parts := strings.SplitN(
		r.URL.Path[len(basePath):],
		"/",
		2,
	)

	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := cache.GetGroup(groupName)

	if group == nil {
		http.Error(w, "group not found", http.StatusNotFound)
		return
	}

	value, err := group.Get(key)
	if err != nil {
		log.Println("group.Get error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, value)
}
