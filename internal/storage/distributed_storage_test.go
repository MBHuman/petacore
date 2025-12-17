package storage_test

import (
	"context"
	"errors"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/storage"
	"sync"
	"testing"
	"time"
)

// TestDistributedStorageBasicReadWrite тестирует базовые операции чтения/записи
func TestDistributedStorageBasicReadWrite(t *testing.T) {
	kvStore := SetupKVStore(t, "test_basic")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start distributed storage: %v", err)
	}
	defer ds.Stop()

	// Ждем синхронизации
	time.Sleep(100 * time.Millisecond)

	// Записываем данные
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("key1", "value1")
		tx.Write("key2", "value2")
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Даем время на синхронизацию
	time.Sleep(100 * time.Millisecond)

	// Читаем данные
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		value1, ok := tx.Read("key1")
		if !ok {
			t.Error("key1 not found")
		}
		if value1 != "value1" {
			t.Errorf("Expected value1, got %s", value1)
		}

		value2, ok := tx.Read("key2")
		if !ok {
			t.Error("key2 not found")
		}
		if value2 != "value2" {
			t.Errorf("Expected value2, got %s", value2)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
}

// TestDistributedStorageReadCommitted тестирует уровень изоляции Read Committed
func TestDistributedStorageReadCommitted(t *testing.T) {
	kvStore := SetupKVStore(t, "test_read_committed")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Записываем начальное значение
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("key", "value1")
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to write initial value: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Обновляем значение
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("key", "value2")
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to update value: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Читаем - должны получить последнее значение
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		value, ok := tx.Read("key")
		if !ok {
			t.Fatal("Key not found")
		}
		if value != "value2" {
			t.Errorf("Expected value2, got %s", value)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
}

// TestDistributedStorageSnapshotIsolation тестирует изоляцию Snapshot
func TestDistributedStorageSnapshotIsolation(t *testing.T) {
	kvStore := SetupKVStore(t, "test_snapshot")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.SnapshotIsolation)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Записываем начальное значение
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("key", "value1")
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Записываем несколько версий
	for i := 2; i <= 5; i++ {
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write("key", "value"+string(rune('0'+i)))
			return nil
		})
		if err != nil {
			t.Fatalf("Failed to write value%d: %v", i, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Проверяем, что получаем актуальную версию
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		value, ok := tx.Read("key")
		if !ok {
			t.Fatal("Key not found")
		}
		// Должны получить последнюю версию
		if value != "value5" {
			t.Errorf("Expected value5, got %s", value)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read: %v", err)
	}
}

// TestDistributedStorageLocalWrites тестирует локальные записи внутри транзакции
func TestDistributedStorageLocalWrites(t *testing.T) {
	kvStore := SetupKVStore(t, "test_local_writes")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		// Записываем локально
		tx.Write("key1", "value1")

		// Читаем - должны получить локальное значение
		value, ok := tx.Read("key1")
		if !ok {
			t.Error("key1 not found in local writes")
		}
		if value != "value1" {
			t.Errorf("Expected value1, got %s", value)
		}

		// Перезаписываем
		tx.Write("key1", "value2")

		// Читаем снова
		value, ok = tx.Read("key1")
		if !ok {
			t.Error("key1 not found after rewrite")
		}
		if value != "value2" {
			t.Errorf("Expected value2, got %s", value)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}
}

// TestDistributedStorageEmptyTransaction тестирует пустую транзакцию
func TestDistributedStorageEmptyTransaction(t *testing.T) {
	kvStore := SetupKVStore(t, "test_empty_tx")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Пустая транзакция не должна вызывать ошибок
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		return nil
	})
	if err != nil {
		t.Fatalf("Empty transaction failed: %v", err)
	}
}

// TestDistributedStorageMultipleKeys тестирует работу с множественными ключами
func TestDistributedStorageMultipleKeys(t *testing.T) {
	kvStore := SetupKVStore(t, "test_multiple_keys")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Записываем много ключей
	const numKeys = 100
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		for i := 0; i < numKeys; i++ {
			key := "key" + string(rune('0'+i%10))
			value := "value" + string(rune('0'+i%10))
			tx.Write(key, value)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to write multiple keys: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	// Читаем и проверяем
	err = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		for i := 0; i < 10; i++ {
			key := "key" + string(rune('0'+i))
			value, ok := tx.Read(key)
			if !ok {
				t.Errorf("Key %s not found", key)
			}
			expectedValue := "value" + string(rune('0'+i))
			if value != expectedValue {
				t.Errorf("Key %s: expected %s, got %s", key, expectedValue, value)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read multiple keys: %v", err)
	}
}

// TestDistributedStorageConcurrentTransactions тестирует конкурентные транзакции
func TestDistributedStorageConcurrentTransactions(t *testing.T) {
	kvStore := SetupKVStore(t, "test_concurrent")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	const numGoroutines = 10
	const numOpsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOpsPerGoroutine; j++ {
				err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					key := "key" + string(rune('0'+id))
					value := "value" + string(rune('0'+j))
					tx.Write(key, value)
					return nil
				})
				if err != nil {
					t.Errorf("Goroutine %d: transaction failed: %v", id, err)
				}
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(200 * time.Millisecond)

	// Проверяем, что все ключи записались
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		for i := 0; i < numGoroutines; i++ {
			key := "key" + string(rune('0'+i))
			_, ok := tx.Read(key)
			if !ok {
				t.Errorf("Key %s not found after concurrent writes", key)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to verify concurrent writes: %v", err)
	}
}

// TestDistributedStorageSyncStatus тестирует статус синхронизации
func TestDistributedStorageSyncStatus(t *testing.T) {
	kvStore := SetupKVStore(t, "test_sync_status")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)

	// До старта должен быть статус Syncing
	status := ds.GetSyncStatus()
	if status != distributed.SyncStatusSyncing {
		t.Errorf("Expected SyncStatusSyncing before start, got %v", status)
	}

	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	// Ждем синхронизации
	time.Sleep(100 * time.Millisecond)

	// После старта должен быть Synced
	if !ds.IsSynced() {
		t.Error("Expected storage to be synced after start")
	}

	status = ds.GetSyncStatus()
	if status != distributed.SyncStatusSynced {
		t.Errorf("Expected SyncStatusSynced after start, got %v", status)
	}
}

// TestDistributedStorageTransactionError тестирует обработку ошибок в транзакции
func TestDistributedStorageTransactionError(t *testing.T) {
	kvStore := SetupKVStore(t, "test_tx_error")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Транзакция с ошибкой не должна применять изменения
	testError := errors.New("test error")
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("key", "value")
		return testError
	})

	if err != testError {
		t.Errorf("Expected test error, got %v", err)
	}
}

// TestDistributedStorageReadNonExistentKey тестирует чтение несуществующего ключа
func TestDistributedStorageReadNonExistentKey(t *testing.T) {
	kvStore := SetupKVStore(t, "test_read_nonexistent")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		_, ok := tx.Read("nonexistent")
		if ok {
			t.Error("Expected nonexistent key to return false")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Transaction failed: %v", err)
	}
}

// TestDistributedStorageInitialData тестирует загрузку начальных данных из ETCD
func TestDistributedStorageInitialData(t *testing.T) {
	kvStore := SetupKVStore(t, "test_initial_data")
	if kvStore == nil {
		t.Skip("KVStore not available")
	}
	defer kvStore.Close()

	// Предварительно заполняем kvStore
	ctx := context.Background()
	kvStore.Put(ctx, "existing1", "value1", 1)
	kvStore.Put(ctx, "existing2", "value2", 2)

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		t.Fatalf("Failed to start: %v", err)
	}
	defer ds.Stop()

	time.Sleep(100 * time.Millisecond)

	// Читаем предзагруженные данные
	err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		value1, ok := tx.Read("existing1")
		if !ok {
			t.Error("existing1 not found")
		}
		if value1 != "value1" {
			t.Errorf("Expected value1, got %s", value1)
		}

		value2, ok := tx.Read("existing2")
		if !ok {
			t.Error("existing2 not found")
		}
		if value2 != "value2" {
			t.Errorf("Expected value2, got %s", value2)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to read initial data: %v", err)
	}
}
