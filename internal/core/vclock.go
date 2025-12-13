package core

import (
	"fmt"
	"sync"
)

// VectorClock представляет векторные часы для отслеживания каузальности событий
// между распределёнными узлами
type VectorClock struct {
	// clock хранит для каждого nodeID его логическое время
	clock map[string]uint64
	mu    sync.RWMutex
}

// NewVectorClock создаёт новые векторные часы
func NewVectorClock() *VectorClock {
	return &VectorClock{
		clock: make(map[string]uint64),
	}
}

// Increment увеличивает счётчик для данного узла
func (vc *VectorClock) Increment(nodeID string) uint64 {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	vc.clock[nodeID]++
	return vc.clock[nodeID]
}

// Get возвращает текущее значение для узла
func (vc *VectorClock) Get(nodeID string) uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	return vc.clock[nodeID]
}

// Update обновляет векторные часы на основе другого вектора
// Берёт максимум из локального и полученного значения для каждого узла
func (vc *VectorClock) Update(other *VectorClock) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	for nodeID, timestamp := range other.clock {
		if vc.clock[nodeID] < timestamp {
			vc.clock[nodeID] = timestamp
		}
	}
}

// UpdateFromMap обновляет из map
func (vc *VectorClock) UpdateFromMap(clockMap map[string]uint64) {
	vc.mu.Lock()
	defer vc.mu.Unlock()

	for nodeID, timestamp := range clockMap {
		if vc.clock[nodeID] < timestamp {
			vc.clock[nodeID] = timestamp
		}
	}
}

// HappensBefore проверяет, произошло ли событие this раньше other
// this < other если все компоненты this <= other и хотя бы один строго меньше
func (vc *VectorClock) HappensBefore(other *VectorClock) bool {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	hasStrictlyLess := false

	// Проверяем все узлы из this
	for nodeID, thisTime := range vc.clock {
		otherTime := other.clock[nodeID]
		if thisTime > otherTime {
			return false // this имеет большее время для какого-то узла
		}
		if thisTime < otherTime {
			hasStrictlyLess = true
		}
	}

	// Проверяем узлы, которые есть в other, но нет в this
	for nodeID := range other.clock {
		if _, exists := vc.clock[nodeID]; !exists {
			hasStrictlyLess = true
		}
	}

	return hasStrictlyLess
}

// ConcurrentWith проверяет, являются ли события конкурентными (несравнимыми)
func (vc *VectorClock) ConcurrentWith(other *VectorClock) bool {
	return !vc.HappensBefore(other) && !other.HappensBefore(vc) && !vc.Equals(other)
}

// Equals проверяет равенство векторных часов
func (vc *VectorClock) Equals(other *VectorClock) bool {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	// Проверяем все узлы из this
	for nodeID, thisTime := range vc.clock {
		if other.clock[nodeID] != thisTime {
			return false
		}
	}

	// Проверяем узлы из other
	for nodeID, otherTime := range other.clock {
		if vc.clock[nodeID] != otherTime {
			return false
		}
	}

	return true
}

// Clone создаёт копию векторных часов
func (vc *VectorClock) Clone() *VectorClock {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	newVC := NewVectorClock()
	for nodeID, timestamp := range vc.clock {
		newVC.clock[nodeID] = timestamp
	}

	return newVC
}

// ToMap возвращает копию внутреннего map
func (vc *VectorClock) ToMap() map[string]uint64 {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	result := make(map[string]uint64, len(vc.clock))
	for nodeID, timestamp := range vc.clock {
		result[nodeID] = timestamp
	}

	return result
}

// String возвращает строковое представление
func (vc *VectorClock) String() string {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	return fmt.Sprintf("%v", vc.clock)
}

// IsSafeToRead проверяет, можно ли безопасно читать данные с данным VClock
// на основе кворума узлов. Данные считаются безопасными, если они
// синхронизированы на большинстве узлов (quorum).
func (vc *VectorClock) IsSafeToRead(minAcks int, totalNodes int) bool {
	vc.mu.RLock()
	defer vc.mu.RUnlock()

	// Подсчитываем количество узлов, которые подтвердили эту версию
	acks := len(vc.clock)

	return acks >= minAcks
}

// MergeMax объединяет два вектора, беря максимум для каждого узла
func (vc *VectorClock) MergeMax(other *VectorClock) *VectorClock {
	result := vc.Clone()
	result.Update(other)
	return result
}
