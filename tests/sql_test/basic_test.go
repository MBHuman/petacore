package petacore_test

import (
	"testing"

	_ "github.com/lib/pq"

	"github.com/jackc/pgx"
)

func NewPgConnection(t testing.TB) *pgx.ConnPool {
	t.Helper()

	// connConfig := pgx.ConnConfig{
	// 	Host:     "localhost",
	// 	Port:     5432,
	// 	User:     "testuser",
	// 	Database: "testdb",
	// }
	conn, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgx.ConnConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				Database: "testdb",
			},
			MaxConnections: 8,
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

// func TestSelectWithCaseWhen(t *testing.T) {
// 	conn := NewPgConnection(t)
// 	defer conn.Close()

// 	// Test CASE WHEN expression
// 	rows, err := conn.Query(`
// 		SELECT id, value,
// 			CASE
// 				WHEN id = 1 THEN 'first'
// 				WHEN id = 2 THEN 'second'
// 				ELSE 'other'
// 			END as position
// 		FROM test_table
// 		ORDER BY id
// 	`)
// 	if err != nil {
// 		t.Fatalf("Failed to execute CASE WHEN query: %v", err)
// 	}
// 	defer rows.Close()

// 	expected := []struct {
// 		id       int
// 		value    string
// 		position string
// 	}{
// 		{1, "test1", "first"},
// 		{2, "test2", "second"},
// 		{3, "test3", "other"},
// 	}

// 	i := 0
// 	for rows.Next() {
// 		var id int
// 		var value, position string
// 		err := rows.Scan(&id, &value, &position)
// 		if err != nil {
// 			t.Fatalf("Failed to scan row %d: %v", i, err)
// 		}

// 		if i >= len(expected) {
// 			t.Fatalf("Too many rows returned")
// 		}

// 		if id != expected[i].id || value != expected[i].value || position != expected[i].position {
// 			t.Fatalf("Row %d: got (%d, %s, %s), want (%d, %s, %s)",
// 				i, id, value, position, expected[i].id, expected[i].value, expected[i].position)
// 		}
// 		i++
// 	}

// 	if i != len(expected) {
// 		t.Fatalf("Expected %d rows, got %d", len(expected), i)
// 	}
// }

// func TestSelectWithJoin(t *testing.T) {
// 	conn := NewPgConnection(t)
// 	defer conn.Close()

// 	// Create a second table for JOIN testing
// 	_, err := conn.Exec(`
// 		CREATE TABLE IF NOT EXISTS categories (
// 			id SERIAL PRIMARY KEY,
// 			name TEXT,
// 			parent_id INTEGER
// 		)
// 	`)
// 	if err != nil {
// 		t.Fatalf("Failed to create categories table: %v", err)
// 	}

// 	_, err = conn.Exec(`TRUNCATE TABLE categories;`)
// 	if err != nil {
// 		t.Fatalf("Failed to truncate categories table: %v", err)
// 	}

// 	_, err = conn.Exec(`
// 		INSERT INTO categories (name, parent_id) VALUES
// 		('Electronics', NULL),
// 		('Phones', 1),
// 		('Laptops', 1)
// 	`)
// 	if err != nil {
// 		t.Fatalf("Failed to insert categories: %v", err)
// 	}

// 	// Test JOIN query
// 	rows, err := conn.Query(`
// 		SELECT c.id, c.name, p.name as parent_name
// 		FROM categories c
// 		LEFT JOIN categories p ON c.parent_id = p.id
// 		ORDER BY c.id
// 	`)
// 	if err != nil {
// 		t.Fatalf("Failed to execute JOIN query: %v", err)
// 	}
// 	defer rows.Close()

// 	expected := []struct {
// 		id          int
// 		name        string
// 		parent_name *string
// 	}{
// 		{1, "Electronics", nil},
// 		{2, "Phones", stringPtr("Electronics")},
// 		{3, "Laptops", stringPtr("Electronics")},
// 	}

// 	i := 0
// 	for rows.Next() {
// 		var id int
// 		var name string
// 		var parentName *string
// 		err := rows.Scan(&id, &name, &parentName)
// 		if err != nil {
// 			t.Fatalf("Failed to scan row %d: %v", i, err)
// 		}

// 		if i >= len(expected) {
// 			t.Fatalf("Too many rows returned")
// 		}

// 		if id != expected[i].id || name != expected[i].name {
// 			t.Fatalf("Row %d: got (%d, %s), want (%d, %s)",
// 				i, id, name, expected[i].id, expected[i].name)
// 		}

// 		if (expected[i].parent_name == nil && parentName != nil) ||
// 			(expected[i].parent_name != nil && parentName == nil) ||
// 			(expected[i].parent_name != nil && parentName != nil && *expected[i].parent_name != *parentName) {
// 			t.Fatalf("Row %d parent_name: got %v, want %v", i, parentName, expected[i].parent_name)
// 		}
// 		i++
// 	}

// 	if i != len(expected) {
// 		t.Fatalf("Expected %d rows, got %d", len(expected), i)
// 	}
// }

// func TestSelectWithSubquery(t *testing.T) {
// 	conn := NewPgConnection(t)
// 	defer conn.Close()

// 	// Test subquery in WHERE clause
// 	var id int
// 	var value string
// 	err := conn.QueryRow(`
// 		SELECT id, value
// 		FROM test_table
// 		WHERE id = (SELECT MAX(id) FROM test_table)
// 	`).Scan(&id, &value)
// 	if err != nil {
// 		t.Fatalf("Failed to execute subquery: %v", err)
// 	}

// 	if id != 3 || value != "test3" {
// 		t.Fatalf("Subquery result: got (%d, %s), want (3, test3)", id, value)
// 	}
// }

func stringPtr(s string) *string {
	return &s
}
