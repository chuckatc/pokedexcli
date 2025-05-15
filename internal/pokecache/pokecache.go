package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{}
	cache.entries = make(map[string]cacheEntry)
	go cache.reapLoop(interval)
	return &cache
}

func (cache *Cache) Add(key string, value []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[key] = entry
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if entry, ok := cache.entries[key]; ok {
		return entry.val, true
	}
	return []byte{}, false
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C
		now := time.Now()

		cache.mu.Lock()
		for key, entry := range cache.entries {
			if now.After(entry.createdAt.Add(interval)) {
				delete(cache.entries, key)
			}
		}
		cache.mu.Unlock()
	}
}
