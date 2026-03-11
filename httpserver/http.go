package httpserver

import (
	"fmt"
	"gocache/cache"
	"net/http"
)

type Server struct {
	Cache *cache.Cache
}

func New(c *cache.Cache) *Server {
	return &Server{
		Cache: c,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	fmt.Println(key)
	if key == "" {
		http.Error(w, "key is required", http.StatusBadRequest)
		return
	}

	value, ok := s.Cache.Get(key)
	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
	}
	fmt.Fprint(w, value)
}
