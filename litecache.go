package litecache

import (
	"container/heap"
	"time"
)

func NewLiteCache(cleanupInterval ...int) *LiteCache {
	var interval int
	if len(cleanupInterval) > 0 {
		interval = cleanupInterval[0]
	} else {
		interval = 1
	}

	expirer := &expiryQueue{
		entries:  make([]*ExpiryEntry, 0),
		keyIndex: make(map[string]int),
	}
	heap.Init(expirer)
	cache := &LiteCache{
		expirer: expirer,
	}
	go cache.startCleanup(interval)
	return cache
}

func (c *LiteCache) Set(key string, value []byte, ttl ...time.Duration) {
	var expiresAt int64
	var version int64

	c.expirer.mutex.Lock()
	if val, ok := c.store.Load(key); ok {
		version = val.(*CacheItem).Version + 1
		if idx, exists := c.expirer.keyIndex[key]; exists {
			heap.Remove(c.expirer, idx)
			delete(c.expirer.keyIndex, key)
		}
	}

	if len(ttl) > 0 && ttl[0] > 0 {
		expiresAt = time.Now().UnixNano() + ttl[0].Nanoseconds()
	} else {
		expiresAt = 0
	}

	item := &CacheItem{Key: key, Value: value, ExpiresAt: expiresAt, Version: version}
	c.store.Store(key, item)

	if expiresAt > 0 {
		entry := &ExpiryEntry{ExpiresAt: expiresAt, Key: key, Version: version}
		heap.Push(c.expirer, entry)
		c.expirer.keyIndex[key] = entry.Index
	}
	c.expirer.mutex.Unlock()
}

func (c *LiteCache) Get(key string) ([]byte, bool) {
	if val, ok := c.store.Load(key); ok {
		item := val.(*CacheItem)
		if item.ExpiresAt > 0 && time.Now().UnixNano() > item.ExpiresAt {
			c.store.Delete(key)
			return nil, false
		}
		return item.Value, true
	}
	return nil, false
}

func (c *LiteCache) startCleanup(cleanupInterval int) {
	ticker := time.NewTicker(time.Duration(cleanupInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		c.cleanup()
	}
}

func (c *LiteCache) cleanup() {
	now := time.Now().UnixNano()
	c.expirer.mutex.Lock()
	for c.expirer.Len() > 0 && c.expirer.entries[0].ExpiresAt <= now {
		entry := heap.Pop(c.expirer).(*ExpiryEntry)
		delete(c.expirer.keyIndex, entry.Key)
		if val, ok := c.store.Load(entry.Key); ok {
			cItem := val.(*CacheItem)
			if cItem.Version == entry.Version && cItem.ExpiresAt > 0 {
				c.store.Delete(entry.Key)
			}
		}
	}
	c.expirer.mutex.Unlock()
}
