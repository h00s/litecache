package litecache

import (
	"sync"

	"github.com/go-raptor/components"
)

type LiteCache struct {
	components.Service
	store   sync.Map
	expirer *expiryQueue
}

type CacheItem struct {
	Key       string
	Value     []byte
	ExpiresAt int64
	Version   int64
}

type ExpiryEntry struct {
	ExpiresAt int64
	Key       string
	Version   int64
	Index     int
}

type expiryQueue struct {
	entries  []*ExpiryEntry
	mutex    sync.Mutex
	keyIndex map[string]int
}
