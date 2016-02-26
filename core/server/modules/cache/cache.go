package cache

import "time"

const (
	DefaultExpiration = time.Duration(0)

	MemoryStore string = "cache_memory"
	RedisStore  string = "cache_redis"
)

type (
	// Cache is the interface that operates the cache data.
	Cache interface {
		// Set adds an item to the cache, replacing any existing item.
		// If the expire is 0, the cache's default expiration time is used (30 min).
		// If it is nil, the item never expires.
		Set(key string, value interface{}, expire ...time.Duration)
		// Get gets an item from the cache. returns the item or nil, and a bool indicating
		Get(key string) (interface{}, bool)
		// Delete deletes an item from the cache.
		Delete(key string)
		// IsExist returns true if cached value exists.
		IsExist(key string) bool
		// Flush deletes all cached data.
		Flush()
		// StartGC starts GC routine based on config string settings.
		StartGC(opt Options)
	}

	Options struct {
		// Store cache store. Default is 'MemoryStore'
		Store string
		// Config stores configuration.
		Config string
		// Interval GC interval time in seconds. Default is 60.
		Interval int
	}
)

var (
	Instance Cache
)

func prepareOptions(options []Options) (opt Options) {
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Store) == 0 {
		opt.Store = MemoryStore
	}
	if opt.Interval == 0 {
		opt.Interval = 60
	}
	return
}

func NewCache(options ...Options) {
	opt := prepareOptions(options)
	if Instance == nil {
		switch opt.Store {
		case MemoryStore:
			Instance = newMemoryCacher()
		}
		Instance.StartGC(opt)
	}
}

func Set(key string, value interface{}, expire ...time.Duration) { Instance.Set(key, value, expire...) }
func Get(key string) (interface{}, bool)                         { return Instance.Get(key) }
func Delete(key string)                                          { Instance.Delete(key) }
func IsExist(key string) bool                                    { return Instance.IsExist(key) }
func Flush()                                                     { Instance.Flush() }
