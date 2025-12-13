package core

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHLCBasic(t *testing.T) {
	clock := NewLClock()

	// Получаем начальное время
	ts1 := clock.Get()
	require.Greater(t, ts1, uint64(0), "Initial timestamp should be greater than 0")

	// Локальное событие
	ts2 := clock.SendOrLocal()
	require.Greater(t, ts2, ts1, "SendOrLocal should increase timestamp")

	// Ещё одно локальное событие
	ts3 := clock.SendOrLocal()
	require.Greater(t, ts3, ts2, "Timestamps should be monotonically increasing")
}

func TestHLCMonotonicity(t *testing.T) {
	clock := NewLClock()

	prev := clock.Get()
	for i := 0; i < 100; i++ {
		current := clock.SendOrLocal()
		require.Greater(t, current, prev, "Timestamps must be strictly monotonic")
		prev = current
	}
}

func TestHLCRecv(t *testing.T) {
	clock1 := NewLClock()
	clock2 := NewLClock()

	// Clock1 отправляет событие
	ts1 := clock1.SendOrLocal()

	// Небольшая задержка для имитации сетевой задержки
	time.Sleep(1 * time.Millisecond)

	// Clock2 получает сообщение с timestamp от clock1
	ts2 := clock2.Recv(ts1)

	// Timestamp clock2 должен быть больше полученного
	require.Greater(t, ts2, ts1, "Receiving timestamp should increase local time")

	// Clock2 продолжает генерировать события
	ts3 := clock2.SendOrLocal()
	require.Greater(t, ts3, ts2, "Local events after Recv should maintain monotonicity")
}

func TestHLCConcurrency(t *testing.T) {
	clock := NewLClock()
	numGoroutines := 10
	eventsPerGoroutine := 100

	var wg sync.WaitGroup
	timestamps := make(chan uint64, numGoroutines*eventsPerGoroutine)

	// Запускаем несколько горутин, генерирующих события
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < eventsPerGoroutine; j++ {
				ts := clock.SendOrLocal()
				timestamps <- ts
			}
		}()
	}

	wg.Wait()
	close(timestamps)

	// Собираем все timestamps
	var allTimestamps []uint64
	for ts := range timestamps {
		allTimestamps = append(allTimestamps, ts)
	}

	// Проверяем, что все timestamps уникальны
	uniqueMap := make(map[uint64]bool)
	for _, ts := range allTimestamps {
		require.False(t, uniqueMap[ts], "All timestamps should be unique")
		uniqueMap[ts] = true
	}

	require.Len(t, uniqueMap, numGoroutines*eventsPerGoroutine, "Should have exactly the expected number of unique timestamps")
}

func TestHLCPhysicalTimeProgression(t *testing.T) {
	clock := NewLClock()

	ts1 := clock.SendOrLocal()

	// Ждем, чтобы физическое время продвинулось
	time.Sleep(10 * time.Millisecond)

	ts2 := clock.SendOrLocal()

	// После задержки новый timestamp должен отражать физическое время
	require.Greater(t, ts2, ts1, "Timestamp should reflect physical time progression")

	// Разница должна быть примерно равна задержке (в наносекундах)
	diff := ts2 - ts1
	expectedMinDiff := uint64(10 * time.Millisecond)
	require.GreaterOrEqual(t, diff, expectedMinDiff, "Timestamp difference should reflect sleep duration")
}

func TestHLCGetTimestamp(t *testing.T) {
	clock := NewLClock()

	// Генерируем несколько событий
	clock.SendOrLocal()
	clock.SendOrLocal()

	// Получаем полный timestamp
	ts := clock.GetTimestamp()
	require.Greater(t, ts.WallTime, uint64(0), "WallTime should be set")
	require.GreaterOrEqual(t, ts.Logical, uint64(0), "Logical should be non-negative")
}

func TestHLCMultipleRecv(t *testing.T) {
	clock := NewLClock()

	// Получаем несколько сообщений с разными timestamps
	ts1 := clock.Recv(1000)
	require.Greater(t, ts1, uint64(0))

	ts2 := clock.Recv(2000)
	require.Greater(t, ts2, ts1, "Second Recv should produce larger timestamp")

	ts3 := clock.Recv(1500) // Более старый timestamp
	require.Greater(t, ts3, ts2, "Even old message timestamp should maintain monotonicity")
}

func BenchmarkHLCGet(b *testing.B) {
	clock := NewLClock()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clock.Get()
	}
}

func BenchmarkHLCSendOrLocal(b *testing.B) {
	clock := NewLClock()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clock.SendOrLocal()
	}
}

func BenchmarkHLCRecv(b *testing.B) {
	clock := NewLClock()
	msgTimestamp := uint64(time.Now().UnixNano())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clock.Recv(msgTimestamp)
		msgTimestamp++
	}
}

func BenchmarkHLCConcurrentSendOrLocal(b *testing.B) {
	clock := NewLClock()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = clock.SendOrLocal()
		}
	})
}
