package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	interval time.Duration

	mutex   *sync.Mutex
	entries map[string]cacheEntry
}

func NewCache(interval time.Duration) Cache {
	var mutex sync.Mutex

	cache := Cache{
		interval: interval,
		mutex:    &mutex,
		entries:  make(map[string]cacheEntry),
	}

	go func() {
		ticker := time.NewTicker(interval)

		for {
			<-ticker.C
			cache.ReapLoop()
		}
	}()

	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	entry, ok := cache.entries[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (cache *Cache) ReapLoop() {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	var toDelete []string

	for key, entry := range cache.entries {
		if entry.createdAt.Add(cache.interval).Before(time.Now()) {
			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(cache.entries, key)
	}
}
