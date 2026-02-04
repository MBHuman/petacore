package petacore_test

import (
	"testing"

	"github.com/jackc/pgx"
)

func PopulateSelectTestData(t testing.TB, conn *pgx.ConnPool) {
	t.Helper()

	tables := []struct {
		name     string
		query    string
		populate string
	}{
		{
			name: "select_test_table",
			query: `
				CREATE TABLE IF NOT EXISTS select_test_table (
					id INT,
					name TEXT,
					age INT,
					salary FLOAT,
					department TEXT,
					PRIMARY KEY (id)
				);
			`,
			populate: `
				INSERT INTO select_test_table (id, name, age, salary, department) VALUES
				(1, 'Alice', 30, 70000.0, 'Engineering'),
				(2, 'Bob', 25, 50000.0, 'Marketing'),
				(3, 'Charlie', 35, 80000.0, 'Engineering'),
				(4, 'Diana', 28, 60000.0, 'Sales'),
				(5, 'Eve', 40, 90000.0, 'Management');
			`,
		},
	}

	for _, table := range tables {
		_, err := conn.Exec(table.query)
		if err != nil {
			t.Fatalf("failed to create table %s: %v", table.name, err)
		}
		_, err = conn.Exec("TRUNCATE TABLE " + table.name)
		if err != nil {
			t.Fatalf("failed to truncate table %s: %v", table.name, err)
		}
		_, err = conn.Exec(table.populate)
		if err != nil {
			t.Fatalf("failed to populate table %s: %v", table.name, err)
		}
	}
}

func TestSelect(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateSelectTestData(t, conn)

	tests := []struct {
		name         string
		query        string
		scanFunc     func(rows *pgx.Rows) ([]interface{}, error)
		expectedRows [][]interface{}
	}{
		{
			name:  "Select all columns",
			query: "SELECT id, name, age, salary, department FROM select_test_table ORDER BY id;",
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var results []interface{}
				for rows.Next() {
					var id int
					var name string
					var age int
					var salary float64
					var department string
					err := rows.Scan(&id, &name, &age, &salary, &department)
					if err != nil {
						return nil, err
					}
					results = append(results, []interface{}{id, name, age, salary, department})
				}
				return results, nil
			},
			expectedRows: [][]interface{}{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{2, "Bob", 25, 50000.0, "Marketing"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{4, "Diana", 28, 60000.0, "Sales"},
				{5, "Eve", 40, 90000.0, "Management"},
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

			results, err := tt.scanFunc(rows)
			if err != nil {
				t.Fatalf("failed to scan rows: %v", err)
			}

			if len(results) != len(tt.expectedRows) {
				t.Fatalf("expected %d rows, got %d", len(tt.expectedRows), len(results))
			}
			for i, expectedRow := range tt.expectedRows {
				resultRow := results[i].([]interface{})
				for j, expectedValue := range expectedRow {
					if resultRow[j] != expectedValue {
						t.Errorf("row %d column %d: expected %v, got %v", i+1, j+1, expectedValue, resultRow[j])
					}
				}
			}
		})
	}
}

func TestSelectWhereClause(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateSelectTestData(t, conn)

	type ExpectedRow struct {
		id         int
		name       string
		age        int
		salary     float64
		department string
	}

	tests := []struct {
		name         string
		query        string
		expectedRows []ExpectedRow
	}{
		{
			name:  "Select with WHERE clause",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age > 30 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with WHERE clause and string comparison",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE department = 'Engineering' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
			},
		},
		{
			name:  "Select with IN operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE department IN ('Engineering', 'Sales') ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{4, "Diana", 28, 60000.0, "Sales"},
			},
		},
		{
			name:  "Select with LIKE operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE name LIKE 'A%' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
			},
		},
		{
			name:  "Select with multiple conditions",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age < 35 AND salary >= 60000.0 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{4, "Diana", 28, 60000.0, "Sales"},
			},
		},
		{
			name:  "Select with OR operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE department = 'Marketing' OR department = 'Sales' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{2, "Bob", 25, 50000.0, "Marketing"},
				{4, "Diana", 28, 60000.0, "Sales"},
			},
		},
		{
			name:  "Select with combined AND/OR conditions",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE (department = 'Engineering' AND age > 32) OR (department = 'Management') ORDER BY id;",
			expectedRows: []ExpectedRow{
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with < operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age < 30 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{2, "Bob", 25, 50000.0, "Marketing"},
				{4, "Diana", 28, 60000.0, "Sales"},
			},
		},
		{
			name:  "Select with <= operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age <= 28 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{2, "Bob", 25, 50000.0, "Marketing"},
				{4, "Diana", 28, 60000.0, "Sales"},
			},
		},
		{
			name:  "Select with > operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age > 30 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with >= operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE age >= 30 ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:         "Select with IS NULL operator",
			query:        "SELECT id, name, age, salary, department FROM select_test_table WHERE department IS NULL ORDER BY id;",
			expectedRows: []ExpectedRow{
				// No rows expected since all departments are non-null
			},
		},
		{
			name:  "Select with NOT operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE NOT department = 'Marketing' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{4, "Diana", 28, 60000.0, "Sales"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with != operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE department != 'Engineering' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{2, "Bob", 25, 50000.0, "Marketing"},
				{4, "Diana", 28, 60000.0, "Sales"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with <> operator",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE department <> 'Engineering' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{2, "Bob", 25, 50000.0, "Marketing"},
				{4, "Diana", 28, 60000.0, "Sales"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with regex operator ~",
			query: "SELECT id, name, age, salary, department FROM select_test_table WHERE name ~ '^A' ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
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

			var results []ExpectedRow
			for rows.Next() {
				var r ExpectedRow
				err := rows.Scan(&r.id, &r.name, &r.age, &r.salary, &r.department)
				if err != nil {
					t.Fatalf("failed to scan row: %v", err)
				}
				results = append(results, r)
			}

			if len(results) != len(tt.expectedRows) {
				t.Fatalf("expected %d rows, got %d", len(tt.expectedRows), len(results))
			}
			for i, expectedRow := range tt.expectedRows {
				resultRow := results[i]
				if resultRow != expectedRow {
					t.Errorf("row %d: expected %+v, got %+v", i+1, expectedRow, resultRow)
				}
			}
		})
	}

}
