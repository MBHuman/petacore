package petacore_test

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx"
)

// ─────────────────────────────────────────────
// Setup
// ─────────────────────────────────────────────

// benchPool — один пул на весь процесс, создаётся один раз.
var benchPool *pgx.ConnPool

// benchRows — количество строк, которые вставляются в тестовую таблице.
// По умолчанию 10, можно изменить при необходимости.
var benchRows = 10

func getBenchPool(b *testing.B) *pgx.ConnPool {
	b.Helper()
	if benchPool != nil {
		return benchPool
	}
	var err error
	benchPool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "testdb",
			User:     "mbhuman",
		},
		MaxConnections: 16,
	})
	if err != nil {
		b.Fatalf("pool: %v", err)
	}
	return benchPool
}

// setupTable создаёт и заполняет таблицу. Вызывается один раз на весь бенч-файл.
func setupTable(b *testing.B, pool *pgx.ConnPool, rows int) {
	b.Helper()

	if _, err := pool.Exec("DROP TABLE bench_users"); err != nil {
		b.Fatalf("drop table: %v", err)
	}

	// Ensure table exists and is empty. CREATE TABLE IF NOT EXISTS makes this idempotent
	// in case another process/node creates it concurrently; TRUNCATE clears old data.
	if _, err := pool.Exec(`CREATE TABLE IF NOT EXISTS bench_users (
		id     INT PRIMARY KEY,
		name   TEXT,
		email  TEXT,
		age    INT,
		score  FLOAT,
		active BOOL
	)`); err != nil {
		b.Fatalf("create table: %v", err)
	}

	// Ensure the table is empty before inserting test data. Some backends may delay
	// visibility after creation, so allow a short retry for TRUNCATE to succeed.
	for i := 0; i < 5; i++ {
		if _, err := pool.Exec("TRUNCATE bench_users"); err != nil {
			if i == 4 {
				b.Fatalf("truncate table: %v", err)
			}
			time.Sleep(50 * time.Millisecond)
			continue
		}
		break
	}

	start := 0
	end := rows
	var sb strings.Builder
	sb.WriteString("INSERT INTO bench_users (id, name, email, age, score, active) VALUES ")
	for i := start; i < end; i++ {
		if i > start {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "(%d,'user_%d','user_%d@example.com',%d,%.2f,%v)",
			i, i, i,
			20+rand.Intn(50),
			rand.Float64()*1000,
			i%2 == 0,
		)
	}
	if _, err := pool.Exec(sb.String()); err != nil {
		b.Fatalf("insert batch start=%d: %v", start, err)
	}

}

// tableOnce гарантирует что setupTable вызывается ровно один раз.
// Используем простой флаг — бенчи в одном процессе запускаются последовательно.
var tableReady bool

func ensureTable(b *testing.B, pool *pgx.ConnPool) {
	b.Helper()
	if tableReady {
		return
	}
	setupTable(b, pool, benchRows)
	tableReady = true
}

// ─────────────────────────────────────────────
// 1. Prepared statement — точечный lookup по PK
// ─────────────────────────────────────────────

func BenchmarkPreparedSelectByID(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	if _, err := conn.Prepare("sel_by_id",
		"SELECT id, name, email, age, score, active FROM bench_users WHERE id = $1",
	); err != nil {
		b.Fatalf("prepare: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var id, age int
		var name, email string
		var score float64
		var active bool
		err := conn.QueryRow("sel_by_id", rand.Intn(benchRows)).
			Scan(&id, &name, &email, &age, &score, &active)
		if err != nil {
			b.Fatalf("query: %v", err)
		}
	}
}

// ─────────────────────────────────────────────
// 2. Simple query — без prepare
// ─────────────────────────────────────────────

func BenchmarkSimpleQueryByID(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var id int
		var name string
		err := conn.QueryRow(
			fmt.Sprintf("SELECT id, name FROM bench_users WHERE id = %d", rand.Intn(benchRows)),
		).Scan(&id, &name)
		if err != nil {
			b.Fatalf("query: %v", err)
		}
	}
}

// ─────────────────────────────────────────────
// 3. SELECT с разными размерами результата
// ─────────────────────────────────────────────

func BenchmarkSelectN(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	// Prepare все варианты один раз до суб-бенчей
	limits := []int{1, 10, 100, 1000, 10_000}
	for _, limit := range limits {
		name := fmt.Sprintf("sel_n_%d", limit)
		sql := fmt.Sprintf(
			"SELECT id, name, email, age, score, active FROM bench_users LIMIT %d", limit)
		if _, err := conn.Prepare(name, sql); err != nil {
			b.Fatalf("prepare limit=%d: %v", limit, err)
		}
	}

	for _, limit := range limits {
		limit := limit
		b.Run(fmt.Sprintf("rows=%d", limit), func(b *testing.B) {
			name := fmt.Sprintf("sel_n_%d", limit)
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				rows, err := conn.Query(name)
				if err != nil {
					b.Fatalf("query: %v", err)
				}
				for rows.Next() {
					var id, age int
					var name, email string
					var score float64
					var active bool
					rows.Scan(&id, &name, &email, &age, &score, &active)
				}
				rows.Close()
			}
		})
	}
}

// ─────────────────────────────────────────────
// 4. SELECT * vs явные колонки
// ─────────────────────────────────────────────

func BenchmarkSelectStar(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	conn.Prepare("sel_star", "SELECT * FROM bench_users LIMIT 100")
	conn.Prepare("sel_explicit",
		"SELECT id, name, email, age, score, active FROM bench_users LIMIT 100")

	b.Run("select_star", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, _ := conn.Query("sel_star")
			rows.Close()
		}
	})

	b.Run("select_explicit", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, _ := conn.Query("sel_explicit")
			rows.Close()
		}
	})
}

// ─────────────────────────────────────────────
// 5. WHERE с разными типами
// ─────────────────────────────────────────────

func BenchmarkWhereTypes(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	conn.Prepare("sel_int", "SELECT id FROM bench_users WHERE age = $1")
	conn.Prepare("sel_bool", "SELECT id FROM bench_users WHERE active = $1::bool")
	conn.Prepare("sel_text", "SELECT id FROM bench_users WHERE name = $1")

	b.Run("where_int", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, err := conn.Query("sel_int", 20+rand.Intn(50))
			if err != nil {
				b.Fatalf("query: %v", err)
			}
			rows.Close()
		}
	})

	b.Run("where_text", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows, err := conn.Query("sel_text",
				fmt.Sprintf("user_%d", rand.Intn(benchRows)))
			if err != nil {
				b.Fatalf("query: %v", err)
			}
			rows.Close()
		}
	})
}

// ─────────────────────────────────────────────
// 6. Параллельные запросы — реальный throughput
// ─────────────────────────────────────────────

func BenchmarkParallelSelect(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	for _, numConns := range []int{1, 2, 4, 8, 16, 1000, 10000} {
		numConns := numConns
		b.Run(fmt.Sprintf("conns=%d", numConns), func(b *testing.B) {
			b.SetParallelism(numConns)
			b.RunParallel(func(pb *testing.PB) {
				// Каждая горутина берёт коннект из пула (не создаёт новый)
				conn, err := pool.Acquire()
				if err != nil {
					b.Errorf("acquire: %v", err)
					return
				}
				defer pool.Release(conn)

				// Имя стейтмента уникально per-connection в pgx v3
				stmtName := fmt.Sprintf("par_sel_%p", conn)
				conn.Prepare(stmtName,
					"SELECT id, name, email FROM bench_users WHERE id = $1")

				var id int
				var name, email string
				for pb.Next() {
					err := conn.QueryRow(stmtName, rand.Intn(benchRows)).
						Scan(&id, &name, &email)
					if err != nil {
						b.Errorf("query: %v", err)
					}
				}
			})
		})
	}
}

// ─────────────────────────────────────────────
// 7. Aggregates
// ─────────────────────────────────────────────

func BenchmarkAggregates(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	queries := []struct {
		name string
		sql  string
	}{
		{"count_star", "SELECT COUNT(*) FROM bench_users"},
		{"count_where", "SELECT COUNT(*) FROM bench_users WHERE active = true"},
		{"avg", "SELECT AVG(score) FROM bench_users"},
		{"group_by", "SELECT active, COUNT(*) FROM bench_users GROUP BY active"},
		{"group_by_age", "SELECT age, AVG(score) FROM bench_users GROUP BY age"},
	}

	// Prepare всё заранее
	for _, q := range queries {
		if _, err := conn.Prepare(q.name, q.sql); err != nil {
			b.Fatalf("prepare %s: %v", q.name, err)
		}
	}

	for _, q := range queries {
		q := q
		b.Run(q.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rows, err := conn.Query(q.name)
				if err != nil {
					b.Fatalf("query %s: %v", q.name, err)
				}
				rows.Close()
			}
		})
	}
}

// ─────────────────────────────────────────────
// 8. Latency percentiles
// ─────────────────────────────────────────────

func BenchmarkLatencyProfile(b *testing.B) {
	pool := getBenchPool(b)
	ensureTable(b, pool)

	conn, err := pool.Acquire()
	if err != nil {
		b.Fatalf("acquire: %v", err)
	}
	defer pool.Release(conn)

	if _, err := conn.Prepare("lat_sel",
		"SELECT id, name, score FROM bench_users WHERE id = $1",
	); err != nil {
		b.Fatalf("prepare: %v", err)
	}

	latencies := make([]int64, 0, b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now().UnixNano()
		var id int
		var name string
		var score float64
		conn.QueryRow("lat_sel", rand.Intn(benchRows)).Scan(&id, &name, &score)
		latencies = append(latencies, time.Now().UnixNano()-start)
	}

	if len(latencies) == 0 {
		return
	}
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	n := len(latencies)
	b.ReportMetric(float64(latencies[n/2]), "p50_ns")
	b.ReportMetric(float64(latencies[n*95/100]), "p95_ns")
	b.ReportMetric(float64(latencies[n*99/100]), "p99_ns")
	b.ReportMetric(float64(latencies[n-1]), "p100_ns")
}
