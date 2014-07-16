package cache

import (
	"container/list"
	"fmt"
	"time"
)

type Cacheable interface {
	Key() string
	Size() int
}

type LRUCache struct {
	capacity int64 // bytes
	items map[string] *LRUCacheItem
	list *list.List
}

func (c *LRUCache) promote(item *LRUCacheItem) {
	c.list.MoveToFront(item.listElement)
}

func (c *LRUCache) prune() {
	for i := 0; i < 50; i++ {
		tail := c.list.Back()
		if tail == nil { return }
		item := c.list.Remove(tail).(*LRUCacheItem)
		delete(c.items, item.Cacheable.Key())
		c.capacity += int64(item.Cacheable.Size())
	}
}

type LRUCacheItem struct {
	Cacheable Cacheable
	Flags int
	Exptime int
	Cas_unique int64
	listElement *list.Element
}

func (c *LRUCache) Get(key string) *LRUCacheItem {
	item, exists := c.items[key]
	if exists == false {
		return nil
	}
	// Passive expiration
	if int64(item.Exptime) < time.Now().Unix() && item.Exptime != 0 {
		c.list.Remove(item.listElement)
		delete(c.items, item.Cacheable.Key())
		return nil
	}
	c.promote(item)
	return item
}

func (c *LRUCache) Set(Cacheable Cacheable) bool {
	if c.capacity < int64(Cacheable.Size()) {
		c.prune()
	}

	//still not enough room, fail
	if c.capacity < int64(Cacheable.Size()) {
		fmt.Printf("Capacity is about %d bytes, but item size is %d bytes (%d 64bit val)", c.capacity, Cacheable.Size(), int64(Cacheable.Size()))
		return false
	}

	item, exists := c.items[Cacheable.Key()]
	if exists {
		item.Cacheable = Cacheable
		c.promote(item)
	} else {
		item = &LRUCacheItem{Cacheable: Cacheable,}
		item.listElement = c.list.PushFront(item)
		c.items[Cacheable.Key()] = item
		c.capacity -= int64(Cacheable.Size())
	}
	return true
}

func (c *LRUCache) Flush(Cacheable Cacheable) bool {
	_, exists := c.items[Cacheable.Key()]
	if exists {
		delete(c.items, Cacheable.Key())
		return true
	} else { return false }
}

func New(capacity int64 /* bytes */) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items: make(map[string] *LRUCacheItem, 10000),
		list: list.New(),
	}
}



