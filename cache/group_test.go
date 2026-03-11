package cache

import "testing"

func TestGroupGet(t *testing.T) {

	db := map[string]string{
		"Tom":  "630",
		"Jack": "589",
	}

	loadCounts := make(map[string]int)

	group := NewGroup(
		"scores",
		2<<10,
		GetterFunc(func(key string) (string, error) {

			loadCounts[key]++

			if v, ok := db[key]; ok {
				return v, nil
			}

			return "", nil
		}),
	)

	for k, v := range db {

		if val, _ := group.Get(k); val != v {
			t.Fatalf("failed to get value")
		}

		if val, _ := group.Get(k); val != v {
			t.Fatalf("cache failed")
		}
	}

	if loadCounts["Tom"] != 1 {
		t.Fatalf("cache not working")
	}
}
