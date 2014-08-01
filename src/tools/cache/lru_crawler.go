package cache

import (
	"time"
	"errors"
	"sync"
)

const (
	// Defines the top ledge of sleeping duration - 1 sec.
	MAX_SLEEP_TIME = 1000000 // (mcs) - 1 sec
)

// Structure for LRU crawler containment
type LRUCrawler struct {
	sleep_period uint32
	enabled bool
	ItemsPerRun uint
}

// Crawler's constructor.
// Receives LRU cache pointer.
func NewCrawler() *LRUCrawler {
	return &LRUCrawler{
		sleep_period: 0,
		enabled: false,
		ItemsPerRun: 0,
	}
}

// Function sets period of time which should pause main loop after each iteration.
// It receives duration parameter specified within micro seconds and required to be between 0 and max allowed sleep time;
// if it is not, function returns an error.
func (c *LRUCrawler) SetSleep(duration int) error {
	if duration >= 0 && duration <= MAX_SLEEP_TIME {
		c.sleep_period = uint32(duration)
		return nil
	}
	return errors.New("Value range mismatch")
}

// Function turns on crawler and runs main loop within thread.
func (c *LRUCache) EnableCrawler() error {
	c.Crawler.enabled = true
	var crawl_sync sync.WaitGroup
	crawl_sync.Add(1)
	go c.crawl(&crawl_sync)
	crawl_sync.Wait()
	if !c.Crawler.enabled {
		return errors.New("Failed to start crawler.")
	}
	return nil
}

// Function disables crawler by turning off its main loop.
func (c *LRUCache) DisableCrawler() {
	c.Crawler.enabled = false
}

// Function loops an infinite cycle and runs through the LRU cache by specified amount of items per loop,
// then falls asleep specified amount of time and runs again, until enabled field will be false
// whether other fields will be corrupted.
func (c *LRUCache) crawl(w_group *sync.WaitGroup) {
	defer w_group.Done()
	defer c.DisableCrawler()
	current_list_elem := c.list.Back()
	if current_list_elem == nil || c.Crawler.ItemsPerRun == 0 {
		return
	}
	w_group.Done()
	w_group.Add(1) // to make counter positive
	for {
		if !c.Crawler.enabled || c.Crawler.ItemsPerRun  == 0 {
			c.DisableCrawler()
			return
		}
		for i := uint(0); i < c.Crawler.ItemsPerRun; i ++ {
			if current_list_elem != nil {
				item := current_list_elem.Value
				if item != nil {
					if c.deleteExpired(item.(*LRUCacheItem).Cacheable) {
						c.Stats.Crawler_reclaimed ++
					}
				}
				current_list_elem = current_list_elem.Prev()
			} else {
				current_list_elem = c.list.Back()
			}
		}
		time.Sleep(time.Microsecond * time.Duration(c.Crawler.sleep_period))
	}
}
