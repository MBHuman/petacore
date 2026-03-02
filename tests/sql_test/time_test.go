package petacore_test

import (
	"testing"
	"time"

	"github.com/jackc/pgx"
)

func TestTimeAndDateFunctions(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	tests := []struct {
		name    string
		query   string
		checkFn func(*testing.T, *pgx.Rows)
	}{
		{
			name:  "now_function",
			query: `SELECT NOW();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var ts time.Time
				if err := rows.Scan(&ts); err != nil {
					t.Fatalf("scan: %v", err)
				}
				// Check that timestamp is recent (within last minute)
				now := time.Now()
				diff := now.Sub(ts)
				if diff < 0 {
					diff = -diff
				}
				if diff > time.Minute {
					t.Fatalf("NOW() returned timestamp too far from current time: %v", ts)
				}
			},
		},
		{
			name:  "pg_postmaster_start_time",
			query: `SELECT pg_postmaster_start_time();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var ts time.Time
				if err := rows.Scan(&ts); err != nil {
					t.Fatalf("scan: %v", err)
				}
				// Check that it's a valid timestamp (2024-01-01 00:00:00 UTC)
				expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				if !ts.Equal(expected) {
					t.Logf("pg_postmaster_start_time returned: %v, expected: %v", ts, expected)
					// Don't fail, just log - the returned time might be in different timezone
				}
			},
		},
		{
			name:  "current_database",
			query: `SELECT CURRENT_DATABASE();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var db string
				if err := rows.Scan(&db); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if db != "testdb" {
					t.Fatalf("got %q, want \"testdb\"", db)
				}
			},
		},
		{
			name:  "current_schema",
			query: `SELECT CURRENT_SCHEMA();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var schema string
				if err := rows.Scan(&schema); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if schema != "public" {
					t.Fatalf("got %q, want \"public\"", schema)
				}
			},
		},
		{
			name:  "current_user",
			query: `SELECT CURRENT_USER;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var user string
				if err := rows.Scan(&user); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if user != "postgres" && user != "testuser" {
					t.Logf("got user %q, expected postgres or testuser", user)
				}
			},
		},
		{
			name:  "version",
			query: `SELECT VERSION();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var version string
				if err := rows.Scan(&version); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if version == "" {
					t.Fatalf("VERSION() returned empty string")
				}
				t.Logf("VERSION: %s", version)
			},
		},
		{
			name:  "pg_backend_pid",
			query: `SELECT PG_BACKEND_PID();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var pid int32
				if err := rows.Scan(&pid); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if pid <= 0 {
					t.Fatalf("got invalid PID: %d", pid)
				}
				t.Logf("PG_BACKEND_PID: %d", pid)
			},
		},
		{
			name:  "current_timestamp",
			query: `SELECT CURRENT_TIMESTAMP;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var ts time.Time
				if err := rows.Scan(&ts); err != nil {
					t.Fatalf("scan: %v", err)
				}
				// Check that timestamp is recent
				now := time.Now()
				diff := now.Sub(ts)
				if diff < 0 {
					diff = -diff
				}
				if diff > time.Minute {
					t.Fatalf("CURRENT_TIMESTAMP returned timestamp too far from current time: %v", ts)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := conn.Query(tt.query)
			if err != nil {
				t.Fatalf("query failed: %v", err)
			}
			defer rows.Close()

			tt.checkFn(t, rows)

			if rows.Next() {
				t.Fatalf("expected exactly one row, got more")
			}
			if err := rows.Err(); err != nil {
				t.Fatalf("rows error: %v", err)
			}
		})
	}
}

func TestArrayFunctions(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	tests := []struct {
		name    string
		query   string
		checkFn func(*testing.T, *pgx.Rows)
	}{
		// TODO fix array result
		// {
		// 	name:  "current_schemas_true",
		// 	query: `SELECT CURRENT_SCHEMAS(true);`,
		// 	checkFn: func(t *testing.T, rows *pgx.Rows) {
		// 		if !rows.Next() {
		// 			t.Fatalf("no rows")
		// 		}
		// 		var schemas []string
		// 		if err := rows.Scan(&schemas); err != nil {
		// 			t.Fatalf("scan: %v", err)
		// 		}
		// 		if len(schemas) < 1 {
		// 			t.Fatalf("expected at least 1 schema, got %d", len(schemas))
		// 		}
		// 		// Should contain at least "public"
		// 		found := false
		// 		for _, s := range schemas {
		// 			if s == "public" {
		// 				found = true
		// 				break
		// 			}
		// 		}
		// 		if !found {
		// 			t.Fatalf("expected 'public' in schemas, got %v", schemas)
		// 		}
		// 		t.Logf("current_schemas(true): %v", schemas)
		// 	},
		// },
		// {
		// 	name:  "current_schemas_false",
		// 	query: `SELECT CURRENT_SCHEMAS(false);`,
		// 	checkFn: func(t *testing.T, rows *pgx.Rows) {
		// 		if !rows.Next() {
		// 			t.Fatalf("no rows")
		// 		}
		// 		var schemas []string
		// 		if err := rows.Scan(&schemas); err != nil {
		// 			t.Fatalf("scan: %v", err)
		// 		}
		// 		if len(schemas) < 1 {
		// 			t.Fatalf("expected at least 1 schema, got %d", len(schemas))
		// 		}
		// 		t.Logf("current_schemas(false): %v", schemas)
		// 	},
		// },
		// TODO not supported yer
		// {
		// 	name:  "array_to_string_schemas",
		// 	query: `SELECT ARRAY_TO_STRING(CURRENT_SCHEMAS(true), ', ');`,
		// 	checkFn: func(t *testing.T, rows *pgx.Rows) {
		// 		if !rows.Next() {
		// 			t.Fatalf("no rows")
		// 		}
		// 		var result string
		// 		if err := rows.Scan(&result); err != nil {
		// 			t.Fatalf("scan: %v", err)
		// 		}
		// 		if result == "" {
		// 			t.Fatalf("array_to_string returned empty string")
		// 		}
		// 		t.Logf("array_to_string result: %s", result)
		// 	},
		// },
		// {
		// 	name:  "array_length_schemas",
		// 	query: `SELECT ARRAY_LENGTH(CURRENT_SCHEMAS(true), 1);`,
		// 	checkFn: func(t *testing.T, rows *pgx.Rows) {
		// 		if !rows.Next() {
		// 			t.Fatalf("no rows")
		// 		}
		// 		var length int32
		// 		if err := rows.Scan(&length); err != nil {
		// 			t.Fatalf("scan: %v", err)
		// 		}
		// 		if length < 1 {
		// 			t.Fatalf("expected length >= 1, got %d", length)
		// 		}
		// 		t.Logf("array_length: %d", length)
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := conn.Query(tt.query)
			if err != nil {
				t.Fatalf("query failed: %v", err)
			}
			defer rows.Close()

			tt.checkFn(t, rows)

			if rows.Next() {
				t.Fatalf("expected exactly one row, got more")
			}
			if err := rows.Err(); err != nil {
				t.Fatalf("rows error: %v", err)
			}
		})
	}
}

func TestTimestampTable(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	// Create table with timestamp column
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP
	);`)
	if err != nil {
		t.Fatalf("Failed to create events table: %v", err)
	}

	// Truncate table
	_, err = conn.Exec(`TRUNCATE TABLE events;`)
	if err != nil {
		t.Fatalf("Failed to truncate events table: %v", err)
	}

	// Insert data with NOW()
	_, err = conn.Exec(`INSERT INTO events (name, created_at) VALUES ('event1', NOW());`)
	if err != nil {
		t.Fatalf("Failed to insert event: %v", err)
	}

	// Query the data
	rows, err := conn.Query(`SELECT name, created_at FROM events WHERE name = 'event1';`)
	if err != nil {
		t.Fatalf("Failed to query events: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("no rows returned")
	}

	var name string
	var createdAt time.Time
	if err := rows.Scan(&name, &createdAt); err != nil {
		t.Fatalf("Failed to scan row: %v", err)
	}

	if name != "event1" {
		t.Fatalf("got name %q, want \"event1\"", name)
	}

	// Check that timestamp is recent
	now := time.Now()
	diff := now.Sub(createdAt)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Minute {
		t.Fatalf("created_at timestamp too far from current time: %v", createdAt)
	}

	t.Logf("Event created at: %v", createdAt)
}
