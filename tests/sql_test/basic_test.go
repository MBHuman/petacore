package petacore_test

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx"
)

func NewPgConnection(t testing.TB) *pgx.ConnPool {
	t.Helper()

	conn, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Database: "testdb",
			},
			MaxConnections: 1,
		},
	)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Create table

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS test_table (
		id SERIAL PRIMARY KEY,
		value TEXT
	);`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	_, err = conn.Exec(`TRUNCATE TABLE test_table;`)
	if err != nil {
		t.Fatalf("Failed to truncate table: %v", err)
	}
	_, err = conn.Exec(`INSERT INTO test_table (value) VALUES ('test1'), ('test2'), ('test3');`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return conn
}

func TestPreparedStatements(t *testing.T) {
	t.Skip()
	conn := NewPgConnection(t)
	defer conn.Close()

	// Prepare a statement
	_, err := conn.Prepare("select_stmt", "SELECT id, value FROM test_table WHERE id = $1;")
	if err != nil {
		t.Fatalf("Failed to prepare statement: %v", err)
	}

	// Execute the prepared statement
	var id int
	var value string
	err = conn.QueryRow("select_stmt", 2).Scan(&id, &value)
	if err != nil {
		t.Fatalf("Failed to execute prepared statement: %v", err)
	}

	// Verify the result
	if id != 2 || value != "test2" {
		t.Fatalf("Unexpected result from prepared statement: got (%d, %s), want (2, test2)", id, value)
	}
}

func TestRegexSelect(t *testing.T) {
	// t.Skip()
	conn := NewPgConnection(t)
	defer conn.Close()

	// Test regex SELECT query
	rows, err := conn.Query(`
		SELECT value
		FROM test_table
		WHERE value OPERATOR(pg_catalog.~) '^test[12]$' COLLATE pg_catalog.default
		ORDER BY value;
	`)
	if err != nil {
		t.Fatalf("Failed to execute regex SELECT query: %v", err)
	}
	defer rows.Close()

	expectedValues := []string{"test1", "test2"}
	i := 0
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			t.Fatalf("Failed to scan row %d: %v", i, err)
		}

		if i >= len(expectedValues) {
			t.Fatalf("Too many rows returned")
		}

		if value != expectedValues[i] {
			t.Fatalf("Row %d: got %s, want %s", i, value, expectedValues[i])
		}
		i++
	}

	if i != len(expectedValues) {
		t.Fatalf("Expected %d rows, got %d", len(expectedValues), i)
	}
}

// Test simple SELECT query

func TestSelectWithCaseWhen(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	// Test CASE WHEN expression
	rows, err := conn.Query(`
		SELECT id, value,
			CASE
				WHEN id = 1 THEN 'first'
				WHEN id = 2 THEN 'second'
				ELSE 'other'
			END as position
		FROM test_table
		ORDER BY id
	`)
	if err != nil {
		t.Fatalf("Failed to execute CASE WHEN query: %v", err)
	}
	defer rows.Close()

	expected := []struct {
		id       int
		value    string
		position string
	}{
		{1, "test1", "first"},
		{2, "test2", "second"},
		{3, "test3", "other"},
	}

	i := 0
	for rows.Next() {
		var id int
		var value, position string
		err := rows.Scan(&id, &value, &position)
		if err != nil {
			t.Fatalf("Failed to scan row %d: %v", i, err)
		}

		if i >= len(expected) {
			t.Fatalf("Too many rows returned")
		}

		if id != expected[i].id || value != expected[i].value || position != expected[i].position {
			t.Fatalf("Row %d: got (%d, %s, %s), want (%d, %s, %s)",
				i, id, value, position, expected[i].id, expected[i].value, expected[i].position)
		}
		i++
	}

	if i != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), i)
	}
}

func TestGroupBySelect(t *testing.T) {
	// t.Skip()
	conn := NewPgConnection(t)
	defer conn.Close()

	// Test GROUP BY query
	rows, err := conn.Query(`
		SELECT value, COUNT(1) as count
		FROM test_table
		GROUP BY value
		ORDER BY value;
	`)
	if err != nil {
		t.Fatalf("Failed to execute GROUP BY query: %v", err)
	}

	defer rows.Close()

	expected := map[string]int{
		"test1": 1,
		"test2": 1,
		"test3": 1,
	}

	i := 0
	for rows.Next() {
		var value string
		var count int
		err := rows.Scan(&value, &count)
		if err != nil {
			t.Fatalf("Failed to scan row %d: %v", i, err)
		}

		expectedCount, ok := expected[value]
		if !ok {
			t.Fatalf("Unexpected value %s in result", value)
		}

		if count != expectedCount {
			t.Fatalf("Value %s: got count %d, want %d", value, count, expectedCount)
		}
		i++
	}

	if i != len(expected) {
		t.Fatalf("Expected %d rows, got %d", len(expected), i)
	}
}

func TestAggregateSelect(t *testing.T) {
	// t.Skip()
	conn := NewPgConnection(t)
	defer conn.Close()

	cases := []struct {
		name     string
		query    string
		expected interface{}
	}{
		{
			name:     "COUNT(1)",
			query:    "SELECT COUNT(1) FROM test_table",
			expected: 3,
		},
		{
			name:     "COUNT(*)",
			query:    "SELECT COUNT(*) FROM test_table",
			expected: 3,
		},
		{
			name:     "MAX(id)",
			query:    "SELECT MAX(id) FROM test_table",
			expected: 3,
		},
		{
			name:     "MIN(id)",
			query:    "SELECT MIN(id) FROM test_table",
			expected: 1,
		},
		{
			name:     "SUM(id)",
			query:    "SELECT SUM(id) FROM test_table",
			expected: 6,
		},
		{
			name:     "AVG(id)",
			query:    "SELECT AVG(id) FROM test_table",
			expected: 2.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			row := conn.QueryRow(tc.query)
			var result interface{}
			switch tc.expected.(type) {
			case int:
				var val int
				err := row.Scan(&val)
				if err != nil {
					t.Fatalf("Failed to scan %s: %v", tc.name, err)
				}
				result = val
			case float64:
				var val float64
				err := row.Scan(&val)
				if err != nil {
					t.Fatalf("Failed to scan %s: %v", tc.name, err)
				}
				result = val
			default:
				t.Fatalf("Unsupported expected type for %s", tc.name)
			}
			if result != tc.expected {
				t.Errorf("%s = %v, want %v", tc.name, result, tc.expected)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
