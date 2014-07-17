package cache

import (
	"container/list"
	"time"
	"fmt"
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
	Exptime int64
	Cas_unique int64
	listElement *list.Element
}

func (c *LRUCache) Get(key string) *LRUCacheItem {
	item, exists := c.items[key]
	if exists == false {
		return nil
	}
	// Passive expiration
	if c.deleteExpired(item.Cacheable) {
		return nil
	}
	c.promote(item)
	return item
}

func (c *LRUCache) Set(Cacheable Cacheable, flags int, expiration_ts int64, cas_unique int64) bool {
	if c.capacity < int64(Cacheable.Size()) {
		c.prune()
	}

	//still not enough room, fail
	if c.capacity < int64(Cacheable.Size()) {
		return false
	}

	//TODO: CAS - Check And Set, need to handle such situation as cas_unique != 0

	item, exists := c.items[Cacheable.Key()]
	if exists {
		item.Cacheable = Cacheable
		item.Cas_unique = cas_unique // TODO: same as above
		item.Flags = flags
		item.Exptime = expiration_ts
		c.promote(item)
	} else {
		item = &LRUCacheItem{Cacheable: Cacheable, Flags: flags, Exptime: expiration_ts, Cas_unique: cas_unique}
		item.listElement = c.list.PushFront(item)
		c.items[Cacheable.Key()] = item
		c.capacity -= int64(Cacheable.Size())
	}
	return true
}

func (c *LRUCache) Flush(Cacheable Cacheable) bool {
	_, exists := c.items[Cacheable.Key()]
	if exists {
		// i really hope, that compiler knows about listElement which keeps deleted Cacheable
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

func (c *LRUCache) deleteExpired(Cacheable Cacheable) bool {
	item, exists := c.items[Cacheable.Key()]
	if exists {
		if item.Exptime < time.Now().Unix() && item.Exptime != 0 {
			delete(c.items, item.Cacheable.Key())
			return true
		} else { return false }
	}
	return false
}
