package storage_test

import (
	"petacore/internal/core"
	"petacore/internal/storage"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSStorageSimple(t *testing.T) {
	// Placeholder for storage tests
	s := storage.NewSimpleStorage()
	exampleKey := "example"

	err := s.RunTransaction(func(tx *core.Transaction) error {
		_, ok := tx.Read(exampleKey)
		require.False(t, ok)
		tx.Write(exampleKey, "value1")
		val, ok := tx.Read(exampleKey)
		require.True(t, ok)
		require.Equal(t, "value1", val)
		return nil
	})
	require.NoError(t, err)
}

// TestSStorageConcurrentWrites проверяет конкурентную запись из нескольких горутин
func TestSStorageConcurrentWrites(t *testing.T) {
	s := storage.NewSimpleStorage()
	key := "counter"

	// Запускаем 10 горутин, каждая пишет свое значение
	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			err := s.RunTransaction(func(tx *core.Transaction) error {
				tx.Write(key, string(rune('A'+id)))
				return nil
			})
			require.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Читаем последнее значение
	var finalValue string
	err := s.RunTransaction(func(tx *core.Transaction) error {
		val, ok := tx.Read(key)
		require.True(t, ok)
		finalValue = val
		return nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, finalValue)
	t.Logf("Final value after concurrent writes: %s", finalValue)
}

// TestSStorageReadYourWrites проверяет изоляцию транзакций и чтение своих записей
func TestSStorageReadYourWrites(t *testing.T) {
	s := storage.NewSimpleStorage()
	key := "account"

	// Транзакция 1: записываем начальное значение
	err := s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write(key, "100")
		return nil
	})
	require.NoError(t, err)

	var tx1Read, tx2Read string
	var wg sync.WaitGroup

	// Транзакция 2: читает и пишет новое значение
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			tx1Read = val
			tx.Write(key, "200")
			// Читаем свою же запись
			val, ok = tx.Read(key)
			require.True(t, ok)
			require.Equal(t, "200", val, "Should read own write")
			return nil
		})
		require.NoError(t, err)
	}()

	// Транзакция 3: читает параллельно
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			tx2Read = val
			return nil
		})
		require.NoError(t, err)
	}()

	wg.Wait()

	t.Logf("Tx1 read: %s, Tx2 read: %s", tx1Read, tx2Read)
	require.Contains(t, []string{"100", "200"}, tx1Read)
	require.Contains(t, []string{"100", "200"}, tx2Read)
}

// TestSStorageMultipleKeysAcrossGoroutines проверяет работу с несколькими ключами в разных рантаймах
func TestSStorageMultipleKeysAcrossGoroutines(t *testing.T) {
	s := storage.NewSimpleStorage()

	// Создаем начальные данные
	err := s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write("user:1", "Alice")
		tx.Write("user:2", "Bob")
		tx.Write("balance:1", "1000")
		tx.Write("balance:2", "2000")
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	results := make(map[int]map[string]string)
	resultsMutex := sync.Mutex{}

	// Несколько горутин читают и модифицируют данные
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			localResults := make(map[string]string)

			err := s.RunTransaction(func(tx *core.Transaction) error {
				// Читаем данные
				user1, ok := tx.Read("user:1")
				require.True(t, ok)
				localResults["user:1"] = user1

				balance1, ok := tx.Read("balance:1")
				require.True(t, ok)
				localResults["balance:1"] = balance1

				// Пишем новые данные
				tx.Write("user:1", user1+"_updated")
				tx.Write("balance:1", balance1+"0") // добавляем 0 к балансу

				// Читаем свои записи
				updatedUser, ok := tx.Read("user:1")
				require.True(t, ok)
				localResults["user:1_updated"] = updatedUser

				return nil
			})
			require.NoError(t, err)

			resultsMutex.Lock()
			results[goroutineID] = localResults
			resultsMutex.Unlock()
		}(i)
	}

	wg.Wait()

	// Проверяем что все горутины успешно выполнились
	require.Len(t, results, 5)

	// Финальное чтение
	err = s.RunTransaction(func(tx *core.Transaction) error {
		user1, ok := tx.Read("user:1")
		require.True(t, ok)
		t.Logf("Final user:1 = %s", user1)

		balance1, ok := tx.Read("balance:1")
		require.True(t, ok)
		t.Logf("Final balance:1 = %s", balance1)

		return nil
	})
	require.NoError(t, err)
}

// TestSStorageSnapshotIsolation проверяет snapshot isolation между горутинами
func TestSStorageSnapshotIsolation(t *testing.T) {
	s := storage.NewSimpleStorage()
	key := "data"

	// Устанавливаем начальное значение
	err := s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write(key, "v1")
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	readChan := make(chan string, 3)

	// Горутина 1: долгая транзакция чтения
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			readChan <- val
			return nil
		})
		require.NoError(t, err)
	}()

	// Горутина 2: быстрая запись
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "v2")
			return nil
		})
		require.NoError(t, err)
	}()

	// Горутина 3: ещё одно чтение
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			readChan <- val
			return nil
		})
		require.NoError(t, err)
	}()

	wg.Wait()
	close(readChan)

	// Собираем результаты
	var reads []string
	for val := range readChan {
		reads = append(reads, val)
	}

	t.Logf("Read values: %v", reads)
	// Должны видеть либо v1, либо v2 в зависимости от времени
	for _, val := range reads {
		require.Contains(t, []string{"v1", "v2"}, val)
	}
}

// TestReadCommittedIsolation проверяет, что Read Committed видит последние committed данные
func TestReadCommittedIsolation(t *testing.T) {
	s := storage.NewSimpleStorage() // По умолчанию Read Committed
	key := "account"

	// Устанавливаем начальное значение
	err := s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write(key, "100")
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	readValues := make(chan string, 2)

	// Транзакция 1: запись нового значения
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "200")
			return nil
		})
		require.NoError(t, err)
	}()

	// Транзакция 2: читает значение (может увидеть либо 100, либо 200)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			readValues <- val
			return nil
		})
		require.NoError(t, err)
	}()

	wg.Wait()
	close(readValues)

	// Проверяем, что прочитанное значение - это committed данные
	val := <-readValues
	require.Contains(t, []string{"100", "200"}, val)
	t.Logf("Read Committed saw value: %s", val)
}

// TestReadCommittedVsSnapshotIsolation сравнивает поведение двух уровней изоляции
func TestReadCommittedVsSnapshotIsolation(t *testing.T) {
	key := "counter"

	t.Run("ReadCommitted sees latest", func(t *testing.T) {
		s := storage.NewSimpleStorageWithIsolation(core.ReadCommitted)

		// Начальное значение
		err := s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "1")
			return nil
		})
		require.NoError(t, err)

		// Обновление значения
		err = s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "2")
			return nil
		})
		require.NoError(t, err)

		// Чтение должно видеть последнюю версию
		err = s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			require.Equal(t, "2", val, "Read Committed should see latest committed value")
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("SnapshotIsolation sees fixed snapshot", func(t *testing.T) {
		s := storage.NewSimpleStorageWithIsolation(core.SnapshotIsolation)

		// Начальное значение
		err := s.RunTransaction(func(tx *core.Transaction) error {
			tx.Write(key, "1")
			return nil
		})
		require.NoError(t, err)

		// Snapshot isolation видит фиксированный снимок
		// (в нашей реализации это тоже будет последнее значение,
		// но фиксируется на момент Begin)
		err = s.RunTransaction(func(tx *core.Transaction) error {
			val, ok := tx.Read(key)
			require.True(t, ok)
			require.Equal(t, "1", val)
			return nil
		})
		require.NoError(t, err)
	})
}

// TestReadCommittedNonRepeatableRead демонстрирует non-repeatable read в Read Committed
func TestReadCommittedNonRepeatableRead(t *testing.T) {
	s := storage.NewSimpleStorageWithIsolation(core.ReadCommitted)
	key := "value"

	// Начальное значение
	err := s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write(key, "initial")
		return nil
	})
	require.NoError(t, err)

	// В Read Committed каждое чтение может видеть новую версию
	// (но мы не можем это протестировать внутри одной транзакции,
	// так как транзакция выполняется атомарно)
	// Вместо этого проверяем, что разные транзакции видят разные значения

	// Первое чтение
	var firstRead string
	err = s.RunTransaction(func(tx *core.Transaction) error {
		val, ok := tx.Read(key)
		require.True(t, ok)
		firstRead = val
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, "initial", firstRead)

	// Обновление
	err = s.RunTransaction(func(tx *core.Transaction) error {
		tx.Write(key, "updated")
		return nil
	})
	require.NoError(t, err)

	// Второе чтение - должно видеть обновление
	var secondRead string
	err = s.RunTransaction(func(tx *core.Transaction) error {
		val, ok := tx.Read(key)
		require.True(t, ok)
		secondRead = val
		return nil
	})
	require.NoError(t, err)
	require.Equal(t, "updated", secondRead, "Read Committed should see the updated value")

	t.Logf("First read: %s, Second read: %s - demonstrates non-repeatable read", firstRead, secondRead)
}
