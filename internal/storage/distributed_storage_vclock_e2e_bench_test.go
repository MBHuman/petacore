package storage

import (
	"context"
	"fmt"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"sync"
	"testing"
	"time"
)

// E2E бенчмарки с реальным ETCD кластером
// Запуск: go test -bench=E2E -benchmem -benchtime=10s ./internal/storage/

var (
	etcdEndpoints = []string{"localhost:2379", "localhost:2479", "localhost:2579"}
	testPrefix    = "petacore-bench-e2e"
)

// setupE2EStorage создает реальное хранилище с ETCD
func setupE2EStorage(b *testing.B, nodeID string, totalNodes int, minAcks int) (*DistributedStorageVClock, func()) {
	b.Helper()

	kvStore, err := distributed.NewETCDStore(etcdEndpoints, testPrefix)
	if err != nil {
		b.Fatalf("Failed to connect to ETCD: %v", err)
	}

	storage := NewDistributedStorageVClock(kvStore, nodeID, totalNodes, core.ReadCommitted, minAcks)

	if err := storage.Start(); err != nil {
		kvStore.Close()
		b.Fatalf("Failed to start storage: %v", err)
	}

	// Ждем начальной синхронизации
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if storage.IsSynced() {
				goto synced
			}
		case <-timeout:
			b.Logf("Warning: initial sync timeout")
			goto synced
		}
	}

synced:
	cleanup := func() {
		storage.Stop()
		kvStore.Close()

		// Очистка тестовых данных
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cleanupStore, _ := distributed.NewETCDStore(etcdEndpoints, testPrefix)
		if cleanupStore != nil {
			// Удаляем все ключи с префиксом
			cleanupStore.DeleteAll(ctx, testPrefix)
			cleanupStore.Close()
		}
	}

	return storage, cleanup
}

// BenchmarkE2E_WriteQuorum тестирует запись с quorum
func BenchmarkE2E_WriteQuorum(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0) // 0 = quorum
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:write:quorum:%d", i)
			value := fmt.Sprintf("value_%d", i)

			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(key, value)
				return nil
			})

			if err != nil {
				b.Errorf("Write failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_WriteAllNodes тестирует запись с требованием всех узлов
func BenchmarkE2E_WriteAllNodes(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, -1) // -1 = all nodes
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:write:all:%d", i)
			value := fmt.Sprintf("value_%d", i)

			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(key, value)
				return nil
			})

			if err != nil {
				b.Errorf("Write failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_WriteWeakConsistency тестирует запись со слабой консистентностью
func BenchmarkE2E_WriteWeakConsistency(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 1)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:write:weak:%d", i)
			value := fmt.Sprintf("value_%d", i)

			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(key, value)
				return nil
			})

			if err != nil {
				b.Errorf("Write failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_ReadWithQuorum тестирует чтение с quorum проверкой
func BenchmarkE2E_ReadWithQuorum(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	// Подготовка: записываем тестовые данные
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bench:read:quorum:%d", i)
		value := fmt.Sprintf("value_%d", i)
		storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
			tx.Write(key, value)
			return nil
		})
	}

	time.Sleep(1 * time.Second) // Ждем синхронизации

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:read:quorum:%d", i%100)

			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				_, ok := tx.Read(key)
				if !ok {
					return fmt.Errorf("key not found: %s", key)
				}
				return nil
			})

			if err != nil {
				b.Errorf("Read failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_MixedWorkload тестирует смешанную нагрузку (70% чтение, 30% запись)
func BenchmarkE2E_MixedWorkload(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	// Подготовка данных
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench:mixed:%d", i)
		value := fmt.Sprintf("value_%d", i)
		storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
			tx.Write(key, value)
			return nil
		})
	}

	time.Sleep(2 * time.Second)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:mixed:%d", i%1000)

			// 70% чтений, 30% записей
			if i%10 < 7 {
				// Чтение
				storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
					_, _ = tx.Read(key)
					return nil
				})
			} else {
				// Запись
				value := fmt.Sprintf("new_value_%d", i)
				storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
					tx.Write(key, value)
					return nil
				})
			}
			i++
		}
	})
}

// BenchmarkE2E_DistributedWrites тестирует распределенные записи на 3 узла
func BenchmarkE2E_DistributedWrites(b *testing.B) {
	// Создаем 3 узла
	storage1, cleanup1 := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup1()

	storage2, cleanup2 := setupE2EStorage(b, "bench-node-2", 3, 0)
	defer cleanup2()

	storage3, cleanup3 := setupE2EStorage(b, "bench-node-3", 3, 0)
	defer cleanup3()

	storages := []*DistributedStorageVClock{storage1, storage2, storage3}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Распределяем записи по узлам
			nodeIdx := i % 3
			key := fmt.Sprintf("bench:distributed:%d", i)
			value := fmt.Sprintf("value_%d", i)

			err := storages[nodeIdx].RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(key, value)
				return nil
			})

			if err != nil {
				b.Errorf("Write failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_VClockSync тестирует скорость синхронизации Vector Clock
func BenchmarkE2E_VClockSync(b *testing.B) {
	storage1, cleanup1 := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup1()

	storage2, cleanup2 := setupE2EStorage(b, "bench-node-2", 3, 0)
	defer cleanup2()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench:vclock:sync:%d", i)
		value := fmt.Sprintf("value_%d", i)

		// Записываем на node-1
		err := storage1.RunTransaction(func(tx *DistributedTransactionVClock) error {
			tx.Write(key, value)
			return nil
		})
		if err != nil {
			b.Errorf("Write failed: %v", err)
		}

		// Ждем небольшое время для синхронизации
		time.Sleep(10 * time.Millisecond)

		// Читаем с node-2
		err = storage2.RunTransaction(func(tx *DistributedTransactionVClock) error {
			readValue, ok := tx.Read(key)
			if !ok {
				return fmt.Errorf("key not synced yet")
			}
			if readValue != value {
				return fmt.Errorf("value mismatch")
			}
			return nil
		})
		if err != nil {
			b.Errorf("Read failed: %v", err)
		}
	}
}

// BenchmarkE2E_ConcurrentWrites тестирует параллельные записи
func BenchmarkE2E_ConcurrentWrites(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	concurrency := []int{1, 10, 50, 100}

	for _, c := range concurrency {
		b.Run(fmt.Sprintf("Concurrency_%d", c), func(b *testing.B) {
			b.SetParallelism(c)
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("bench:concurrent:%d:%d", c, i)
					value := fmt.Sprintf("value_%d", i)

					err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
						tx.Write(key, value)
						return nil
					})

					if err != nil {
						b.Errorf("Write failed: %v", err)
					}
					i++
				}
			})
		})
	}
}

// BenchmarkE2E_BatchWrites тестирует батчевую запись
func BenchmarkE2E_BatchWrites(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	batchSizes := []int{1, 10, 50, 100}

	for _, size := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
					for j := 0; j < size; j++ {
						key := fmt.Sprintf("bench:batch:%d:%d:%d", size, i, j)
						value := fmt.Sprintf("value_%d_%d", i, j)
						tx.Write(key, value)
					}
					return nil
				})

				if err != nil {
					b.Errorf("Batch write failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkE2E_ConsistencyCheck тестирует проверку консистентности
func BenchmarkE2E_ConsistencyCheck(b *testing.B) {
	storage1, cleanup1 := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup1()

	storage2, cleanup2 := setupE2EStorage(b, "bench-node-2", 3, 0)
	defer cleanup2()

	storage3, cleanup3 := setupE2EStorage(b, "bench-node-3", 3, 0)
	defer cleanup3()

	storages := []*DistributedStorageVClock{storage1, storage2, storage3}

	// Записываем данные на первый узел
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("bench:consistency:%d", i)
		value := fmt.Sprintf("value_%d", i)
		storage1.RunTransaction(func(tx *DistributedTransactionVClock) error {
			tx.Write(key, value)
			return nil
		})
	}

	time.Sleep(2 * time.Second) // Ждем синхронизации

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench:consistency:%d", i%100)

		// Читаем со всех узлов
		values := make([]string, 3)
		for idx, storage := range storages {
			storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				value, ok := tx.Read(key)
				if ok {
					values[idx] = value
				}
				return nil
			})
		}

		// Проверяем консистентность
		if values[0] != "" && values[1] != "" && values[2] != "" {
			if values[0] != values[1] || values[0] != values[2] {
				b.Errorf("Inconsistent values for key %s: %v", key, values)
			}
		}
	}
}

// BenchmarkE2E_HighContentionWrite тестирует запись в один ключ с высокой конкуренцией
func BenchmarkE2E_HighContentionWrite(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	hotKey := "bench:hotkey"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			value := fmt.Sprintf("value_%d", i)

			err := storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(hotKey, value)
				return nil
			})

			if err != nil {
				b.Errorf("Write failed: %v", err)
			}
			i++
		}
	})
}

// BenchmarkE2E_Latency измеряет латентность операций
func BenchmarkE2E_Latency(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	// Записываем данные
	key := "bench:latency:test"
	storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
		tx.Write(key, "test_value")
		return nil
	})

	time.Sleep(1 * time.Second)

	var totalLatency time.Duration
	var mu sync.Mutex

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			start := time.Now()

			storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				_, _ = tx.Read(key)
				return nil
			})

			latency := time.Since(start)
			mu.Lock()
			totalLatency += latency
			mu.Unlock()
		}
	})

	avgLatency := totalLatency / time.Duration(b.N)
	b.ReportMetric(float64(avgLatency.Microseconds()), "µs/op")
}

// BenchmarkE2E_Throughput измеряет общую пропускную способность
func BenchmarkE2E_Throughput(b *testing.B) {
	storage, cleanup := setupE2EStorage(b, "bench-node-1", 3, 0)
	defer cleanup()

	b.ResetTimer()
	start := time.Now()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("bench:throughput:%d", i)
			value := fmt.Sprintf("value_%d", i)

			storage.RunTransaction(func(tx *DistributedTransactionVClock) error {
				tx.Write(key, value)
				return nil
			})
			i++
		}
	})

	elapsed := time.Since(start)
	throughput := float64(b.N) / elapsed.Seconds()
	b.ReportMetric(throughput, "ops/sec")
}
