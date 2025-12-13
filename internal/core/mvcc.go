package core

import (
	"sync"
)

// ConcurrentHashMap is a thread-safe hash map using sync.Map
type ConcurrentHashMap struct {
	data sync.Map
}

// NewConcurrentHashMap creates a new concurrent hash map
func NewConcurrentHashMap() *ConcurrentHashMap {
	return &ConcurrentHashMap{}
}

// Get retrieves a skip list map for a given key
func (chm *ConcurrentHashMap) Get(key string) (*ConcurrentSkipListMap, bool) {
	val, ok := chm.data.Load(key)
	if !ok {
		return nil, false
	}
	return val.(*ConcurrentSkipListMap), true
}

// ComputeIfAbsent gets or creates a skip list map for a given key
func (chm *ConcurrentHashMap) ComputeIfAbsent(key string, mappingFunc func() *ConcurrentSkipListMap) *ConcurrentSkipListMap {
	if val, ok := chm.data.Load(key); ok {
		return val.(*ConcurrentSkipListMap)
	}

	newVal := mappingFunc()
	actual, loaded := chm.data.LoadOrStore(key, newVal)
	if loaded {
		return actual.(*ConcurrentSkipListMap)
	}
	return newVal
}

// Put adds or updates a skip list map
func (chm *ConcurrentHashMap) Put(key string, value *ConcurrentSkipListMap) {
	chm.data.Store(key, value)
}

// Delete removes a key from the map
func (chm *ConcurrentHashMap) Delete(key string) {
	chm.data.Delete(key)
}

// Size returns the number of keys in the map
func (chm *ConcurrentHashMap) Size() int {
	count := 0
	chm.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// MVCC implements Multi-Version Concurrency Control
type MVCC struct {
	versions *ConcurrentHashMap
}

// NewMVCC creates a new MVCC instance
func NewMVCC() *MVCC {
	return &MVCC{
		versions: NewConcurrentHashMap(),
	}
}

// Read retrieves the value for a key at a specific version
// Returns the most recent version less than or equal to the requested version
func (mvcc *MVCC) Read(key string, version int64) (string, bool) {
	skipList, ok := mvcc.versions.Get(key)
	if !ok {
		return "", false
	}

	// Find the greatest version <= requested version
	_, value, found := skipList.FloorEntry(version)
	return value, found
}

// Pool для skip list map
var skipListMapPool = sync.Pool{
	New: func() interface{} {
		return NewConcurrentSkipListMap()
	},
}

// Write stores a value for a key at a specific version
func (mvcc *MVCC) Write(key string, value string, version int64) {
	skipList := mvcc.versions.ComputeIfAbsent(key, func() *ConcurrentSkipListMap {
		sl := skipListMapPool.Get().(*ConcurrentSkipListMap)
		// Reset skip list state
		sl.level = 1
		// Очищаем head forward pointers
		for i := 0; i < maxLevel; i++ {
			sl.head.forward[i] = nil
		}
		return sl
	})

	skipList.Put(version, value)
}

// Delete removes a key from the MVCC store
func (mvcc *MVCC) Delete(key string) {
	mvcc.versions.Delete(key)
}

// GetVersionCount returns the number of versions stored for a key
func (mvcc *MVCC) GetVersionCount(key string) int {
	skipList, ok := mvcc.versions.Get(key)
	if !ok {
		return 0
	}
	return skipList.Size()
}
