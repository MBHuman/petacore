package storage_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/storage"
	"strings"
	"sync"
	"testing"
	"time"
)

// SetupKVStore создает тестовый KVStore на основе env переменных
// DB_TYPE: etcd (default), postgres
// ETCD_ENDPOINTS: для etcd, default "localhost:2379"
// PG_CONN_STRING: для postgres, default "postgres://postgres:password@localhost/petacore_test?sslmode=disable"
func SetupKVStore(tb testing.TB, namespace string) distributed.KVStore {
	tb.Helper()

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "etcd"
	}

	var kvStore distributed.KVStore
	var err error

	switch strings.ToLower(dbType) {
	case "etcd":
		endpointsStr := os.Getenv("ETCD_ENDPOINTS")
		if endpointsStr == "" {
			endpointsStr = "localhost:2379"
		}
		endpoints := strings.Split(endpointsStr, ",")
		for i, ep := range endpoints {
			endpoints[i] = strings.TrimSpace(ep)
		}
		kvStore, err = distributed.NewETCDStore(endpoints, namespace)
		if err != nil {
			tb.Skipf("ETCD не доступен (%v), пропускаем тест: %v", endpoints, err)
			return nil
		}

	case "postgres":
		connStr := os.Getenv("PG_CONN_STRING")
		if connStr == "" {
			connStr = "postgres://postgres:password@localhost/petacore_test?sslmode=disable"
		}
		kvStore, err = distributed.NewPGStore(connStr, namespace)
		if err != nil {
			tb.Skipf("PostgreSQL не доступен (%s), пропускаем тест: %v", connStr, err)
			return nil
		}

	default:
		tb.Fatalf("Неизвестный DB_TYPE: %s", dbType)
		return nil
	}

	return kvStore
}

// SetupDistributedStorage создает тестовое распределенное хранилище на основе env переменных
// DB_TYPE: etcd (default), postgres
// ETCD_ENDPOINTS: для etcd, default "localhost:2379"
// PG_CONN_STRING: для postgres, default "postgres://postgres:password@localhost/petacore_test?sslmode=disable"
func SetupDistributedStorage(tb testing.TB) *storage.DistributedStorage {
	tb.Helper()

	namespace := fmt.Sprintf("%s_%d", tb.Name(), time.Now().UnixNano())

	kvStore := SetupKVStore(tb, namespace)
	if kvStore == nil {
		return nil
	}

	ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
	if err := ds.Start(); err != nil {
		tb.Fatalf("Не удалось запустить синхронизацию: %v", err)
	}

	// Ждем синхронизации
	timeout := time.After(10 * time.Second)
	for !ds.IsSynced() {
		select {
		case <-timeout:
			tb.Fatal("Таймаут ожидания синхронизации")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	tb.Cleanup(func() {
		ds.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if kvStoreT, ok := kvStore.(distributed.KVStoreT); ok {
			_ = kvStoreT.DeleteAll(ctx, "")
		}
		kvStore.Close()
	})

	return ds
}

// SetupDistributedStorageVClock создает тестовое распределенное хранилище VClock на основе env переменных
func SetupDistributedStorageVClock(tb testing.TB, nodeID string, totalNodes int, isolation core.IsolationLevel, minAcks int) *storage.DistributedStorageVClock {
	tb.Helper()

	namespace := fmt.Sprintf("%s_%d", tb.Name(), time.Now().UnixNano())

	kvStore := SetupKVStore(tb, namespace)
	if kvStore == nil {
		return nil
	}

	ds := storage.NewDistributedStorageVClock(kvStore, nodeID, totalNodes, isolation, minAcks)
	if err := ds.Start(); err != nil {
		tb.Fatalf("Не удалось запустить синхронизацию: %v", err)
	}

	// Ждем синхронизации
	timeout := time.After(10 * time.Second)
	for !ds.IsSynced() {
		select {
		case <-timeout:
			tb.Fatal("Таймаут ожидания синхронизации")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	tb.Cleanup(func() {
		ds.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if kvStoreT, ok := kvStore.(distributed.KVStoreT); ok {
			_ = kvStoreT.DeleteAll(ctx, "")
		}
		kvStore.Close()
	})

	return ds
}

// setupDistributedStorage создает тестовое распределенное хранилище (legacy, используй SetupDistributedStorage)

// BenchmarkDistributedWrite измеряет производительность записей
func BenchmarkDistributedWrite(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "value")
			return nil
		})
		if err != nil {
			b.Fatalf("Ошибка записи: %v", err)
		}
	}
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "writes/sec")
}

// BenchmarkDistributedRead измеряет производительность чтений
func BenchmarkDistributedRead(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			_, _ = tx.Read(key)
			return nil
		})
		if err != nil {
			b.Fatalf("Ошибка чтения: %v", err)
		}
	}
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "reads/sec")
}

// BenchmarkDistributedReadWrite измеряет производительность смешанных операций
func BenchmarkDistributedReadWrite(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Предварительно записываем данные
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "initial_value")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readKey := fmt.Sprintf("key_%d", i%100)
		writeKey := fmt.Sprintf("key_%d", (i+1)%100)
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			_, _ = tx.Read(readKey)
			tx.Write(writeKey, "updated_value")
			return nil
		})
		if err != nil {
			b.Fatalf("Ошибка транзакции: %v", err)
		}
	}
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tx/sec")
}

// BenchmarkDistributedTransaction измеряет производительность сложных транзакций
func BenchmarkDistributedTransaction(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Инициализация
	_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("account:alice", "1000")
		tx.Write("account:bob", "500")
		return nil
	})

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			alice, _ := tx.Read("account:alice")
			bob, _ := tx.Read("account:bob")

			// Симулируем перевод средств
			tx.Write("account:alice", alice+"-100")
			tx.Write("account:bob", bob+"+100")
			return nil
		})
		if err != nil {
			b.Fatalf("Ошибка транзакции: %v", err)
		}
	}
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tx/sec")
}

// BenchmarkDistributedConcurrentWrites измеряет конкурентные записи
func BenchmarkDistributedConcurrentWrites(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i)
			err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
				tx.Write(key, "value")
				return nil
			})
			if err != nil {
				log.Printf("Ошибка записи: %v", err)
			}
			i++
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "writes/sec")
}

// BenchmarkDistributedConcurrentReads измеряет конкурентные чтения
func BenchmarkDistributedConcurrentReads(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i%1000)
			err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
				_, _ = tx.Read(key)
				return nil
			})
			if err != nil {
				log.Printf("Ошибка чтения: %v", err)
			}
			i++
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "reads/sec")
}

// BenchmarkDistributedHotKey измеряет работу с горячим ключом
func BenchmarkDistributedHotKey(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Инициализируем счетчик
	_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
		tx.Write("counter", "0")
		return nil
	})

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
				value, _ := tx.Read("counter")
				tx.Write("counter", value+"1")
				return nil
			})
			if err != nil {
				log.Printf("Ошибка: %v", err)
			}
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "updates/sec")
}

// BenchmarkDistributedIsolationLevels сравнивает уровни изоляции
func BenchmarkDistributedIsolationLevels(b *testing.B) {
	levels := []struct {
		name  string
		level core.IsolationLevel
	}{
		{"ReadCommitted", core.ReadCommitted},
		{"SnapshotIsolation", core.SnapshotIsolation},
	}

	for _, lvl := range levels {
		b.Run(lvl.name, func(b *testing.B) {
			namespace := fmt.Sprintf("bench_iso_%s_%d", lvl.name, time.Now().UnixNano())
			kvStore := SetupKVStore(b, namespace)
			if kvStore == nil {
				return
			}
			defer kvStore.Close()

			ds := storage.NewDistributedStorage(kvStore, lvl.level)
			if err := ds.Start(); err != nil {
				b.Fatalf("Не удалось запустить: %v", err)
			}
			defer ds.Stop()

			timeout := time.After(10 * time.Second)
			for !ds.IsSynced() {
				select {
				case <-timeout:
					b.Fatal("Таймаут синхронизации")
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}

			// Предварительно записываем данные
			for i := 0; i < 100; i++ {
				key := fmt.Sprintf("key_%d", i)
				_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					tx.Write(key, "value")
					return nil
				})
			}

			time.Sleep(500 * time.Millisecond)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key_%d", i%100)
				err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					_, _ = tx.Read(key)
					tx.Write(key, "updated")
					return nil
				})
				if err != nil {
					b.Fatalf("Ошибка: %v", err)
				}
			}
			b.StopTimer()

			b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "tx/sec")
		})
	}
}

// BenchmarkDistributedMultiNode симулирует несколько узлов
func BenchmarkDistributedMultiNode(b *testing.B) {
	nodeCount := 3
	nodes := make([]*storage.DistributedStorage, nodeCount)
	namespace := fmt.Sprintf("bench_multi_%d", time.Now().UnixNano())

	// Создаем несколько узлов
	for i := 0; i < nodeCount; i++ {
		kvStore := SetupKVStore(b, namespace)
		if kvStore == nil {
			return
		}

		ds := storage.NewDistributedStorage(kvStore, core.ReadCommitted)
		if err := ds.Start(); err != nil {
			b.Fatalf("Не удалось запустить узел %d: %v", i, err)
		}

		nodes[i] = ds

		b.Cleanup(func() {
			ds.Stop()
			if i == 0 {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if kvStoreT, ok := kvStore.(distributed.KVStoreT); ok {
					_ = kvStoreT.DeleteAll(ctx, "")
				}
				cancel()
			}
			kvStore.Close()
		})
	}

	// Ждем синхронизации всех узлов
	for i, node := range nodes {
		timeout := time.After(10 * time.Second)
		for !node.IsSynced() {
			select {
			case <-timeout:
				b.Fatalf("Таймаут синхронизации узла %d", i)
			default:
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	b.ResetTimer()

	// Каждый узел выполняет свою часть работы
	var wg sync.WaitGroup
	opsPerNode := b.N / nodeCount

	for nodeIdx := 0; nodeIdx < nodeCount; nodeIdx++ {
		wg.Add(1)
		go func(idx int, ds *storage.DistributedStorage) {
			defer wg.Done()
			for i := 0; i < opsPerNode; i++ {
				key := fmt.Sprintf("node%d_key_%d", idx, i)
				err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					tx.Write(key, "value")
					return nil
				})
				if err != nil {
					log.Printf("Узел %d ошибка: %v", idx, err)
				}
			}
		}(nodeIdx, nodes[nodeIdx])
	}

	wg.Wait()
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "writes/sec")
	b.ReportMetric(float64(nodeCount), "nodes")
}

// BenchmarkDistributedBatchWrites измеряет пакетные записи
func BenchmarkDistributedBatchWrites(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	batchSizes := []int{1, 10, 50, 100}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("BatchSize%d", batchSize), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					for j := 0; j < batchSize; j++ {
						key := fmt.Sprintf("key_%d_%d", i, j)
						tx.Write(key, "value")
					}
					return nil
				})
				if err != nil {
					b.Fatalf("Ошибка пакетной записи: %v", err)
				}
			}
			b.StopTimer()

			totalOps := b.N * batchSize
			b.ReportMetric(float64(totalOps)/b.Elapsed().Seconds(), "writes/sec")
		})
	}
}

// BenchmarkDistributedLatency измеряет задержки операций
func BenchmarkDistributedLatency(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	latencies := make([]time.Duration, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			key := fmt.Sprintf("key_%d", i)
			tx.Write(key, "value")
			return nil
		})
		latencies[i] = time.Since(start)

		if err != nil {
			b.Fatalf("Ошибка: %v", err)
		}
	}
	b.StopTimer()

	// Вычисляем статистику
	var total time.Duration
	min := latencies[0]
	max := latencies[0]

	for _, lat := range latencies {
		total += lat
		if lat < min {
			min = lat
		}
		if lat > max {
			max = lat
		}
	}

	avg := total / time.Duration(len(latencies))

	b.ReportMetric(float64(min.Microseconds()), "min_μs")
	b.ReportMetric(float64(avg.Microseconds()), "avg_μs")
	b.ReportMetric(float64(max.Microseconds()), "max_μs")
}

// BenchmarkDistributedContention измеряет конкуренцию за ключи
func BenchmarkDistributedContention(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Инициализируем набор ключей
	keyCount := 10
	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("shared_key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "0")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("shared_key_%d", i%keyCount)
			err := ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
				value, _ := tx.Read(key)
				tx.Write(key, value+"1")
				return nil
			})
			if err != nil {
				log.Printf("Ошибка: %v", err)
			}
			i++
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/sec")
}

// BenchmarkDistributedReadHeavy измеряет производительность при большой нагрузке на чтение
func BenchmarkDistributedReadHeavy(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 90% чтений, 10% записей
			if i%10 == 0 {
				key := fmt.Sprintf("key_%d", i%1000)
				_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					tx.Write(key, "updated")
					return nil
				})
			} else {
				key := fmt.Sprintf("key_%d", i%1000)
				_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					_, _ = tx.Read(key)
					return nil
				})
			}
			i++
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/sec")
}

// BenchmarkDistributedWriteHeavy измеряет производительность при большой нагрузке на запись
func BenchmarkDistributedWriteHeavy(b *testing.B) {
	ds := SetupDistributedStorage(b)
	if ds == nil {
		return
	}

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	time.Sleep(500 * time.Millisecond)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 90% записей, 10% чтений
			if i%10 == 0 {
				key := fmt.Sprintf("key_%d", i%1000)
				_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					_, _ = tx.Read(key)
					return nil
				})
			} else {
				key := fmt.Sprintf("key_%d", i%1000)
				_ = ds.RunTransaction(func(tx *storage.DistributedTransaction) error {
					tx.Write(key, "updated")
					return nil
				})
			}
			i++
		}
	})
	b.StopTimer()

	b.ReportMetric(float64(b.N)/b.Elapsed().Seconds(), "ops/sec")
}
