package core

import (
	"bytes"
	"math/rand"
	"sync"
)

const (
	maxLevel    = 16
	probability = 0.5
)

// // Pool для переиспользования update массивов
// var updatePool = sync.Pool{
// 	New: func() interface{} {
// 		return make([]*SkipListNode[interface{}], maxLevel)
// 	},
// }

// // Pool для переиспользования forward массивов разных размеров
// var forwardPools = [maxLevel + 1]sync.Pool{}

// func init() {
// 	for i := 1; i <= maxLevel; i++ {
// 		size := i
// 		forwardPools[i].New = func() interface{} {
// 			return make([]*SkipListNode[interface{}], size)
// 		}
// 	}
// }

// func getForwardSlice[V any](size int) []*SkipListNode[V] {
// 	if size <= 0 || size > maxLevel {
// 		return make([]*SkipListNode[V], size)
// 	}
// 	return forwardPools[size].Get().([]*SkipListNode[V])
// }

// func putForwardSlice[V any](slice []*SkipListNode[V]) {
// 	size := len(slice)
// 	if size <= 0 || size > maxLevel {
// 		return
// 	}
// 	// Очищаем slice перед возвратом
// 	for i := range slice {
// 		slice[i] = nil
// 	}
// 	forwardPools[size].Put(slice)
// }

// SkipListNode represents a node in the skip list
type SkipListNode[V any] struct {
	// key         int64
	key         []byte
	value       V
	forward     []*SkipListNode[V]
	prev0       *SkipListNode[V]
	versionData *MVCCVersion // Для хранения полной версии с VectorClock
}

// Для обратной совместимости
type skipListNode[V any] = SkipListNode[V]

// ConcurrentSkipListMap is a thread-safe skip list map
type ConcurrentSkipListMap[V any] struct {
	head     *SkipListNode[V]
	level    int
	mu       sync.RWMutex
	rng      *rand.Rand
	rngMutex sync.Mutex
}

// NewConcurrentSkipListMap creates a new concurrent skip list map
func NewConcurrentSkipListMap[V any]() *ConcurrentSkipListMap[V] {
	head := &SkipListNode[V]{
		key:     []byte{},
		forward: make([]*SkipListNode[V], maxLevel),
	}
	return &ConcurrentSkipListMap[V]{
		head:  head,
		level: 1,
		rng:   rand.New(rand.NewSource(42)),
	}
}

// randomLevel generates a random level for a new node
func (sl *ConcurrentSkipListMap[V]) randomLevel() int {
	sl.rngMutex.Lock()
	defer sl.rngMutex.Unlock()

	level := 1
	for level < maxLevel && sl.rng.Float64() < probability {
		level++
	}
	return level
}

// Put inserts or updates a key-value pair
func (sl *ConcurrentSkipListMap[V]) Put(key []byte, value V) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	insertKey := make([]byte, len(key))
	copy(insertKey, key)

	var update [maxLevel]*SkipListNode[V]
	current := sl.head

	// Find position to insert
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, insertKey) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	// Update existing node
	if current != nil && bytes.Equal(current.key, insertKey) {
		current.value = value
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

	forwardSlice := make([]*SkipListNode[V], newLevel)
	newNode := &SkipListNode[V]{
		key:     insertKey,
		value:   value,
		forward: forwardSlice,
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	// Для обработного прохода по уровню 0 устанавливаем prev0 ссылки
	prev := update[0]          // узел перед newNode на уровне 0
	next := newNode.forward[0] // узел после newNode на уровне 0

	newNode.prev0 = prev
	if next != nil {
		next.prev0 = newNode
	}

	// Очищаем update перед возвратом в пул
	for i := 0; i < maxLevel; i++ {
		update[i] = nil
	}
}

// Get retrieves a value by key
func (sl *ConcurrentSkipListMap[V]) Get(key []byte) (V, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil && bytes.Equal(current.key, key) {
		return current.value, true
	}
	var zero V
	return zero, false
}

// FloorEntry returns the greatest key-value pair less than or equal to the given key
func (sl *ConcurrentSkipListMap[V]) FloorEntry(key []byte) ([]byte, V, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.head
	var lastValid *SkipListNode[V]

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, key) <= 0 {
			current = current.forward[i]
			if bytes.Compare(current.key, key) <= 0 {
				lastValid = current
			}
		}
	}

	if lastValid != nil {
		return lastValid.key, lastValid.value, true
	}
	var zero V
	return nil, zero, false
}

// CeilingEntry returns the smallest key-value pair greater than or equal to the given key
func (sl *ConcurrentSkipListMap[V]) CeilingEntry(key []byte) ([]byte, V, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	if current != nil {
		return current.key, current.value, true
	}
	var zero V
	return nil, zero, false
}

func (sl *ConcurrentSkipListMap[V]) Delete(key []byte) bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	var update [maxLevel]*SkipListNode[V]
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]
	if current == nil || !bytes.Equal(current.key, key) {
		return false
	}

	// перешиваем prev0 (обратную цепочку) на уровне 0
	prev := update[0]
	next := current.forward[0]
	if next != nil {
		next.prev0 = prev
	}

	// разрываем forward ссылки
	for i := 0; i < sl.level; i++ {
		if update[i].forward[i] != current {
			break
		}
		update[i].forward[i] = current.forward[i]
	}

	// уменьшаем уровень если нужно
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}

	return true
}

// Size returns the number of elements in the skip list
func (sl *ConcurrentSkipListMap[V]) Size() int {
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

func (sl *ConcurrentSkipListMap[V]) ComputeIfAbsent(key []byte, mappingFunc func() V) V {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Check if exists (do search under same lock)
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && bytes.Compare(cur.forward[i].key, key) < 0 {
			cur = cur.forward[i]
		}
	}
	cur = cur.forward[0]
	if cur != nil && bytes.Equal(cur.key, key) {
		return cur.value
	}

	// Create and insert once
	v := mappingFunc()

	k := make([]byte, len(key))
	copy(k, key)

	sl.putUnsafeLocked(k, v)
	return v
}

// putUnsafeLocked inserts or updates a key-value pair.
// IMPORTANT: caller must hold sl.mu.Lock().
func (sl *ConcurrentSkipListMap[V]) putUnsafeLocked(key []byte, value V) {
	// Find position to insert (fill update)
	var update [maxLevel]*SkipListNode[V]
	current := sl.head

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && bytes.Compare(current.forward[i].key, key) < 0 {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	// Update existing node
	if current != nil && bytes.Equal(current.key, key) {
		current.value = value
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

	forwardSlice := make([]*SkipListNode[V], newLevel)
	newNode := &SkipListNode[V]{
		key:     key, // key must already be an immutable copy
		value:   value,
		forward: forwardSlice,
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	// prev0 wiring (level 0)
	prev := update[0]
	next := newNode.forward[0]
	newNode.prev0 = prev
	if next != nil {
		next.prev0 = newNode
	}
}

// findFirstGE returns first node with key >= target (or nil). Caller must hold sl.mu (RLock or Lock).
func (sl *ConcurrentSkipListMap[V]) findFirstGE(target []byte) *SkipListNode[V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && bytes.Compare(cur.forward[i].key, target) < 0 {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// findLastLE returns last node with key <= target (or nil). Caller must hold sl.mu (RLock or Lock).
func (sl *ConcurrentSkipListMap[V]) findLastLE(target []byte) *SkipListNode[V] {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && bytes.Compare(cur.forward[i].key, target) <= 0 {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

type IteratorType int

const (
	IteratorTypeEq IteratorType = iota
	IteratorTypeReq
	IteratorTypeAll
	IteratorTypeGE
	IteratorTypeGT
	IteratorTypeLE
	IteratorTypeLT
)

type SkipListIterator[V any] struct {
	sl  *ConcurrentSkipListMap[V]
	typ IteratorType

	// current element
	cur   *SkipListNode[V]
	ended bool

	// direction
	forward bool
}

func (sl *ConcurrentSkipListMap[V]) NewIterator(seekKey []byte, typ IteratorType) *SkipListIterator[V] {
	it := &SkipListIterator[V]{sl: sl, typ: typ}

	sl.mu.RLock()
	defer sl.mu.RUnlock()

	switch typ {
	case IteratorTypeAll:
		it.forward = true
		it.cur = sl.findFirstGE(seekKey)

	case IteratorTypeEq:
		it.forward = true
		n := sl.findFirstGE(seekKey)
		if n != nil && bytes.Equal(n.key, seekKey) {
			it.cur = n
		} else {
			it.ended = true
		}

	case IteratorTypeGE:
		it.forward = true
		it.cur = sl.findFirstGE(seekKey)

	case IteratorTypeGT:
		it.forward = true
		n := sl.findFirstGE(seekKey)
		if n != nil && bytes.Equal(n.key, seekKey) {
			n = n.forward[0]
		}
		it.cur = n

	case IteratorTypeLE:
		it.forward = false
		it.cur = sl.findLastLE(seekKey)

	case IteratorTypeLT:
		it.forward = false
		n := sl.findLastLE(seekKey)
		if n != nil && bytes.Equal(n.key, seekKey) {
			n = n.prev0
			if n == sl.head {
				n = nil
			}
		}
		it.cur = n

	default:
		it.ended = true
	}

	if it.cur == nil {
		it.ended = true
	}
	return it
}

// Valid reports whether iterator currently points to an element.
func (it *SkipListIterator[V]) Valid() bool {
	return !it.ended && it.cur != nil
}

// Next advances iterator by one. Returns true if now valid, false if ended.
func (it *SkipListIterator[V]) Next() bool {
	if it.ended {
		return false
	}

	// EQ is single-element
	if it.typ == IteratorTypeEq {
		it.ended = true
		return false
	}

	// Step under short RLock (avoid blocking writers for long scans)
	it.sl.mu.RLock()
	if it.forward {
		it.cur = it.cur.forward[0]
	} else {
		it.cur = it.cur.prev0
		if it.cur == it.sl.head {
			it.cur = nil
		}
	}
	it.sl.mu.RUnlock()

	if it.cur == nil {
		it.ended = true
		return false
	}
	return true
}

func (it *SkipListIterator[V]) Key() []byte {
	if !it.Valid() {
		return nil
	}
	return it.cur.key
}

func (it *SkipListIterator[V]) Value() V {
	if !it.Valid() {
		var zero V
		return zero
	}
	return it.cur.value
}

// Version returns versionData if you store MVCCVersion in nodes.
func (it *SkipListIterator[V]) Version() *MVCCVersion {
	if !it.Valid() {
		return nil
	}
	return it.cur.versionData
}

func (it *SkipListIterator[V]) Close() {
	// no-op for now
}
