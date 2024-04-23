package helpers

import (
	"sync"
	"time"
)

var CacheManager *Cache

type Cache struct {
	data    map[string]cacheItem
	mu      sync.RWMutex
	cleanup chan struct{}
	stop    chan struct{}
}

type cacheItem struct {
	value     interface{}
	createdAt time.Time
	ttl       time.Duration
}

func NewCache() *Cache {
	CacheManager = &Cache{
		data:    make(map[string]cacheItem),
		cleanup: make(chan struct{}),
		stop:    make(chan struct{}),
	}
	go CacheManager.cleanupExpiredItems()
	return CacheManager
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		value:     value,
		createdAt: time.Now(),
		ttl:       ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.data[key]
	if !found {
		return nil, false
	}
	// Check if item has expired
	if time.Since(item.createdAt) > item.ttl {
		// Delete the expired item
		delete(c.data, key)
		return nil, false
	}
	return item.value, true
}

func (c *Cache) cleanupExpiredItems() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, item := range c.data {
				if time.Since(item.createdAt) > item.ttl {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		case <-c.cleanup:
			return
		}
	}
}

func (c *Cache) StopCleanup() {
	close(c.cleanup)
}
