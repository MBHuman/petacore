package core

import (
	"fmt"
	"petacore/internal/utils"
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
	versions *ConcurrentSkipListMap[*ConcurrentSkipListMap[*MVCCVersion]]
}

// NewMVCCWithVClock создаёт новый MVCC с Vector Clock
func NewMVCCWithVClock() *MVCCWithVClock {
	return &MVCCWithVClock{
		versions: NewConcurrentSkipListMap[*ConcurrentSkipListMap[*MVCCVersion]](),
	}
}

// WriteWithVClock записывает значение с Vector Clock
func (m *MVCCWithVClock) WriteWithVClock(key []byte, value string, timestamp uint64, vclock *VectorClock) {
	skipList := m.versions.ComputeIfAbsent(key, func() *ConcurrentSkipListMap[*MVCCVersion] {
		return NewConcurrentSkipListMap[*MVCCVersion]()
	})

	version := &MVCCVersion{
		Value:       value,
		VectorClock: vclock.Clone(),
		Timestamp:   timestamp,
	}

	// Сохраняем MVCCVersion под timestamp
	skipList.Put([]byte(fmt.Sprintf("%d", timestamp)), version)
}

// ReadWithSnapshot читает значение с snapshot isolation
// snapshotVClock != nil для SI, snapshotTimestamp для RC
func (m *MVCCWithVClock) ReadWithSnapshot(key []byte, snapshotVClock *VectorClock, snapshotTimestamp uint64, minAcks int, totalNodes int, currentNodeID string) (string, *VectorClock, bool) {
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
func (m *MVCCWithVClock) ReadLatest(key []byte) (string, *VectorClock, bool) {
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
func (m *MVCCWithVClock) GetVectorClock(key []byte) (*VectorClock, bool) {
	_, vclock, ok := m.ReadLatest(key)
	return vclock, ok
}

// UpdateVectorClock обновляет Vector Clock для существующей версии
// Это нужно когда другой узел подтверждает запись
func (m *MVCCWithVClock) UpdateVectorClock(key []byte, timestamp uint64, nodeID string) bool {
	skipList, ok := m.versions.Get(key)
	if !ok {
		return false
	}

	version, ok := skipList.Get(fmt.Appendf(nil, "%d", timestamp))
	if !ok {
		return false
	}

	// Обновляем Vector Clock
	version.VectorClock.Increment(nodeID)
	return true
}

// GetIterator возвращает итератор для версий ключа
func (m *MVCCWithVClock) GetIterator(key []byte, it IteratorType) *SkipListIterator[*ConcurrentSkipListMap[*MVCCVersion]] {
	iterator := m.versions.NewIterator(key, it)
	return iterator
}

// VersionIterator итератор для версий в ConcurrentSkipListMap
type VersionIterator[V any] struct {
	current *skipListNode[V]
	mu      *sync.RWMutex // Для управления блокировкой
}

// NewVersionIterator создаёт итератор для версий (от старых к новым)
func (sl *ConcurrentSkipListMap[V]) NewVersionIterator() *VersionIterator[V] {
	sl.mu.RLock() // Блокируем на чтение для всего итератора
	return &VersionIterator[V]{
		current: sl.head,
		mu:      &sl.mu,
	}
}

// Next переходит к следующей версии, возвращает true если есть следующая
func (it *VersionIterator[V]) Next() bool {
	it.current = it.current.forward[0]
	return it.current != nil
}

// Value возвращает текущую версию
func (it *VersionIterator[V]) Value() V {
	if it.current == nil {
		var zero V
		return zero
	}
	return it.current.value
}

// Close разблокирует мьютекс
func (it *VersionIterator[V]) Close() {
	it.mu.RUnlock()
}

// ScanWithSnapshot сканирует ключи с префиксом и возвращает последние безопасные версии
func (m *MVCCWithVClock) ScanWithSnapshot(prefix []byte, it IteratorType, snapshotVClock *VectorClock, snapshotTimestamp uint64, minAcks int, totalNodes int, currentNodeID string, limit int) map[string]string {
	iterator := m.versions.NewIterator(prefix, it)
	if iterator == nil {
		return make(map[string]string)
	}
	defer iterator.Close()

	result := make(map[string]string)
	count := 0
	for iterator.Next() {
		// log.Printf("DEBUG: Scanning key: %s", iterator.Key())
		if limit > 0 && count >= limit {
			break
		}
		key := iterator.Key()
		if !utils.HasPrefix(key, prefix) {
			break
		}
		// Для каждого ключа применяем логику ReadWithSnapshot
		skipList := iterator.Value()
		versionIterator := skipList.NewVersionIterator()
		defer versionIterator.Close()

		var latestSafe *MVCCVersion
		for versionIterator.Next() {
			version := versionIterator.Value()
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
			result[string(key)] = latestSafe.Value
			count++
		}
	}
	return result
}
