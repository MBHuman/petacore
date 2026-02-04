package core

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
)

// ConcurrentHashMap is a thread-safe hash map using sync.Map
type ConcurrentHashMap struct {
	data sync.Map
	size atomic.Uint64
}

// NewConcurrentHashMap creates a new concurrent hash map
func NewConcurrentHashMap() *ConcurrentHashMap {
	return &ConcurrentHashMap{
		size: atomic.Uint64{},
	}
}

// Get retrieves a skip list map for a given key
func (chm *ConcurrentHashMap) Get(key []byte) (*ConcurrentSkipListMap[interface{}], bool) {
	val, ok := chm.data.Load(string(key))
	if !ok {
		return nil, false
	}
	return val.(*ConcurrentSkipListMap[interface{}]), true
}

// ComputeIfAbsent gets or creates a skip list map for a given key
func (chm *ConcurrentHashMap) ComputeIfAbsent(key []byte, mappingFunc func() *ConcurrentSkipListMap[interface{}]) *ConcurrentSkipListMap[interface{}] {
	insertKey := make([]byte, len(key))
	copy(insertKey, key)

	if val, ok := chm.data.Load(string(insertKey)); ok {
		return val.(*ConcurrentSkipListMap[interface{}])
	}

	newVal := mappingFunc()
	actual, loaded := chm.data.LoadOrStore(string(insertKey), newVal)
	if loaded {
		return actual.(*ConcurrentSkipListMap[interface{}])
	}
	return newVal
}

// Put adds or updates a skip list map
func (chm *ConcurrentHashMap) Put(key []byte, value *ConcurrentSkipListMap[interface{}]) {
	chm.size.Add(1)
	chm.data.Store(key, value)
}

// Delete removes a key from the map
func (chm *ConcurrentHashMap) Delete(key []byte) {
	_, loaded := chm.data.LoadAndDelete(string(key))
	if loaded {
		chm.size.Add(^uint64(0))
	}
}

// Size returns the number of keys in the map
func (chm *ConcurrentHashMap) Size() uint64 {
	count := chm.size.Load()
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
func (mvcc *MVCC) Read(key []byte, version int64) (string, bool) {
	skipList, ok := mvcc.versions.Get(key)
	if !ok {
		return "", false
	}

	// Find the greatest version <= requested version
	_, value, found := skipList.FloorEntry([]byte(fmt.Sprintf("%d", version)))
	return value.(string), found
}

// Pool для skip list map
var skipListMapPool = sync.Pool{
	New: func() interface{} {
		return NewConcurrentSkipListMap[interface{}]()
	},
}

// Write stores a value for a key at a specific version
func (mvcc *MVCC) Write(key []byte, value string, version int64) {
	skipList := mvcc.versions.ComputeIfAbsent(key, func() *ConcurrentSkipListMap[interface{}] {
		sl := skipListMapPool.Get().(*ConcurrentSkipListMap[interface{}])
		// Reset skip list state
		sl.level = 1
		// Очищаем head forward pointers
		for i := 0; i < maxLevel; i++ {
			sl.head.forward[i] = nil
		}
		return sl
	})

	var versionBytes [8]byte
	binary.BigEndian.PutUint64(versionBytes[:], uint64(version))
	skipList.Put(versionBytes[:], value)
}

// Delete removes a key from the MVCC store
func (mvcc *MVCC) Delete(key []byte) {
	mvcc.versions.Delete(key)
}

// GetVersionCount returns the number of versions stored for a key
func (mvcc *MVCC) GetVersionCount(key []byte) int {
	skipList, ok := mvcc.versions.Get(key)
	if !ok {
		return 0
	}
	return skipList.Size()
}
