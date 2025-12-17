package core

import (
	"sync"
)

// MVCCVersion представляет версию данных с Vector Clock
type MVCCVersion struct {
	Value       string
	VectorClock *VectorClock
	Timestamp   uint64 // Для совместимости с HLC
}

// MVCCWithVClock расширенный MVCC с поддержкой Vector Clock
type MVCCWithVClock struct {
	// версии данных: key -> (timestamp -> MVCCVersion)
	versions *ConcurrentHashMap
}

// NewMVCCWithVClock создаёт новый MVCC с Vector Clock
func NewMVCCWithVClock() *MVCCWithVClock {
	return &MVCCWithVClock{
		versions: NewConcurrentHashMap(),
	}
}

// WriteWithVClock записывает значение с Vector Clock
func (m *MVCCWithVClock) WriteWithVClock(key string, value string, timestamp uint64, vclock *VectorClock) {
	skipList := m.versions.ComputeIfAbsent(key, func() *ConcurrentSkipListMap {
		return NewConcurrentSkipListMap()
	})

	version := &MVCCVersion{
		Value:       value,
		VectorClock: vclock.Clone(),
		Timestamp:   timestamp,
	}

	// Сохраняем MVCCVersion под timestamp
	skipList.PutVersion(int64(timestamp), version)
}

// ReadWithSnapshot читает значение с snapshot isolation
// snapshotVClock != nil для SI, snapshotTimestamp для RC
func (m *MVCCWithVClock) ReadWithSnapshot(key string, snapshotVClock *VectorClock, snapshotTimestamp uint64, minAcks int, totalNodes int, currentNodeID string) (string, *VectorClock, bool) {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return "", nil, false
	}

	// Используем итератор для поиска последней безопасной версии
	// которая произошла до snapshotVClock
	iterator := skipList.NewVersionIterator()
	defer iterator.Close()

	var latestSafe *MVCCVersion

	// Итерируем от старых к новым, ищем максимальную безопасную версию
	for iterator.Next() {
		version := iterator.Value()
		safe := version.VectorClock.IsSafeToRead(minAcks, totalNodes, currentNodeID)
		if snapshotVClock != nil {
			// SI: check VClock
			safe = safe && !version.VectorClock.HappensAfter(snapshotVClock)
		} else {
			// RC: check timestamp
			safe = safe && version.Timestamp <= snapshotTimestamp
		}
		if safe {
			latestSafe = version
		}
	}

	if latestSafe != nil {
		return latestSafe.Value, latestSafe.VectorClock, true
	}

	// Нет безопасных версий в snapshot
	return "", nil, false
}

// ReadLatest читает последнюю версию без проверки quorum
// (используется для внутренних целей)
func (m *MVCCWithVClock) ReadLatest(key string) (string, *VectorClock, bool) {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return "", nil, false
	}

	iterator := skipList.NewVersionIterator()
	defer iterator.Close()

	var latest *MVCCVersion

	// Итерируем до конца, чтобы найти последнюю версию
	for iterator.Next() {
		latest = iterator.Value()
	}

	if latest == nil {
		return "", nil, false
	}

	return latest.Value, latest.VectorClock, true
}

// GetVectorClock возвращает Vector Clock для последней версии ключа
func (m *MVCCWithVClock) GetVectorClock(key string) (*VectorClock, bool) {
	_, vclock, ok := m.ReadLatest(key)
	return vclock, ok
}

// UpdateVectorClock обновляет Vector Clock для существующей версии
// Это нужно когда другой узел подтверждает запись
func (m *MVCCWithVClock) UpdateVectorClock(key string, timestamp uint64, nodeID string) bool {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return false
	}

	version := skipList.GetVersion(int64(timestamp))
	if version == nil {
		return false
	}

	// Обновляем Vector Clock
	version.VectorClock.Increment(nodeID)
	return true
}

// ConcurrentSkipListMap расширенный для хранения MVCCVersion
func (sl *ConcurrentSkipListMap) PutVersion(key int64, version *MVCCVersion) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	update := make([]*skipListNode, maxLevel)
	current := sl.head

	// Find position to insert
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].key < key {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	// Update existing key
	if current != nil && current.key == key {
		current.versionData = version
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

	newNode := &skipListNode{
		key:         key,
		value:       version.Value,
		versionData: version,
		forward:     make([]*skipListNode, maxLevel),
	}

	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}

// GetVersion получает MVCCVersion по timestamp
func (sl *ConcurrentSkipListMap) GetVersion(key int64) *MVCCVersion {
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
		return current.versionData
	}

	return nil
}

// VersionIterator итератор для версий в ConcurrentSkipListMap
type VersionIterator struct {
	current *skipListNode
	mu      *sync.RWMutex // Для управления блокировкой
}

// NewVersionIterator создаёт итератор для версий (от старых к новым)
func (sl *ConcurrentSkipListMap) NewVersionIterator() *VersionIterator {
	sl.mu.RLock() // Блокируем на чтение для всего итератора
	return &VersionIterator{
		current: sl.head.forward[0],
		mu:      &sl.mu,
	}
}

// Next переходит к следующей версии, возвращает true если есть следующая
func (it *VersionIterator) Next() bool {
	if it.current == nil {
		return false
	}
	it.current = it.current.forward[0]
	return it.current != nil && it.current.versionData != nil
}

// Value возвращает текущую версию
func (it *VersionIterator) Value() *MVCCVersion {
	if it.current == nil {
		return nil
	}
	return it.current.versionData
}

// Close разблокирует мьютекс
func (it *VersionIterator) Close() {
	it.mu.RUnlock()
}
