package lfu

import (
	"testing"
)

func TestNewLfuCache(t *testing.T) {
	// 测试正常创建
	cache := NewLfuCache(3)
	if cache == nil {
		t.Fatal("NewLfuCache returned nil")
	}
	if cache.maxSize != 3 {
		t.Errorf("Expected maxSize = 3, got %d", cache.maxSize)
	}
	if cache.minfreq != 0 {
		t.Errorf("Expected minfreq = 0 initially, got %d", cache.minfreq)
	}
	if len(cache.cache) != 0 {
		t.Errorf("New cache should be empty, got %d items", len(cache.cache))
	}

	// 测试容量为0的情况
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for maxSize <= 0")
		}
	}()
	NewLfuCache(0)
}

func TestLfuCache_AddAndGet(t *testing.T) {
	cache := NewLfuCache(3)

	// 测试添加和获取
	cache.Add("key1", "value1")
	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("Get(key1) = %s, %v, want 'value1', true", val, ok)
	}
	if cache.Len() != 1 {
		t.Errorf("Cache length should be 1, got %d", cache.Len())
	}

	// 测试更新现有值
	cache.Add("key1", "value1-updated")
	if val, ok := cache.Get("key1"); !ok || val != "value1-updated" {
		t.Errorf("Get(key1) after update = %s, %v, want 'value1-updated', true", val, ok)
	}
	if cache.Len() != 1 {
		t.Errorf("Cache length should still be 1 after update, got %d", cache.Len())
	}
}

func TestLfuCache_EvictionWhenFull(t *testing.T) {
	cache := NewLfuCache(3)

	// 填满缓存
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	if cache.Len() != 3 {
		t.Errorf("Cache should be full with 3 items, got %d", cache.Len())
	}

	// 添加第4个key，应该触发淘汰
	cache.Add("key4", "value4")

	if cache.Len() != 3 {
		t.Errorf("Cache length should be 3 after eviction, got %d", cache.Len())
	}

	// 由于所有key都没有被访问过，应该淘汰key1（第一个添加的）
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should have been evicted (least frequently used)")
	}

	// 其他key应该还在
	if val, ok := cache.Get("key2"); !ok || val != "value2" {
		t.Errorf("key2 should still be in cache")
	}
	if val, ok := cache.Get("key3"); !ok || val != "value3" {
		t.Errorf("key3 should still be in cache")
	}
	if val, ok := cache.Get("key4"); !ok || val != "value4" {
		t.Errorf("key4 should be in cache")
	}
}

func TestLfuCache_FrequencyBasedEviction(t *testing.T) {
	cache := NewLfuCache(3)

	// 添加初始数据
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	// 增加key1的访问频率
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key1") // key1 频率=4

	// 增加key2的访问频率
	cache.Get("key2")
	cache.Get("key2") // key2 频率=3

	// key3频率=1（最低）

	// 添加新key，应该淘汰频率最低的key3
	cache.Add("key4", "value4")

	// 验证key3被淘汰
	if _, ok := cache.Get("key3"); ok {
		t.Error("key3 should have been evicted (lowest frequency)")
	}

	// 验证其他key都在
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still be in cache (highest frequency)")
	}
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 should still be in cache")
	}
	if _, ok := cache.Get("key4"); !ok {
		t.Error("key4 should be in cache")
	}
}

func TestLfuCache_ComplexFrequencyScenario(t *testing.T) {
	cache := NewLfuCache(4)

	// 添加数据
	cache.Add("A", "1")
	cache.Add("B", "2")
	cache.Add("C", "3")
	cache.Add("D", "4")

	// 访问模式: A(2次), B(3次), C(1次), D(1次)
	cache.Get("A")
	cache.Get("A")
	cache.Get("B")
	cache.Get("B")
	cache.Get("B")

	// 此时频率: A:3, B:4, C:2, D:2 (添加时频率为1，之后每个Get加1)

	// 添加新key，应该淘汰频率最低的C或D
	cache.Add("E", "5")

	// 检查C和D是否有一个被淘汰
	count := 0
	if _, ok := cache.Get("C"); ok {
		count++
	}
	if _, ok := cache.Get("D"); ok {
		count++
	}
	if count != 1 {
		t.Errorf("Exactly one of C or D should be evicted, but both or none are present")
	}

	// A和B应该都在
	if _, ok := cache.Get("A"); !ok {
		t.Error("A should still be in cache")
	}
	if _, ok := cache.Get("B"); !ok {
		t.Error("B should still be in cache")
	}
}

func TestLfuCache_SameFrequencyEviction(t *testing.T) {
	cache := NewLfuCache(3)

	// 添加数据
	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	// 所有key访问一次，使频率相同
	cache.Get("key1")
	cache.Get("key2")
	cache.Get("key3")

	// 现在所有key频率=2
	// 添加新key，应该淘汰最先添加的（LFU + FIFO）
	cache.Add("key4", "value4")

	// 应该淘汰key1
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should have been evicted (same frequency, oldest)")
	}

	// key2, key3, key4应该在
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 should still be in cache")
	}
	if _, ok := cache.Get("key3"); !ok {
		t.Error("key3 should still be in cache")
	}
	if _, ok := cache.Get("key4"); !ok {
		t.Error("key4 should be in cache")
	}
}

func TestLfuCache_Clear(t *testing.T) {
	cache := NewLfuCache(3)

	cache.Add("key1", "value1")
	cache.Add("key2", "value2")
	cache.Get("key1")

	if cache.Len() != 2 {
		t.Errorf("Cache should have 2 items before clear")
	}

	cache.Clear()

	if cache.Len() != 0 {
		t.Errorf("Cache should be empty after clear, got length: %d", cache.Len())
	}
	if cache.minfreq != 0 {
		t.Errorf("minfreq should be 0 after clear, got %d", cache.minfreq)
	}

	// 验证所有数据都被清空
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should not exist after clear")
	}
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should not exist after clear")
	}

	// 清空后可以重新使用
	cache.Add("newkey", "newvalue")
	if val, ok := cache.Get("newkey"); !ok || val != "newvalue" {
		t.Error("Cache should work after clear")
	}
}

func TestLfuCache_SingleItemCache(t *testing.T) {
	cache := NewLfuCache(1)

	cache.Add("key1", "value1")
	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("Single item cache should work")
	}
	if cache.Len() != 1 {
		t.Errorf("Single item cache length should be 1")
	}

	// 添加第二个key，应该淘汰第一个
	cache.Add("key2", "value2")
	if cache.Len() != 1 {
		t.Errorf("Cache should still have 1 item")
	}
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should be evicted")
	}
	if val, ok := cache.Get("key2"); !ok || val != "value2" {
		t.Error("key2 should be in cache")
	}
}

func TestLfuCache_MinFrequencyTracking(t *testing.T) {
	cache := NewLfuCache(3)

	cache.Add("A", "1")
	cache.Add("B", "2")
	cache.Add("C", "3")

	// 初始minfreq应该是1
	if cache.minfreq != 1 {
		t.Errorf("Initial minfreq should be 1, got %d", cache.minfreq)
	}

	// 访问A，使其频率增加
	cache.Get("A")
	if cache.minfreq != 1 {
		t.Errorf("minfreq should still be 1, got %d", cache.minfreq)
	}

	// 访问B和C，使所有key频率=2
	cache.Get("B")
	cache.Get("C")
	// 此时所有key频率=2，minfreq应该更新为2
	if cache.minfreq != 2 {
		t.Errorf("minfreq should be 2, got %d", cache.minfreq)
	}

	// 淘汰一个key
	cache.Add("D", "4")
	// 淘汰后minfreq应该重置为1
	if cache.minfreq != 1 {
		t.Errorf("minfreq should be 1 after adding new key, got %d", cache.minfreq)
	}
}

func TestLfuCache_ConcurrentBasic(t *testing.T) {
	cache := NewLfuCache(100)
	done := make(chan bool)

	// 简单并发测试
	for i := 0; i < 10; i++ {
		go func(id int) {
			key := string('A' + byte(id%26))
			cache.Add(key, "value")
			cache.Get(key)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证缓存中至少有数据
	if cache.Len() == 0 {
		t.Error("Cache should have some items after concurrent operations")
	}
}

func TestLfuCache_GetNonExistent(t *testing.T) {
	cache := NewLfuCache(3)

	if val, ok := cache.Get("nonexistent"); ok {
		t.Errorf("Get(nonexistent) = %s, %v, want '', false", val, ok)
	}

	cache.Add("key1", "value1")
	if _, ok := cache.Get("key2"); ok {
		t.Error("Get should return false for non-existent key")
	}
}

func TestLfuCache_MultipleOperations(t *testing.T) {
	cache := NewLfuCache(5)

	// 复杂操作序列
	cache.Add("A", "1")
	cache.Get("A")
	cache.Add("B", "2")
	cache.Get("A")
	cache.Add("C", "3")
	cache.Get("B")
	cache.Add("D", "4")
	cache.Get("A")
	cache.Get("C")
	cache.Add("E", "5")
	cache.Get("D")
	cache.Add("F", "6") // 应该淘汰频率最低的

	// 此时频率: A:4, B:2, C:2, D:2, E:1, F:1
	// 容量5，所以应该淘汰一个频率1的

	if cache.Len() != 5 {
		t.Errorf("Cache should have 5 items, got %d", cache.Len())
	}

	// 验证高频的A一定在
	if _, ok := cache.Get("A"); !ok {
		t.Error("A (highest frequency) should be in cache")
	}
}

func TestLfuCache_UpdateFrequencyOnUpdate(t *testing.T) {
	cache := NewLfuCache(3)

	cache.Add("key1", "value1")
	cache.Get("key1") // 频率=2

	// 更新值不应该增加频率
	cache.Add("key1", "value1-updated")

	// 添加其他key
	cache.Add("key2", "value2")
	cache.Add("key3", "value3")

	// 添加第4个key
	cache.Add("key4", "value4")

	// key1应该还在（频率高）
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still be in cache after update")
	}
	// key2或key3应该有一个被淘汰
	if _, ok1 := cache.Get("key2"); !ok1 {
		// key2被淘汰，key3应该在
		if _, ok2 := cache.Get("key3"); !ok2 {
			t.Error("One of key2 or key3 should be in cache")
		}
	}
}

// 基准测试
func BenchmarkLfuCache_Add(b *testing.B) {
	cache := NewLfuCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := string('A' + byte(i%52))
		cache.Add(key, "value")
	}
}

func BenchmarkLfuCache_Get(b *testing.B) {
	cache := NewLfuCache(1000)

	// 先填充缓存
	for i := 0; i < 1000; i++ {
		key := string('A' + byte(i%52))
		cache.Add(key, "value")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := string('A' + byte(i%52))
		cache.Get(key)
	}
}

func BenchmarkLfuCache_MixedOperations(b *testing.B) {
	cache := NewLfuCache(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := string('A' + byte(i%52))
		if i%4 == 0 {
			cache.Add(key, "value")
		} else {
			cache.Get(key)
		}
	}
}
