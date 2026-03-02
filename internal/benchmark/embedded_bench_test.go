package benchmark

import (
	"fmt"
	"math/rand"
	"petacore/internal/core"
	"petacore/internal/distributed"
	"petacore/internal/runtime/executor"
	"petacore/internal/runtime/rsql/visitor"
	"petacore/internal/storage"
	"petacore/sdk/pmem"
	"sort"
	"testing"
	"time"
)

// Embedded бенчмарк - тестирует производительность напрямую через executor
// без сетевого слоя (без pgx, wire protocol, haproxy)

var (
	benchStorage   *storage.DistributedStorageVClock
	benchAllocator pmem.Allocator

	benchRows = 10
)

// setupEmbeddedBench создает storage и заполняет тестовую таблицу
func setupEmbeddedBench(b *testing.B) {
	b.Helper()

	if benchStorage != nil {
		return // Уже инициализировано
	}

	// Используем in-memory store для чистоты эксперимента
	kv := distributed.NewInMemoryStore()
	benchStorage = storage.NewDistributedStorageVClock(
		kv,
		"bench_node1",
		1, // totalNodes
		core.SnapshotIsolation,
		1, // minAcks
	)

	if err := benchStorage.Start(); err != nil {
		b.Fatalf("failed to start storage: %v", err)
	}

	// Создаем arena и allocator
	arena, err := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	if err != nil {
		b.Fatalf("failed to create arena: %v", err)
	}
	benchAllocator = arena

	// Создаем таблицу
	createTableSQL := `CREATE TABLE bench_users (
		id INT PRIMARY KEY,
		name TEXT,
		email TEXT,
		age INT,
		score FLOAT,
		active BOOL
	)`

	stmt, err := visitor.ParseSQL(benchAllocator, createTableSQL)
	if err != nil {
		b.Fatalf("failed to parse CREATE TABLE: %v", err)
	}

	sessionParams := make(map[string]string)
	_, err = executor.ExecuteStatement(benchAllocator, stmt, benchStorage, sessionParams)
	if err != nil {
		b.Logf("create table error (may already exist): %v", err)
	}

	// Заполняем таблицу
	for i := 0; i < benchRows; i++ {
		insertSQL := fmt.Sprintf(
			"INSERT INTO bench_users (id, name, email, age, score, active) VALUES (%d, 'user_%d', 'user_%d@example.com', %d, %.2f, %v)",
			i, i, i,
			20+rand.Intn(50),
			rand.Float64()*1000,
			i%2 == 0,
		)

		stmt, err := visitor.ParseSQL(benchAllocator, insertSQL)
		if err != nil {
			b.Fatalf("failed to parse INSERT: %v", err)
		}

		_, err = executor.ExecuteStatement(benchAllocator, stmt, benchStorage, sessionParams)
		if err != nil {
			b.Fatalf("failed to insert row %d: %v", i, err)
		}
	}

	b.Logf("Setup complete: inserted %d rows", benchRows)
}

// BenchmarkEmbeddedSelectByID - точечный SELECT по PK без сетевого слоя
func BenchmarkEmbeddedSelectByID(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// Создаем свежий allocator для каждого запроса
		alloc, err := pmem.NewMmapArena(1024 * 1024 * 10) // 10MB

		id := rand.Intn(benchRows)
		query := fmt.Sprintf("SELECT id, name, email, age, score, active FROM bench_users WHERE id = %d", id)

		stmt, err := visitor.ParseSQL(alloc, query)
		if err != nil {
			b.Fatalf("failed to parse SELECT: %v", err)
		}

		result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
		if err != nil {
			b.Fatalf("failed to execute SELECT: %v", err)
		}

		if result == nil || len(result.Rows) == 0 {
			b.Fatalf("expected result, got nil or empty")
		}

		alloc.Reset()
	}
}

// BenchmarkEmbeddedSimpleSelect - простой SELECT двух колонок
func BenchmarkEmbeddedSimpleSelect(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)

	b.ResetTimer()
	b.ReportAllocs()

	alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	for i := 0; i < b.N; i++ {

		id := rand.Intn(benchRows)
		query := fmt.Sprintf("SELECT id, name FROM bench_users WHERE id = %d", id)

		stmt, err := visitor.ParseSQL(alloc, query)
		if err != nil {
			b.Fatalf("failed to parse SELECT: %v", err)
		}

		result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
		if err != nil {
			b.Fatalf("failed to execute SELECT: %v", err)
		}

		if result == nil || len(result.Rows) == 0 {
			b.Fatalf("expected result, got nil or empty")
		}

		alloc.Reset()
	}
}

// BenchmarkEmbeddedSelectN - SELECT с разными LIMIT
func BenchmarkEmbeddedSelectN(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)
	limits := []int{1, 10}

	for _, limit := range limits {
		limit := limit
		b.Run(fmt.Sprintf("rows=%d", limit), func(b *testing.B) {
			query := fmt.Sprintf("SELECT id, name, email, age, score, active FROM bench_users LIMIT %d", limit)
			alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {

				stmt, err := visitor.ParseSQL(alloc, query)
				if err != nil {
					b.Fatalf("failed to parse SELECT: %v", err)
				}

				result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
				if err != nil {
					b.Fatalf("failed to execute SELECT: %v", err)
				}

				if result == nil {
					b.Fatalf("expected result, got nil")
				}

				alloc.Reset()
			}
		})
	}
}

// BenchmarkEmbeddedWhere - WHERE с разными типами
func BenchmarkEmbeddedWhere(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)

	tests := []struct {
		name  string
		query string
	}{
		{"where_int", "SELECT id FROM bench_users WHERE age = 25"},
		{"where_text", "SELECT id FROM bench_users WHERE name = 'user_5'"},
		{"where_bool", "SELECT id FROM bench_users WHERE active = true"},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
			for i := 0; i < b.N; i++ {

				stmt, err := visitor.ParseSQL(alloc, tc.query)
				if err != nil {
					b.Fatalf("failed to parse: %v", err)
				}

				result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
				if err != nil {
					b.Fatalf("failed to execute: %v", err)
				}

				_ = result
				alloc.Reset()
			}
		})
	}
}

// BenchmarkEmbeddedSelectStar - SELECT * vs явные колонки
func BenchmarkEmbeddedSelectStar(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)

	tests := []struct {
		name  string
		query string
	}{
		{"select_star", "SELECT * FROM bench_users LIMIT 10"},
		{"select_explicit", "SELECT id, name, email, age, score, active FROM bench_users LIMIT 10"},
	}

	for _, tc := range tests {
		tc := tc
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
			for i := 0; i < b.N; i++ {

				stmt, err := visitor.ParseSQL(alloc, tc.query)
				if err != nil {
					b.Fatalf("failed to parse: %v", err)
				}

				result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
				if err != nil {
					b.Fatalf("failed to execute: %v", err)
				}

				_ = result
				alloc.Reset()
			}
		})
	}
}

// BenchmarkEmbeddedLatencyProfile - профиль латентности
func BenchmarkEmbeddedLatencyProfile(b *testing.B) {
	setupEmbeddedBench(b)

	sessionParams := make(map[string]string)
	query := "SELECT id, name, score FROM bench_users WHERE id = 5"

	latencies := make([]int64, 0, b.N)

	b.ResetTimer()

	alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	for i := 0; i < b.N; i++ {

		start := time.Now().UnixNano()

		stmt, err := visitor.ParseSQL(alloc, query)
		if err != nil {
			b.Fatalf("failed to parse: %v", err)
		}

		result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
		if err != nil {
			b.Fatalf("failed to execute: %v", err)
		}

		_ = result
		latencies = append(latencies, time.Now().UnixNano()-start)
		alloc.Reset()
	}

	if len(latencies) == 0 {
		return
	}

	// Сортируем для вычисления перцентилей
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	n := len(latencies)
	p50 := latencies[n/2]
	p95 := latencies[n*95/100]
	p99 := latencies[n*99/100]
	p100 := latencies[n-1]

	b.ReportMetric(float64(p50), "p50_ns")
	b.ReportMetric(float64(p95), "p95_ns")
	b.ReportMetric(float64(p99), "p99_ns")
	b.ReportMetric(float64(p100), "p100_ns")
}

// BenchmarkEmbeddedParseOnly - только парсинг (без выполнения)
func BenchmarkEmbeddedParseOnly(b *testing.B) {
	setupEmbeddedBench(b)

	query := "SELECT id, name, email, age, score, active FROM bench_users WHERE id = 5"

	b.ResetTimer()
	b.ReportAllocs()

	alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	for i := 0; i < b.N; i++ {

		stmt, err := visitor.ParseSQL(alloc, query)
		if err != nil {
			b.Fatalf("failed to parse: %v", err)
		}
		_ = stmt

		alloc.Reset()
	}
}

// BenchmarkEmbeddedExecuteOnly - только выполнение (с предпарсенным statement)
func BenchmarkEmbeddedExecuteOnly(b *testing.B) {
	setupEmbeddedBench(b)

	query := "SELECT id, name, email, age, score, active FROM bench_users WHERE id = 5"
	sessionParams := make(map[string]string)

	// Парсим один раз с отдельным allocator
	parseAlloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	stmt, err := visitor.ParseSQL(parseAlloc, query)
	if err != nil {
		b.Fatalf("failed to parse: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
	for i := 0; i < b.N; i++ {

		result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
		if err != nil {
			b.Fatalf("failed to execute: %v", err)
		}
		_ = result

		alloc.Reset()
	}
}

// BenchmarkEmbeddedParallel - параллельные запросы
func BenchmarkEmbeddedParallel(b *testing.B) {
	setupEmbeddedBench(b)

	for _, numConns := range []int{1, 2, 4, 8, 16} {
		numConns := numConns
		b.Run(fmt.Sprintf("goroutines=%d", numConns), func(b *testing.B) {
			b.SetParallelism(numConns)
			b.RunParallel(func(pb *testing.PB) {
				sessionParams := make(map[string]string)
				query := "SELECT id, name, email FROM bench_users WHERE id = 5"
				alloc, _ := pmem.NewMmapArena(10 * 1024 * 1024) // 10MB
				for pb.Next() {

					stmt, err := visitor.ParseSQL(alloc, query)
					if err != nil {
						b.Errorf("failed to parse: %v", err)
						return
					}

					result, err := executor.ExecuteStatement(alloc, stmt, benchStorage, sessionParams)
					if err != nil {
						b.Errorf("failed to execute: %v", err)
						return
					}

					_ = result
					alloc.Reset()
				}
			})
		})
	}
}
