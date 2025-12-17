package storage_test

import (
	"fmt"
	"petacore/internal/core"
	"petacore/internal/storage"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkVClock_SingleWrite измеряет скорость одиночных записей
func BenchmarkVClock_SingleWrite(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}
}

// BenchmarkVClock_SingleRead измеряет скорость одиночных чтений с quorum
func BenchmarkVClock_SingleRead(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	// Подготовка данных
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%1000)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Read(key)
			return nil
		})
	}
}

// BenchmarkVClock_ReadWrite измеряет скорость смешанных операций
func BenchmarkVClock_ReadWrite(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	// Подготовка данных
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key%d", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			// Чтение
			key := fmt.Sprintf("key%d", i%100)
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read(key)
				return nil
			})
		} else {
			// Запись
			key := fmt.Sprintf("key%d", i%100)
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(key, "newvalue")
				return nil
			})
		}
	}
}

// BenchmarkVClock_ConcurrentWrites измеряет throughput параллельных записей
func BenchmarkVClock_ConcurrentWrites(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(key, "value")
				return nil
			})
			i++
		}
	})
}

// BenchmarkVClock_ConcurrentReads измеряет throughput параллельных чтений
func BenchmarkVClock_ConcurrentReads(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	// Подготовка данных
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read(key)
				return nil
			})
			i++
		}
	})
}

// BenchmarkVClock_Transaction измеряет скорость транзакций
func BenchmarkVClock_Transaction(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(fmt.Sprintf("key%d", i), "value")
			tx.Read(fmt.Sprintf("key%d", i))
			return nil
		})
	}
}

// BenchmarkVClock_LargeTransaction измеряет скорость больших транзакций
func BenchmarkVClock_LargeTransaction(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			// 10 операций в одной транзакции
			for j := 0; j < 10; j++ {
				key := fmt.Sprintf("key%d_%d", i, j)
				tx.Write(key, "value")
			}
			return nil
		})
	}
}

// BenchmarkVClock_HotKey измеряет contention на горячем ключе
func BenchmarkVClock_HotKey(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	// Инициализация горячего ключа
	ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
		tx.Write("hotkey", "0")
		return nil
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Read("hotkey")
				tx.Write("hotkey", "updated")
				return nil
			})
		}
	})
}

// BenchmarkVClock_QuorumOverhead измеряет overhead от quorum проверки
func BenchmarkVClock_QuorumOverhead(b *testing.B) {
	scenarios := []struct {
		name       string
		totalNodes int
	}{
		{"1node", 1},
		{"3nodes", 3},
		{"5nodes", 5},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			ds := SetupDistributedStorageVClock(b, "node1", scenario.totalNodes, core.ReadCommitted, 0)
			if ds == nil {
				return
			}

			// Подготовка данных
			for i := 0; i < 100; i++ {
				key := fmt.Sprintf("key%d", i)
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					tx.Write(key, "value")
					return nil
				})
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := fmt.Sprintf("key%d", i%100)
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					tx.Read(key)
					return nil
				})
			}
		})
	}
}

// BenchmarkVClock_VectorClockSize измеряет влияние размера VectorClock
func BenchmarkVClock_VectorClockSize(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 5, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(key, "value")
				return nil
			})
			i++
		}
	})
}

// BenchmarkVClock_Latency измеряет распределение латентности
func BenchmarkVClock_Latency(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	latencies := make([]time.Duration, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(fmt.Sprintf("key%d", i), "value")
			return nil
		})
		latencies[i] = time.Since(start)
	}

	b.StopTimer()

	// Вычисляем перцентили
	if b.N > 0 {
		// Простая сортировка для малых N
		for i := 0; i < len(latencies)-1; i++ {
			for j := i + 1; j < len(latencies); j++ {
				if latencies[i] > latencies[j] {
					latencies[i], latencies[j] = latencies[j], latencies[i]
				}
			}
		}

		p50 := latencies[b.N/2]
		p95 := latencies[b.N*95/100]
		p99 := latencies[b.N*99/100]

		b.ReportMetric(float64(p50.Microseconds()), "p50_us")
		b.ReportMetric(float64(p95.Microseconds()), "p95_us")
		b.ReportMetric(float64(p99.Microseconds()), "p99_us")
	}
}

// BenchmarkVClock_Throughput измеряет максимальный throughput
func BenchmarkVClock_Throughput(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	var ops int64
	done := make(chan struct{})

	// Запускаем измерение throughput
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				current := atomic.LoadInt64(&ops)
				b.ReportMetric(float64(current), "ops/sec")
				atomic.StoreInt64(&ops, 0)
			case <-done:
				return
			}
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(fmt.Sprintf("key%d", i), "value")
				return nil
			})
			atomic.AddInt64(&ops, 1)
			i++
		}
	})

	close(done)
}

// BenchmarkVClock_MemoryUsage измеряет использование памяти
func BenchmarkVClock_MemoryUsage(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	b.ResetTimer()

	// Записываем много данных для измерения памяти
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value_%d_with_some_data", i)
		ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
			tx.Write(key, value)
			return nil
		})
	}
}

// BenchmarkVClock_Scalability проверяет масштабируемость с разным числом горутин
func BenchmarkVClock_Scalability(b *testing.B) {
	goroutines := []int{1, 2, 4, 8, 16, 32}

	for _, numGoroutines := range goroutines {
		b.Run(fmt.Sprintf("goroutines_%d", numGoroutines), func(b *testing.B) {
			ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
			if ds == nil {
				return
			}

			b.SetParallelism(numGoroutines)
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("key%d", i)
					ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
						tx.Write(key, "value")
						return nil
					})
					i++
				}
			})
		})
	}
}

// BenchmarkVClock_BatchWrite измеряет скорость пакетной записи
func BenchmarkVClock_BatchWrite(b *testing.B) {
	batchSizes := []int{1, 10, 50, 100}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("batch_%d", batchSize), func(b *testing.B) {
			ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
			if ds == nil {
				return
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
					for j := 0; j < batchSize; j++ {
						key := fmt.Sprintf("key%d_%d", i, j)
						tx.Write(key, "value")
					}
					return nil
				})
			}
		})
	}
}

// BenchmarkVClock_ContentionLevel измеряет производительность при разных уровнях contention
func BenchmarkVClock_ContentionLevel(b *testing.B) {
	contentionLevels := []struct {
		name     string
		keyRange int
	}{
		{"low_contention", 10000},
		{"medium_contention", 100},
		{"high_contention", 10},
	}

	for _, level := range contentionLevels {
		b.Run(level.name, func(b *testing.B) {
			ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
			if ds == nil {
				return
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					key := fmt.Sprintf("key%d", i%level.keyRange)
					ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
						tx.Write(key, "value")
						return nil
					})
					i++
				}
			})
		})
	}
}

// BenchmarkVClock_SynchronizationOverhead измеряет overhead синхронизации
func BenchmarkVClock_SynchronizationOverhead(b *testing.B) {
	ds := SetupDistributedStorageVClock(b, "node1", 1, core.ReadCommitted, 0)
	if ds == nil {
		return
	}

	var wg sync.WaitGroup

	b.ResetTimer()

	// Измеряем время с учетом синхронизации watch
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ds.RunTransaction(func(tx *storage.DistributedTransactionVClock) error {
				tx.Write(fmt.Sprintf("key%d", i), "value")
				return nil
			})
		}(i)
	}

	wg.Wait()
}
