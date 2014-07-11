package tools

import (
	"container/list"
)

type Cacheable interface {
	Key() string
	Size() int
}

type LRUCache struct {
	capacity int
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
		c.capacity += item.cacheable.Size()
	}
}

type LRUCacheItem struct {
	cacheable Cacheable
	listElement *list.Element
}
