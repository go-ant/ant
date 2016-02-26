package cache

import (
	"sync"
	"time"
)

type memoryItem struct {
	object     interface{}
	expiration *time.Time
}

type memoryCacher struct {
	sync.RWMutex
	items    map[string]*memoryItem
	interval int
}

// newMemoryCacher creates and returns a new memory cacher.
func newMemoryCacher() *memoryCacher {
	return &memoryCacher{items: make(map[string]*memoryItem)}
}

// Expired returns true if the item has expired.
func (c *memoryItem) Expired() bool {
	if c.expiration == nil {
		return false
	}
	return c.expiration.Before(time.Now())
}

func (c *memoryCacher) Set(key string, obj interface{}, expire ...time.Duration) {
	c.Lock()
	defer c.Unlock()

	var e *time.Time
	if len(expire) > 0 {
		if expire[0] == 0 {
			expire[0] = DefaultExpiration
		}
		t := time.Now().Add(expire[0])
		e = &t
	}

	c.items[key] = &memoryItem{
		object:     obj,
		expiration: e,
	}
}

func (c *memoryCacher) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if item.Expired() {
		go c.Delete(key)
		return nil, false
	}
	return item.object, true
}

func (c *memoryCacher) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *memoryCacher) IsExist(key string) bool {
	c.RLock()
	defer c.RUnlock()
	_, found := c.items[key]
	return found
}

func (c *memoryCacher) Flush() {
	c.Lock()
	defer c.Unlock()
	c.items = make(map[string]*memoryItem)
}

func (c *memoryCacher) StartGC(opt Options) {
	c.interval = opt.Interval
	go c.startGC()
}

func (c *memoryCacher) startGC() {
	if c.interval < 1 {
		return
	}
	if c.items != nil {
		for key, _ := range c.items {
			c.checkExpiration(key)
		}
	}
	time.AfterFunc(time.Duration(c.interval)*time.Second, func() { c.startGC() })
}

func (c *memoryCacher) checkExpiration(key string) {
	c.Lock()
	defer c.Unlock()

	item, found := c.items[key]
	if !found {
		return
	}
	if item.Expired() {
		delete(c.items, key)
	}
}
