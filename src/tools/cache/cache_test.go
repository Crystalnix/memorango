package cache

import (
	"testing"
	"tools"
)

func TestCacheCreationSuite1(t *testing.T){
	if New(0) != nil || New(-1000) != nil {
		t.Fatalf("Capacity is invalid.")
	}
}

func TestCacheCreationSuite2(t *testing.T){
	if New(1000) == nil {
		t.Fatalf("Unexpected nil value instead of LRUCahce instance.")
	}
}

func TestCacheSetSuite1(t *testing.T){
	cache := New(50)
	if !cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0) {
		t.Fatalf("Unexpected value.")
	}
	if cache.list.Len() == 0 {
		t.Fatalf("Error occured during setting of element.")
	}
}

func TestCacheSetSuite2(t *testing.T){
	cache := New(50)
	cache.Set(tools.NewStoredData([]byte("TEST1"), "key1"), 0, 0, 0)
	l_elem := cache.items["key1"].listElement
	cache.Set(tools.NewStoredData([]byte("TEST2"), "key2"), 0, 0, 0)
	l := cache.list.Len()
	if !cache.Set(tools.NewStoredData([]byte("CHANGED"), "key1"), 0, 0, 0) {
		t.Fatalf("Unexpected value.")
	}
	if cache.list.Len() != l {
		t.Fatalf("Error occured during updating of item.")
	}
	if l_elem != cache.list.Front() {
		t.Fatalf("Error occured during promoting of item.")
	}
}

func TestCacheSetSuite3(t *testing.T){
	cache := New(4)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0)
	l := cache.list.Len()
	cache.Set(tools.NewStoredData([]byte("TEST"), "not_key"), 0, 0, 0)
	if cache.list.Len() != l {
		t.Fatalf("Error occured during appending of exceeding item.")
	}
}

func TestCacheSetSuite4(t *testing.T){
	cache := New(4)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0)
	if cache.Set(tools.NewStoredData([]byte("HUGE AMOUNT OF DATA"), "not_key"), 0, 0, 0) {
		t.Fatalf("Error occured during appending item of unappropriate size.")
	}
	if cache.list.Len() != 0 {
		t.Fatalf("Error occured during appending of exceeding item.")
	}
}

func TestCacheGetSuite1(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0)
	res := cache.Get("key")
	if res == nil {
		t.Fatalf("Unexpected value.", res)
	}
	extr := tools.ExtractStoredData(res.Cacheable)
	if string(extr) != "TEST" {
		t.Fatalf("Wrong returned value: %s", string(extr))
	}
}

func TestCacheGetSuite2(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key1"), 0, 0, 0)
	l_elem := cache.items["key1"].listElement
	cache.Set(tools.NewStoredData([]byte("TEST"), "key2"), 0, 0, 0)
	if l_elem == cache.list.Front() {
		t.Fatalf("Wrong list element position.")
	}
	cache.Get("key1")
	if l_elem != cache.list.Front() {
		t.Fatalf("Wrong list element position after retrieving.")
	}
}

func TestCacheGetSuite3(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 1111111, 0) // Should be immediately expired
	res := cache.Get("key")
	if res != nil {
		t.Fatalf("Unexpected value.", res)
	}
}

func TestCacheGetSuite4(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0)
	res := cache.Get("key1")
	if res != nil {
		t.Fatalf("Unexpected value.", res)
	}
}

func TestCacheFlushAll(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key1"), 0, 0, 0)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key2"), 0, 0, 0)
	if cache.list.Len() != 2 {
		t.Fatalf("Error occurred during setting of elements.")
	}
	cache.FlushAll()
	if cache.list.Len() != 0 {
		t.Fatalf("Error occured during flushing all elements.")
	}

}

func TestCacheFlushItemSuite1(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key1"), 0, 0, 0)
	if cache.list.Len() == 0 {
		t.Fatalf("Error occurred during setting of element.")
	}
	if !cache.Flush("key1") {
		t.Fatalf("Unexpected result of flushing.")
	}
	if cache.list.Len() != 0 {
		t.Fatalf("The length of list still same.")
	}
}

func TestCacheFlushItemSuite2(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key1"), 0, 0, 0)
	if cache.list.Len() == 0 {
		t.Fatalf("Error occurred during setting of element.")
	}
	if cache.Flush("key2") {
		t.Fatalf("Unexpected result of flushing.")
	}
	if cache.list.Len() == 0 {
		t.Fatalf("The length of list was changed.")
	}
}

func TestCacheSetCasSuite(t *testing.T){
	cache := New(10)
	cache.Set(tools.NewStoredData([]byte("TEST"), "key"), 0, 0, 0)
	if cache.SetCas("not_key", 424242) || !cache.SetCas("key", 424242) {
		t.Fatalf("Unexpected behavior")
	}
	if cache.Get("key").Cas_unique != 424242 {
		t.Fatalf("Cas unique wasn't set")
	}
}
