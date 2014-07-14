package cache

import (
	"container/list"
)

type Cacheable interface {
	Key() string
	Size() int //bytes
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
		delete(c.items, item.cacheable.Key())
		c.capacity += int64(item.cacheable.Size())
	}
}

type LRUCacheItem struct {
	cacheable Cacheable
	listElement *list.Element
}

func (c *LRUCache) Get(key string) Cacheable {
	item, exists := c.items[key]
	if exists == false {
		return nil
	}
	c.promote(item)
	return item.cacheable
}

func (c *LRUCache) Set(cacheable Cacheable) bool {
	if c.capacity < int64(cacheable.Size()) {
		c.prune()
	}

	//still not enough room, fail
	if c.capacity < int64(cacheable.Size()) { return false }

	item, exists := c.items[cacheable.Key()]
	if exists {
		item.cacheable = cacheable
		c.promote(item)
	} else {
		item = &LRUCacheItem{cacheable: cacheable,}
		item.listElement = c.list.PushFront(item)
		c.items[cacheable.Key()] = item
		c.capacity -= int64(cacheable.Size())
	}
	return true
}

func (c *LRUCache) Flush(cacheable Cacheable) bool {
	_, exists := c.items[cacheable.Key()]
	if exists {
		delete(c.items, cacheable.Key())
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



