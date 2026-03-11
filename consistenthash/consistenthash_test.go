package consistenthash

import "testing"

func TestHashing(t *testing.T) {

	hash := NewMap(3)

	hash.Add("NodeA", "NodeB", "NodeC")

	key := "Tom"

	node := hash.Get(key)

	if node == "" {
		t.Fatalf("node should not be empty")
	}
}
