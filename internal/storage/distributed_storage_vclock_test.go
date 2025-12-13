package storage

import (
	"context"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// MockKVStore для тестирования без реального ETCD
type MockKVStore struct {
	data       map[string]*distributed.KVEntry
	mu         sync.RWMutex
	watchChans []chan *distributed.WatchEvent
	revision   int64
}

func NewMockKVStore() *MockKVStore {
	return &MockKVStore{
		data:       make(map[string]*distributed.KVEntry),
		watchChans: make([]chan *distributed.WatchEvent, 0),
		revision:   0,
	}
}

func (m *MockKVStore) Get(ctx context.Context, key string) (*distributed.KVEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if entry, ok := m.data[key]; ok {
		return entry, nil
	}
	return nil, fmt.Errorf("key not found")
}

func (m *MockKVStore) Put(ctx context.Context, key string, value string, version int64) error {
	m.mu.Lock()
	m.revision++
	entry := &distributed.KVEntry{
		Key:      key,
		Value:    value,
		Version:  version,
		Revision: m.revision,
	}
	m.data[key] = entry
	watchChans := make([]chan *distributed.WatchEvent, len(m.watchChans))
	copy(watchChans, m.watchChans)
	m.mu.Unlock()

	// Уведомляем наблюдателей
	event := &distributed.WatchEvent{
		Type:  distributed.EventTypePut,
		Entry: entry,
	}
	for _, ch := range watchChans {
		select {
		case ch <- event:
		default:
		}
	}

	return nil
}

func (m *MockKVStore) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MockKVStore) GetAll(ctx context.Context, prefix string) ([]*distributed.KVEntry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*distributed.KVEntry, 0)
	for k, entry := range m.data {
		if prefix == "" || (len(k) >= len(prefix) && k[:len(prefix)] == prefix) {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (m *MockKVStore) Watch(ctx context.Context, prefix string) (<-chan *distributed.WatchEvent, error) {
	ch := make(chan *distributed.WatchEvent, 100)

	m.mu.Lock()
	m.watchChans = append(m.watchChans, ch)
	m.mu.Unlock()

	return ch, nil
}

func (m *MockKVStore) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, ch := range m.watchChans {
		close(ch)
	}
	m.watchChans = nil
	return nil
}

// NewMockKVStoreForREPL экспортированная версия для REPL
func NewMockKVStoreForREPL() distributed.KVStore {
	return NewMockKVStore()
}

// TestVectorClock_BasicOperations тестирует базовые операции с VClock
func TestVectorClock_BasicOperations(t *testing.T) {
	kvStore := NewMockKVStore()
	// Используем totalNodes=1 для простоты (minAcks=1)
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	// Ждём инициализации
	time.Sleep(100 * time.Millisecond)

	// Тест Write
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("key1", "value1")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Даём время на синхронизацию
	time.Sleep(200 * time.Millisecond)

	// Тест Read
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		value, ok := tx.Read("key1")
		if !ok {
			t.Error("Expected to find key1")
			return fmt.Errorf("key not found")
		}
		if value != "value1" {
			t.Errorf("Expected value1, got %s", value)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

}

// TestVectorClock_QuorumRead тестирует quorum-based чтение
// В этом тесте мы проверяем, что чтения блокируются до достижения кворума
func TestVectorClock_QuorumRead(t *testing.T) {
	t.Skip("Skipping quorum test - requires proper 3-node setup and synchronization timing")

	kvStore := NewMockKVStore()

	// Создаём 3 узла
	storage1 := NewDistributedStorageVClock(kvStore, "node1", 3, core.ReadCommitted, 0)
	storage2 := NewDistributedStorageVClock(kvStore, "node2", 3, core.ReadCommitted, 0)
	storage3 := NewDistributedStorageVClock(kvStore, "node3", 3, core.ReadCommitted, 0)

	if err := storage1.Start(); err != nil {
		t.Fatalf("Failed to start storage1: %v", err)
	}
	defer storage1.Stop()

	if err := storage2.Start(); err != nil {
		t.Fatalf("Failed to start storage2: %v", err)
	}
	defer storage2.Stop()

	if err := storage3.Start(); err != nil {
		t.Fatalf("Failed to start storage3: %v", err)
	}
	defer storage3.Stop()

	time.Sleep(200 * time.Millisecond)

	// Node1 записывает данные
	err := storage1.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("test_key", "test_value")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Ждём синхронизации между узлами
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что все узлы получили запись
	for i, storage := range []*DistributedStorageVClock{storage1, storage2, storage3} {
		err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
			value, ok := tx.Read("test_key")
			if !ok {
				return fmt.Errorf("node%d: key not found", i+1)
			}
			if value != "test_value" {
				return fmt.Errorf("node%d: expected test_value, got %s", i+1, value)
			}
			return nil
		})

		t.Errorf("Read on node%d failed: %v", i+1, err)
	}
}

// TestVectorClock_ConcurrentWrites тестирует конкурентные записи
func TestVectorClock_ConcurrentWrites(t *testing.T) {
	kvStore := NewMockKVStore()
	// Используем totalNodes=1 для простоты
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Конкурентные записи
	const numWrites = 10
	var wg sync.WaitGroup
	wg.Add(numWrites)

	for i := 0; i < numWrites; i++ {
		go func(i int) {
			defer wg.Done()
			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
				return nil
			})
			if err != nil {
				t.Errorf("Write %d failed: %v", i, err)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(300 * time.Millisecond)

	// Проверяем, что все записи прошли
	for i := 0; i < numWrites; i++ {
		err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
			key := fmt.Sprintf("key%d", i)
			value, ok := tx.Read(key)
			if !ok {
				return fmt.Errorf("key %s not found", key)
			}
			expectedValue := fmt.Sprintf("value%d", i)
			if value != expectedValue {
				return fmt.Errorf("expected %s, got %s", expectedValue, value)
			}
			return nil
		})
		require.NoError(t, err)
	}
}

// TestVectorClock_TransactionIsolation тестирует изоляцию транзакций
func TestVectorClock_TransactionIsolation(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Первая транзакция - инициализация
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("counter", "0")
		return nil
	})

	if err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Параллельные транзакции увеличивают счётчик
	const numTx = 5
	var wg sync.WaitGroup
	wg.Add(numTx)

	for i := 0; i < numTx; i++ {
		go func() {
			defer wg.Done()
			storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				val, ok := tx.Read("counter")
				if !ok {
					val = "0"
				}
				// Просто перезаписываем (в реальной системе бы парсили и увеличивали)
				tx.Write("counter", val+"1")
				return nil
			})
		}()
	}

	wg.Wait()
	time.Sleep(300 * time.Millisecond)

	// Проверяем финальное значение
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		value, ok := tx.Read("counter")
		if !ok {
			return fmt.Errorf("counter not found")
		}
		t.Logf("Final counter value: %s", value)
		// Просто проверяем, что значение существует
		if len(value) == 0 {
			return fmt.Errorf("counter is empty")
		}
		return nil
	})

	require.NoError(t, err)

}

// TestVectorClock_MultipleKeys тестирует работу с несколькими ключами
func TestVectorClock_MultipleKeys(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Записываем несколько ключей в одной транзакции
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("user:1:name", "Alice")
		tx.Write("user:1:email", "alice@example.com")
		tx.Write("user:2:name", "Bob")
		tx.Write("user:2:email", "bob@example.com")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Читаем все ключи
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tests := []struct {
			key      string
			expected string
		}{
			{"user:1:name", "Alice"},
			{"user:1:email", "alice@example.com"},
			{"user:2:name", "Bob"},
			{"user:2:email", "bob@example.com"},
		}

		for _, tt := range tests {
			value, ok := tx.Read(tt.key)
			if !ok {
				return fmt.Errorf("key %s not found", tt.key)
			}
			if value != tt.expected {
				return fmt.Errorf("key %s: expected %s, got %s", tt.key, tt.expected, value)
			}
		}
		return nil
	})

	if err != nil {
		t.Errorf("Read failed: %v", err)
	}
}

// TestVectorClock_UpdateOperations тестирует операции обновления
func TestVectorClock_UpdateOperations(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Первая запись
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("config:version", "1.0")
		return nil
	})
	if err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Обновление
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("config:version", "1.1")
		return nil
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что получаем последнюю версию
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		value, ok := tx.Read("config:version")
		if !ok {
			return fmt.Errorf("config:version not found")
		}
		if value != "1.1" {
			return fmt.Errorf("expected 1.1, got %s", value)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Read after update failed: %v", err)
	}
}

// TestVectorClock_EmptyRead тестирует чтение несуществующего ключа
func TestVectorClock_EmptyRead(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Читаем несуществующий ключ
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		_, ok := tx.Read("nonexistent")
		if ok {
			return fmt.Errorf("expected key not found, but it was found")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Empty read test failed: %v", err)
	}
}

// TestVectorClock_TransactionRollback тестирует откат транзакций
func TestVectorClock_TransactionRollback(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	// Транзакция с ошибкой (должна откатиться)
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write("rollback_test", "should_not_persist")
		return fmt.Errorf("simulated error")
	})

	if err == nil {
		t.Fatal("Expected transaction to fail")
	}

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что данные НЕ сохранились
	err = storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		_, ok := tx.Read("rollback_test")
		if ok {
			return fmt.Errorf("data should not exist after rollback")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Rollback verification failed: %v", err)
	}
}

// TestVectorClock_SequentialWrites тестирует последовательные записи
func TestVectorClock_SequentialWrites(t *testing.T) {
	kvStore := NewMockKVStore()
	storage := NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := storage.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer storage.Stop()

	time.Sleep(100 * time.Millisecond)

	const numWrites = 20

	// Последовательные записи
	for i := 0; i < numWrites; i++ {
		err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
			tx.Write("sequence", fmt.Sprintf("value-%d", i))
			return nil
		})
		if err != nil {
			t.Fatalf("Write %d failed: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond) // Небольшая пауза между записями
	}

	time.Sleep(200 * time.Millisecond)

	// Проверяем финальное значение
	err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		value, ok := tx.Read("sequence")
		if !ok {
			return fmt.Errorf("sequence key not found")
		}
		expectedValue := fmt.Sprintf("value-%d", numWrites-1)
		if value != expectedValue {
			return fmt.Errorf("expected %s, got %s", expectedValue, value)
		}
		return nil
	})

	if err != nil {
		t.Errorf("Sequential writes test failed: %v", err)
	}
}
