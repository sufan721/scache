package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gocache/cache"
	"gocache/httpserver"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

func initDB() {
	var err error
	dbConn, err = sql.Open(
		"mysql",
		"root:123456@tcp(127.0.0.1:3306)/test",
	)
	if err != nil {
		panic(err)
	}
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("mysql connected")
}

func createGroup() *cache.Group {
	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) (string, error) {
			log.Println("[MySQL] search key", key)
			var value string
			err := dbConn.QueryRow(
				"SELECT score FROM scores WHERE name=?",
				key,
			).Scan(&value)

			if err != nil {
				return "", err
			}

			return value, nil
		}))
}

func startCacheServer(addr string, addrs []string, g *cache.Group) {

	peers := httpserver.NewHTTPPool(addr)
	peers.Set(addrs...)
	g.RegisterPeer(peers)
	log.Println("cache node running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func main() {

	initDB()
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

	rows, err := dbConn.Query("SELECT name FROM scores")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		group.Bloom().Add(name)
	}
	startCacheServer(addrMap[port], addrs, group)
}
