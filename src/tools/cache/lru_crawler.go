package cache

import (
	"time"
	"errors"
	"container/list"
)

const (
	SLAB_SIZE = 1024 * 1024 // 1MiB
	MAX_SLEEP_TIME = 1000000 // (mcs) - 1 sec
)

// Structure for LRU crawler containment
type LRUCrawler struct {
	sleep_period uint64
	enabled bool
	storage *LRUCache
	ItemsChunk uint
	Slabs []int
}

// Crawler's constructor.
func NewCrawler(cache *LRUCache){
	return &LRUCrawler{
		sleep_period: 0,
		enabled: false,
		storage: cache,
		ItemsChunk: 0,
		Slabs: nil,
	}
}

func (c *LRUCrawler) SetSleep(duration int64) error {
	if duration >= 0 && duration <= MAX_SLEEP_TIME {
		c.sleep_period = uint64(duration)
		return nil
	}
	return errors.New("Value range mismatch")
}

func (c *LRUCrawler) AddSlabs(args ...int) {
	c.Slabs = append(c.Slabs, args...)
}

func (c *LRUCrawler) Enable() {
	c.enabled = true
}

func (c *LRUCrawler) Disable() {
	c.enabled = false
}

func (c *LRUCrawler) run() {
	var last_item = 0
	var last_slab = 0 // index of last elem of c.Slubs
	var current_list_elem *list.Element
	if c.storage != nil {
		current_list_elem = c.storage.list.Back()
	} else {
		return
	}

//	if c.Slabs == nil {
//		last_slab = 0
//	} else {
//		last_slab = c.Slabs[0]
//	}

	for{
		if !c.enabled || c.ItemsChunk == 0 || c.storage == nil {
			return
		}

		for i := int64(0); i < c.ItemsChunk; i ++ {
			if c.Slabs != nil {
				if last_item == SLAB_SIZE - 1 {
					last_item = 0
					if last_slab == len(c.Slabs) - 1 {
						last_slab = 0
						//current_list_elem = c.storage.list.
					}
				}
			} else {
				// all slabs
			}
		}


		time.Sleep(time.Microsecond * int64(c.sleep_period))
	}
}
