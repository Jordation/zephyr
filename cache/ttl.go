package cache

import (
	"sync"
	"time"
)

type TTLCache[K comparable, T any] interface {
	Set(K, T)
	Get(K) (T, bool)
	Start()
}

type cacheItem[T any] struct {
	lastMod int64
	value   T
}

type ttlCache[K comparable, T any] struct {
	cache     map[K]*cacheItem[T]
	m         *sync.Mutex
	t         *time.Ticker
	ttl       int64 // nanos
	lastCheck time.Time
}

func NewTTL[K comparable, T any](ttl int64) TTLCache[K, T] {
	c := &ttlCache[K, T]{
		ttl:   ttl,
		cache: make(map[K]*cacheItem[T]),
		m:     &sync.Mutex{},
	}
	return c
}

// Start is blocking
func (c *ttlCache[K, T]) Start() {
	t := time.NewTicker(time.Second * 5)
	for ; true; <-t.C {
		c.m.Lock()
		now := time.Now().UnixNano()

		for k, item := range c.cache {
			if now-item.lastMod >= c.ttl {
				delete(c.cache, k)
			}
		}

		c.m.Unlock()
	}
}

func (c *ttlCache[K, T]) Get(key K) (T, bool) {
	c.m.Lock()
	defer c.m.Unlock()

	v, ok := c.cache[key]
	v.lastMod = time.Now().UnixNano()
	return v.value, ok
}

func (c *ttlCache[K, T]) Set(key K, value T) {
	c.m.Lock()
	defer c.m.Unlock()

	c.cache[key] = &cacheItem[T]{
		lastMod: time.Now().UnixNano(),
		value:   value,
	}
}
