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
	mu       sync.RWMutex
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

// ReadWithVClock читает значение с проверкой Vector Clock
// Возвращает последнюю безопасную версию (с подтверждённым quorum)
func (m *MVCCWithVClock) ReadWithVClock(key string, minAcks int, totalNodes int) (string, *VectorClock, bool) {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return "", nil, false
	}

	// Получаем все версии и ищем последнюю безопасную
	versions := skipList.GetAllVersions()
	if len(versions) == 0 {
		return "", nil, false
	}

	// Идём от самой новой к самой старой
	for i := len(versions) - 1; i >= 0; i-- {
		version := versions[i]

		// Проверяем, безопасна ли эта версия для чтения
		if version.VectorClock.IsSafeToRead(minAcks, totalNodes) {
			return version.Value, version.VectorClock, true
		}
	}

	// Нет безопасных версий
	return "", nil, false
}

// ReadLatest читает последнюю версию без проверки quorum
// (используется для внутренних целей)
func (m *MVCCWithVClock) ReadLatest(key string) (string, *VectorClock, bool) {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return "", nil, false
	}

	versions := skipList.GetAllVersions()
	if len(versions) == 0 {
		return "", nil, false
	}

	latest := versions[len(versions)-1]
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

// GetAllVersions возвращает все версии в порядке возрастания timestamp
func (sl *ConcurrentSkipListMap) GetAllVersions() []*MVCCVersion {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	var versions []*MVCCVersion
	current := sl.head.forward[0]

	for current != nil {
		if current.versionData != nil {
			versions = append(versions, current.versionData)
		}
		current = current.forward[0]
	}

	return versions
}
