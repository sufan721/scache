package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

// 测试 Get 方法
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("value1"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "value1" {
		t.Fatalf("cache 获取元素失败 ")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache未命中缓存")
	}
}

// 测试，当使用内存超过了设定值时，是否会触发“无用”节点的移除
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := k1 + k2 + v1 + v2
	lru := New(int64(len(cap)), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("移除老数据失败")
	}
}

// 测试回调函数能否被调用
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

}
