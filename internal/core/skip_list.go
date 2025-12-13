package core

import (
	"math/rand"
	"sync"
)

const (
	maxLevel    = 16
	probability = 0.5
)

// Pool для переиспользования update массивов
var updatePool = sync.Pool{
	New: func() interface{} {
		return make([]*SkipListNode, maxLevel)
	},
}

// Pool для переиспользования forward массивов разных размеров
var forwardPools = [maxLevel + 1]sync.Pool{}

func init() {
	for i := 1; i <= maxLevel; i++ {
		size := i
		forwardPools[i].New = func() interface{} {
			return make([]*SkipListNode, size)
		}
	}
}

func getForwardSlice(size int) []*SkipListNode {
	if size <= 0 || size > maxLevel {
		return make([]*SkipListNode, size)
	}
	return forwardPools[size].Get().([]*SkipListNode)
}

func putForwardSlice(slice []*SkipListNode) {
	size := len(slice)
	if size <= 0 || size > maxLevel {
		return
	}
	// Очищаем slice перед возвратом
	for i := range slice {
		slice[i] = nil
	}
	forwardPools[size].Put(slice)
}

// SkipListNode represents a node in the skip list
type SkipListNode struct {
	key     int64
	value   string
	forward []*SkipListNode
}

// ConcurrentSkipListMap is a thread-safe skip list map
type ConcurrentSkipListMap struct {
	head     *SkipListNode
	level    int
	mu       sync.RWMutex
	rng      *rand.Rand
	rngMutex sync.Mutex
}

// NewConcurrentSkipListMap creates a new concurrent skip list map
func NewConcurrentSkipListMap() *ConcurrentSkipListMap {
	head := &SkipListNode{
		key:     -1,
		forward: make([]*SkipListNode, maxLevel),
	}
	return &ConcurrentSkipListMap{
		head:  head,
		level: 1,
		rng:   rand.New(rand.NewSource(42)),
	}
}

// randomLevel generates a random level for a new node
func (sl *ConcurrentSkipListMap) randomLevel() int {
	sl.rngMutex.Lock()
	defer sl.rngMutex.Unlock()

	level := 1
	for level < maxLevel && sl.rng.Float64() < probability {
		level++
	}
	return level
}

// Put inserts or updates a key-value pair
func (sl *ConcurrentSkipListMap) Put(key int64, value string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	update := updatePool.Get().([]*SkipListNode)
	defer updatePool.Put(update)

	current := sl.head

	// Find position to insert
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	// Update existing node
	if current != nil && current.key == key {
		current.value = value
		// Очищаем update перед возвратом в пул
		for i := 0; i < maxLevel; i++ {
			update[i] = nil
		}
		return
	}

	// Insert new node
	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.head
		}
		sl.level = newLevel
	}

	forwardSlice := getForwardSlice(newLevel)
	newNode := &SkipListNode{
		key:     key,
		value:   value,
		forward: forwardSlice,
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	// Очищаем update перед возвратом в пул
	for i := 0; i < maxLevel; i++ {
		update[i] = nil
	}
}

// Get retrieves a value by key
func (sl *ConcurrentSkipListMap) Get(key int64) (string, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil && current.key == key {
		return current.value, true
	}
	return "", false
}

// FloorEntry returns the greatest key-value pair less than or equal to the given key
func (sl *ConcurrentSkipListMap) FloorEntry(key int64) (int64, string, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.head
	var lastValid *SkipListNode

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key <= key {
			current = current.forward[i]
			if current.key <= key {
				lastValid = current
			}
		}
	}

	if lastValid != nil {
		return lastValid.key, lastValid.value, true
	}
	return 0, "", false
}

// Delete removes a key-value pair
func (sl *ConcurrentSkipListMap) Delete(key int64) bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	update := updatePool.Get().([]*SkipListNode)
	defer func() {
		// Очищаем update перед возвратом в пул
		for i := 0; i < maxLevel; i++ {
			update[i] = nil
		}
		updatePool.Put(update)
	}()

	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]
	if current == nil || current.key != key {
		return false
	}

	// Возвращаем forward slice в пул перед удалением
	putForwardSlice(current.forward)

	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] != current {
			break
		}
		update[i].forward[i] = current.forward[i]
	}

	// Update level
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}

	return true
}

// Size returns the number of elements in the skip list
func (sl *ConcurrentSkipListMap) Size() int {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	count := 0
	current := sl.head.forward[0]
	for current != nil {
		count++
		current = current.forward[0]
	}
	return count
}
