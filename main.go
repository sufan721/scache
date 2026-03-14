package main

import (
	"flag"
	"fmt"
	"gocache/cache"
	"gocache/httpserver"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) (string, error) {

			log.Println("[Getter] DB lookup key =", key)

			if v, ok := db[key]; ok {
				log.Println("[Getter] HIT DB:", key, "->", v)
				return v, nil
			}

			log.Println("[Getter] MISS DB:", key)
			return "", fmt.Errorf("key not exist")
		}))
}

func startCacheServer(addr string, addrs []string, g *cache.Group) {

	// 创建节点池
	peers := httpserver.NewHTTPPool(addr)

	// 注册所有节点
	peers.Set(addrs...)

	// 只注册一次 peerPicker
	g.RegisterPeer(peers)

	log.Println("cache node running at", addr)

	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func main() {

	var port int
	flag.IntVar(&port, "port", 8001, "cache server port")

	flag.Parse()

	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	group := createGroup()
	for k := range db {
		group.Bloom().Add(k)
	}
	startCacheServer(addrMap[port], addrs, group)
}
