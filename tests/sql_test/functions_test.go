package petacore_test

import (
	"testing"

	"github.com/jackc/pgx"
)

func PopulateFunctionsByData(t testing.TB, conn *pgx.ConnPool) {
	t.Helper()

	tables := []struct {
		name     string
		query    string
		populate string
	}{
		{
			name: "numbers",
			query: `
				CREATE TABLE IF NOT EXISTS numbers (
					id SERIAL PRIMARY KEY,
					value INTEGER NOT NULL
				);
			`,
			populate: `
				INSERT INTO numbers (value) VALUES
				(1), (2), (3), (4), (5);
			`,
		},
		{
			name: "strings",
			query: `
				CREATE TABLE IF NOT EXISTS strings (
					id SERIAL PRIMARY KEY,
					value TEXT NOT NULL
				);
			`,
			populate: `
				INSERT INTO strings (value) VALUES
				('foo'), ('bar'), ('baz'), ('qux'), ('quux');
			`,
		},
	}

	for _, table := range tables {
		_, err := conn.Exec(table.query)
		if err != nil {
			t.Fatalf("Failed to create table %s: %v", table.name, err)
		}

		_, err = conn.Exec("TRUNCATE TABLE " + table.name + ";")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table.name, err)
		}

		_, err = conn.Exec(table.populate)
		if err != nil {
			t.Fatalf("Failed to populate table %s: %v", table.name, err)
		}
	}
}

func TestFunctions(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateFunctionsByData(t, conn)

	tests := []struct {
		name    string
		query   string
		checkFn func(*testing.T, *pgx.Rows)
	}{
		{
			name:  "upper",
			query: `SELECT UPPER('hello');`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "HELLO" {
					t.Fatalf("got %v, want HELLO", s)
				}
			},
		},
		{
			name:  "count_numbers",
			query: `SELECT COUNT(value) FROM numbers;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var c int
				if err := rows.Scan(&c); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if c != 5 {
					t.Fatalf("got %v, want 5", c)
				}
			},
		},
		{
			name:  "sum_numbers",
			query: `SELECT SUM(value) FROM numbers;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s int
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != 15 {
					t.Fatalf("got %v, want 15", s)
				}
			},
		},
		{
			name:  "avg_numbers",
			query: `SELECT AVG(value) FROM numbers;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var a float64
				if err := rows.Scan(&a); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if a != 3.0 {
					t.Fatalf("got %v, want 3.0", a)
				}
			},
		},
		{
			name:  "max_numbers",
			query: `SELECT MAX(value) FROM numbers;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var m int
				if err := rows.Scan(&m); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if m != 5 {
					t.Fatalf("got %v, want 5", m)
				}
			},
		},
		{
			name:  "min_numbers",
			query: `SELECT MIN(value) FROM numbers;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var m int
				if err := rows.Scan(&m); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if m != 1 {
					t.Fatalf("got %v, want 1", m)
				}
			},
		},
		{
			name:  "max_text",
			query: `SELECT MAX(value) FROM strings;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "qux" {
					t.Fatalf("got %v, want qux", s)
				}
			},
		},
		{
			name:  "min_text",
			query: `SELECT MIN(value) FROM strings;`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "bar" {
					t.Fatalf("got %v, want bar", s)
				}
			},
		},
		{
			name:  "quote_ident",
			query: `SELECT QUOTE_IDENT('My.Id');`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if len(s) == 0 || s[0] != '"' {
					t.Fatalf("got %v, want quoted string", s)
				}
			},
		},
		{
			name:  "length",
			query: `SELECT LENGTH('привет');`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var l int
				if err := rows.Scan(&l); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if l != 6 {
					t.Fatalf("got %v, want 6", l)
				}
			},
		},
		{
			name:  "substring",
			query: `SELECT SUBSTRING('abcdef', 1, 3);`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "abc" {
					t.Fatalf("got %v, want abc", s)
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
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "testdb" {
					t.Fatalf("got %v, want testdb", s)
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
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if s != "public" {
					t.Fatalf("got %v, want public", s)
				}
			},
		},
		// todo remove from antlr grammar
		// {
		// 	name:  "current_user",
		// 	query: `SELECT CURRENT_USER();`,
		// 	checkFn: func(t *testing.T, rows *pgx.Rows) {
		// 		if !rows.Next() {
		// 			t.Fatalf("no rows")
		// 		}
		// 		var s string
		// 		if err := rows.Scan(&s); err != nil {
		// 			t.Fatalf("scan: %v", err)
		// 		}
		// 		if s != "postgres" {
		// 			t.Fatalf("got %v, want postgres", s)
		// 		}
		// 	},
		// },
		{
			name:  "version",
			query: `SELECT VERSION();`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var s string
				if err := rows.Scan(&s); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if len(s) == 0 {
					t.Fatalf("version empty")
				}
			},
		},
		{
			name:  "pg_table_is_visible",
			query: `SELECT PG_TABLE_IS_VISIBLE(1);`,
			checkFn: func(t *testing.T, rows *pgx.Rows) {
				if !rows.Next() {
					t.Fatalf("no rows")
				}
				var b bool
				if err := rows.Scan(&b); err != nil {
					t.Fatalf("scan: %v", err)
				}
				if b != true {
					t.Fatalf("got %v, want true", b)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := conn.Query(tt.query)
			if err != nil {
				t.Fatalf("Failed to execute query: %v", err)
			}
			defer rows.Close()

			tt.checkFn(t, rows)
		})
	}
}
