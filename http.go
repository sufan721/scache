package scache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_scache/"

type HttpPool struct {
	self     string
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{self: self, basePath: defaultBasePath}
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	ret := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(ret) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	GroupName := ret[0]
	Key := ret[1]
	Group := GetGroup(GroupName)

	if Group == nil {
		http.Error(w, "no such group: "+GroupName, http.StatusNotFound)
		return
	}

	view, err := Group.Get(Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(view.ByteSlice())
}
