package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type cacheElement struct {
	key          string
	value        interface{}
	lastAccessed time.Time // not really used, but we may wish to time out elements later
}

type LRUCache struct {
	mutex    sync.Mutex
	lookup   map[string]*list.Element
	list     *list.List
	size     uint
	capacity uint
	requests uint64
	hits     uint64
}

type Stats struct {
	Size     uint
	Capacity uint
	Usage    float32
	Hits     uint64
	Requests uint64
	HitRate  float32
}

func (s *Stats) String() string {
	return fmt.Sprintf("Caching %d of %d elements (%d%%), Fullfilled %d of %d requests (%d%%)", s.Size, s.Capacity, int(100*s.Usage), s.Hits, s.Requests, int(100*s.HitRate))
}

func NewLRUCache(capacity uint) *LRUCache {
	return &LRUCache{
		lookup:   make(map[string]*list.Element),
		list:     list.New(),
		capacity: capacity,
	}
}

func (lru *LRUCache) Stats() *Stats {
	stats := new(Stats)
	stats.Size = lru.size
	stats.Capacity = lru.capacity
	if lru.capacity > 0 {
		stats.Usage = float32(lru.size) / float32(lru.capacity)
	}
	stats.Hits = lru.hits
	stats.Requests = lru.requests
	if stats.Requests > 0 {
		stats.HitRate = float32(lru.hits) / float32(lru.requests)
	}
	return stats
}

// Adds an element to
func (lru *LRUCache) Add(key string, value interface{}) {

	if lru.Contains(key) {
		lru.Delete(key)
	}

	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	cacheEl := &cacheElement{key, value, time.Now()}
	listEl := lru.list.PushFront(cacheEl)
	lru.lookup[key] = listEl
	lru.size++
	lru.ensureCapacity()
}

func (lru *LRUCache) Get(key string) (interface{}, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	lru.requests++

	listEl := lru.lookup[key]
	if listEl == nil {
		return nil, false
	}

	lru.moveToFront(listEl)
	cacheItem := listEl.Value.(*cacheElement)

	lru.hits++

	return cacheItem.value, true
}

func (lru *LRUCache) Contains(key string) bool {
	return lru.lookup[key] != nil
}

func (lru *LRUCache) Delete(key string) bool {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	listEl := lru.lookup[key]
	if listEl == nil {
		return false
	}

	lru.list.Remove(listEl)
	delete(lru.lookup, key)
	lru.size--
	return true
}

// ensures that our cache size has not outgrown its capacity
func (lru *LRUCache) ensureCapacity() {
	if lru.size > lru.capacity {
		listEl := lru.list.Back()
		lru.list.Remove(listEl)
		cacheEl := listEl.Value.(*cacheElement)
		delete(lru.lookup, cacheEl.key)
		lru.size--
	}
}

// moves an element to the front of the list and updates its last accessed time
func (lru *LRUCache) moveToFront(listEl *list.Element) {
	lru.list.MoveToFront(listEl)
	cacheEl := listEl.Value.(*cacheElement)
	cacheEl.lastAccessed = time.Now()
}
