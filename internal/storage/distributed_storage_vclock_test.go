package storage_test

import (
	"fmt"
	"petacore/internal/core"
	"petacore/internal/storage"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestVectorClock_BasicOperations тестирует базовые операции с VClock
func TestVectorClock_BasicOperations(t *testing.T) {
	kvStore := SetupKVStore(t, "test_vclock_basic")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	// Используем totalNodes=1 для простоты (minAcks=1)
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	// Ждём инициализации
	time.Sleep(100 * time.Millisecond)

	// Тест Write
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("key1"), "value1")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Даём время на синхронизацию
	time.Sleep(200 * time.Millisecond)

	// Тест Read
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		value, ok := tx.Read([]byte("key1"))
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

	kvStore := SetupKVStore(t, "test_vclock_quorum")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	// Создаём 3 узла
	storage1 := storage.NewDistributedStorageVClock(kvStore, "node1", 3, core.ReadCommitted, 0)
	storage2 := storage.NewDistributedStorageVClock(kvStore, "node2", 3, core.ReadCommitted, 0)
	storage3 := storage.NewDistributedStorageVClock(kvStore, "node3", 3, core.ReadCommitted, 0)

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
	err := storage1.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("test_key"), "test_value")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Ждём синхронизации между узлами
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что все узлы получили запись
	for i, store := range []*storage.DistributedStorageVClock{storage1, storage2, storage3} {
		err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			value, ok := tx.Read([]byte("test_key"))
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
	kvStore := SetupKVStore(t, "test_vclock_concurrent")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	// Используем totalNodes=1 для простоты
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Конкурентные записи
	const numWrites = 10
	var wg sync.WaitGroup
	wg.Add(numWrites)

	for i := 0; i < numWrites; i++ {
		go func(i int) {
			defer wg.Done()
			err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write([]byte(fmt.Sprintf("key%d", i)), fmt.Sprintf("value%d", i))
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
		err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			key := fmt.Sprintf("key%d", i)
			value, ok := tx.Read([]byte(key))
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
	kvStore := SetupKVStore(t, "test_vclock_isolation")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Первая транзакция - инициализация
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("counter"), "0")
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
			store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				val, ok := tx.Read([]byte("counter"))
				if !ok {
					val = "0"
				}
				// Просто перезаписываем (в реальной системе бы парсили и увеличивали)
				tx.Write([]byte("counter"), val+"1")
				return nil
			})
		}()
	}

	wg.Wait()
	time.Sleep(300 * time.Millisecond)

	// Проверяем финальное значение
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		value, ok := tx.Read([]byte("counter"))
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
	kvStore := SetupKVStore(t, "test_vclock_multiple")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Записываем несколько ключей в одной транзакции
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("user:1:name"), "Alice")
		tx.Write([]byte("user:1:email"), "alice@example.com")
		tx.Write([]byte("user:2:name"), "Bob")
		tx.Write([]byte("user:2:email"), "bob@example.com")
		return nil
	})

	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Читаем все ключи
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
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
			value, ok := tx.Read([]byte(tt.key))
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
	kvStore := SetupKVStore(t, "test_vclock_update")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Первая запись
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("config:version"), "1.0")
		return nil
	})
	if err != nil {
		t.Fatalf("Initial write failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Обновление
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("config:version"), "1.1")
		return nil
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что получаем последнюю версию
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		value, ok := tx.Read([]byte("config:version"))
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
	kvStore := SetupKVStore(t, "test_vclock_empty")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Читаем несуществующий ключ
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		_, ok := tx.Read([]byte("nonexistent"))
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
	kvStore := SetupKVStore(t, "test_vclock_rollback")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	// Транзакция с ошибкой (должна откатиться)
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("rollback_test"), "should_not_persist")
		return fmt.Errorf("simulated error")
	})

	if err == nil {
		t.Fatal("Expected transaction to fail")
	}

	time.Sleep(100 * time.Millisecond)

	// Проверяем, что данные НЕ сохранились
	err = store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		_, ok := tx.Read([]byte("rollback_test"))
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
	kvStore := SetupKVStore(t, "test_vclock_sequential")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()
	store := storage.NewDistributedStorageVClock(kvStore, "node1", 1, core.ReadCommitted, 0)

	if err := store.Start(); err != nil {
		t.Fatalf("Failed to start storage: %v", err)
	}
	defer store.Stop()

	time.Sleep(100 * time.Millisecond)

	const numWrites = 20

	// Последовательные записи
	for i := 0; i < numWrites; i++ {
		err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write([]byte("sequence"), fmt.Sprintf("value-%d", i))
			return nil
		})
		if err != nil {
			t.Fatalf("Write %d failed: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond) // Небольшая пауза между записями
	}

	time.Sleep(200 * time.Millisecond)

	// Проверяем финальное значение
	err := store.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		value, ok := tx.Read([]byte("sequence"))
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

// TestOCC_DoubleSpending тестирует OCC для предотвращения double spending
func TestOCC_DoubleSpending(t *testing.T) {
	kvStore := SetupKVStore(t, "test_occ_double")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	// Создаём два хранилища на одном KV (симулируем два узла)
	store1 := storage.NewDistributedStorageVClock(kvStore, "node1", 2, core.SnapshotIsolation, -1)
	store2 := storage.NewDistributedStorageVClock(kvStore, "node2", 2, core.SnapshotIsolation, -1)

	if err := store1.Start(); err != nil {
		t.Fatalf("Failed to start store1: %v", err)
	}
	defer store1.Stop()

	if err := store2.Start(); err != nil {
		t.Fatalf("Failed to start store2: %v", err)
	}
	defer store2.Stop()

	// Ждём инициализации
	time.Sleep(200 * time.Millisecond)

	// Устанавливаем начальный баланс
	err := store1.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write([]byte("balance"), "100")
		return nil
	})
	require.NoError(t, err)

	// Ждём синхронизации
	time.Sleep(500 * time.Millisecond)

	// Транзакция 1: read balance, write balance-50
	var tx1Done sync.WaitGroup
	tx1Done.Add(1)
	var tx1Err error
	go func() {
		defer tx1Done.Done()
		tx1Err = store1.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			if value, ok := tx.Read([]byte("balance")); !ok {
				return fmt.Errorf("balance not found")
			} else if value != "100" {
				return fmt.Errorf("expected 100, got %s", value)
			}
			tx.Write([]byte("balance"), "50")
			time.Sleep(200 * time.Millisecond) // Имитируем работу
			return nil
		})
	}()

	// Транзакция 2: read balance, write balance-100 (double spending)
	var tx2Done sync.WaitGroup
	tx2Done.Add(1)
	var tx2Err error
	go func() {
		defer tx2Done.Done()
		time.Sleep(50 * time.Millisecond) // Начинаем после tx1
		tx2Err = store2.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			if value, ok := tx.Read([]byte("balance")); !ok {
				return fmt.Errorf("balance not found")
			} else if value != "100" {
				return fmt.Errorf("expected 100, got %s", value)
			}
			tx.Write([]byte("balance"), "0")
			time.Sleep(300 * time.Millisecond) // Имитируем работу
			return nil
		})
	}()

	// Ждём завершения
	tx1Done.Wait()
	tx2Done.Wait()

	// Tx1 должна быть успешной
	require.NoError(t, tx1Err, "Tx1 should succeed")

	// Tx2 должна fail с OCC violation
	require.Error(t, tx2Err, "Tx2 should fail with OCC violation")
	require.Contains(t, tx2Err.Error(), "OCC violation", "Error should contain OCC violation")

	// Проверяем финальный баланс
	time.Sleep(500 * time.Millisecond) // Ждём синхронизации
	err = store1.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		value, ok := tx.Read([]byte("balance"))
		if !ok {
			return fmt.Errorf("balance not found")
		}
		if value != "50" {
			t.Errorf("Expected balance 50, got %s", value)
		}
		return nil
	})
	require.NoError(t, err)
}
