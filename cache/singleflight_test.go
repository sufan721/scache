package cache

import (
	"sync"
	"testing"
)

func TestSingleflight(t *testing.T) {

	count := 0

	g := NewGroup("test", 10, GetterFunc(func(key string) (string, error) {
		count++
		return "value", nil
	}))

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {

		wg.Add(1)

		go func() {
			defer wg.Done()
			g.Get("key")
		}()

	}

	wg.Wait()

	if count > 1 {
		t.Fatal("singleflight failed")
	}
}
