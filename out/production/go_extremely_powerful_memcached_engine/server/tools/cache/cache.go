package cache

import (
	"container/list"
)

func (c *LRUCache) Get(key string) Cacheable {
	item, exists := c.items[key]
	if exists == false { return nil }
	c.promote(item)
	return item.cacheable
}

func (c *LRUCache) Set(cacheable Cacheable) bool {
	if c.capacity < cacheable.Size() { c.prune() }

	//stil not enough room, fail
	if c.capacity < cacheable.Size() { return false }

	item, exists := c.items[cacheable.Key()]
	if exists {
		item.cacheable = cacheable
		c.promote(item)
	} else {
		item = &LRUCacheItem{cacheable: cacheable,}
		item.listElement = c.list.PushFront(item)
		c.items[cacheable.Key()] = item
		c.capacity -= cacheable.Size()
	}
	return true
}

func New(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items: make(map[string]*LRUCacheItem, 10000),
		list: list.New(),
	}
}



