package lru

import (
	"testing"
)

func TestGet(t *testing.T) {

	lru := NewLruCache(2)

	lru.Add("key1", "123")
	lru.Add("key2", "456")

	if v, ok := lru.Get("key1"); !ok || v != "123" {
		t.Fatalf("cache hit key1 failed")
	}

	if v, ok := lru.Get("key2"); !ok || v != "456" {
		t.Fatalf("cache hit key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {

	lru := NewLruCache(2)

	lru.Add("key1", "1")
	lru.Add("key2", "2")
	lru.Add("key3", "3")

	if _, ok := lru.Get("key1"); ok {
		t.Fatalf("key1 should be removed")
	}
}
