package petacore_test

import (
	"testing"

	"github.com/jackc/pgx"
)

func PopulateSubqueryTestData(t testing.TB, conn *pgx.ConnPool) {
	t.Helper()

	tables := []struct {
		name     string
		query    string
		populate string
	}{
		{
			name: "subquery_test_table",
			query: `
				CREATE TABLE IF NOT EXISTS subquery_test_table (
					id INT,
					name TEXT,
					age INT,
					salary FLOAT,
					department TEXT,
					PRIMARY KEY (id)
				);
			`,
			populate: `
				INSERT INTO subquery_test_table (id, name, age, salary, department) VALUES
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

type ExpectedRow struct {
	id         int
	name       string
	age        int
	salary     float64
	department string
}

func TestSubquerySelect(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateSubqueryTestData(t, conn)

	tests := []struct {
		name         string
		query        string
		expectedRows []ExpectedRow
	}{
		{
			name:  "Select with scalar subquery in WHERE",
			query: "SELECT id, name, age, salary, department FROM subquery_test_table WHERE age > (SELECT AVG(age) FROM subquery_test_table) ORDER BY id;",
			expectedRows: []ExpectedRow{
				{3, "Charlie", 35, 80000.0, "Engineering"},
				{5, "Eve", 40, 90000.0, "Management"},
			},
		},
		{
			name:  "Select with IN subquery",
			query: "SELECT id, name, age, salary, department FROM subquery_test_table WHERE id IN (SELECT id FROM subquery_test_table WHERE department = 'Engineering') ORDER BY id;",
			expectedRows: []ExpectedRow{
				{1, "Alice", 30, 70000.0, "Engineering"},
				{3, "Charlie", 35, 80000.0, "Engineering"},
			},
		},
		{
			name:  "Select with EXISTS subquery",
			query: "SELECT id, name, age, salary, department FROM subquery_test_table WHERE EXISTS (SELECT 1 FROM subquery_test_table WHERE subquery_test_table.age > 38 AND subquery_test_table.id = id) ORDER BY id;",
			expectedRows: []ExpectedRow{
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
