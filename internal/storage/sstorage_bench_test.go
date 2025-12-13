package storage_test

import (
	"fmt"
	"petacore/internal/core"
	"petacore/internal/storage"
	"sync"
	"testing"
)

// BenchmarkWrite измеряет производительность последовательных записей
func BenchmarkWrite(b *testing.B) {
	s := storage.NewSimpleStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "value")
			return nil
		})
	}
}

// BenchmarkRead измеряет производительность последовательных чтений
func BenchmarkRead(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key_%d", i%1000)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			_, _ = tx.Read(key)
			return nil
		})
	}
}

// BenchmarkReadWrite измеряет производительность смешанных операций чтения/записи
func BenchmarkReadWrite(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "initial_value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		readKey := fmt.Sprintf("key_%d", i%100)
		writeKey := fmt.Sprintf("key_%d", (i+1)%100)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			_, _ = tx.Read(readKey)
			tx.Write(writeKey, "updated_value")
			return nil
		})
	}
}

// BenchmarkConcurrentWrites измеряет производительность конкурентных записей
func BenchmarkConcurrentWrites(b *testing.B) {
	s := storage.NewSimpleStorage()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i)
			_ = s.RunTransaction(func(tx *core.Transaction) error {
				tx.Write(key, "value")
				return nil
			})
			i++
		}
	})
}

// BenchmarkConcurrentReads измеряет производительность конкурентных чтений
func BenchmarkConcurrentReads(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i%1000)
			_ = s.RunTransaction(func(tx *core.Transaction) error {
				_, _ = tx.Read(key)
				return nil
			})
			i++
		}
	})
}

// BenchmarkConcurrentReadWrite измеряет производительность конкурентных смешанных операций
func BenchmarkConcurrentReadWrite(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "initial_value")
			return nil
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			readKey := fmt.Sprintf("key_%d", i%100)
			writeKey := fmt.Sprintf("key_%d", (i+1)%100)
			_ = s.RunTransaction(func(tx *core.Transaction) error {
				_, _ = tx.Read(readKey)
				tx.Write(writeKey, "updated_value")
				return nil
			})
			i++
		}
	})
}

// BenchmarkMultiKeyTransaction измеряет производительность транзакций с несколькими ключами
func BenchmarkMultiKeyTransaction(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			// Читаем 5 ключей
			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("key_%d", (i*5+j)%100)
				_, _ = tx.Read(key)
			}
			// Пишем 5 ключей
			for j := 0; j < 5; j++ {
				key := fmt.Sprintf("key_%d", (i*5+j)%100)
				tx.Write(key, "new_value")
			}
			return nil
		})
	}
}

// BenchmarkHighContentionWrites измеряет производительность при высокой конкуренции
func BenchmarkHighContentionWrites(b *testing.B) {
	s := storage.NewSimpleStorage()
	const hotKey = "hot_key"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.RunTransaction(func(tx *core.Transaction) error {
				tx.Write(hotKey, "value")
				return nil
			})
		}
	})
}

// BenchmarkMVCCVersionAccumulation измеряет производительность с накоплением версий
func BenchmarkMVCCVersionAccumulation(b *testing.B) {
	s := storage.NewSimpleStorage()
	key := "versioned_key"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, fmt.Sprintf("version_%d", i))
			return nil
		})
	}
}

// BenchmarkLargeTransaction измеряет производительность больших транзакций
func BenchmarkLargeTransaction(b *testing.B) {
	s := storage.NewSimpleStorage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			// Пишем 100 разных ключей в одной транзакции
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("large_tx_key_%d_%d", i, j)
				tx.Write(key, "value")
			}
			return nil
		})
	}
}

// BenchmarkRealisticWorkload моделирует реалистичную нагрузку (80% чтение, 20% запись)
func BenchmarkRealisticWorkload(b *testing.B) {
	s := storage.NewSimpleStorage()

	// Предварительно записываем данные
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key_%d", i)
		_ = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "value")
			return nil
		})
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key_%d", i%1000)
			_ = s.RunTransaction(func(tx *core.Transaction) error {
				if i%5 == 0 {
					// 20% записей
					tx.Write(key, "updated_value")
				} else {
					// 80% чтений
					_, _ = tx.Read(key)
				}
				return nil
			})
			i++
		}
	})
}

// BenchmarkConcurrentTransactionBatch измеряет пропускную способность при пакетной обработке
func BenchmarkConcurrentTransactionBatch(b *testing.B) {
	s := storage.NewSimpleStorage()

	b.ResetTimer()
	var wg sync.WaitGroup
	workersCount := 10
	txPerWorker := b.N / workersCount

	for w := 0; w < workersCount; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < txPerWorker; i++ {
				key := fmt.Sprintf("worker_%d_key_%d", workerID, i)
				_ = s.RunTransaction(func(tx *core.Transaction) error {
					tx.Write(key, "value")
					return nil
				})
			}
		}(w)
	}
	wg.Wait()
}
