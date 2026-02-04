package petacore_test

import (
	"testing"

	"github.com/jackc/pgx"
)

func PopulateJoinData(t testing.TB, conn *pgx.ConnPool) {
	t.Helper()

	tables := []struct {
		name     string
		query    string
		populate string
	}{
		{
			name: "join_employees",
			query: `
				CREATE TABLE IF NOT EXISTS join_employees (
					id INT,
					name TEXT,
					department_id INT,
					PRIMARY KEY (id)
				);
			`,
			populate: `
				INSERT INTO join_employees (id, name, department_id) VALUES
				(1, 'Alice', 10),
				(2, 'Bob', 20),
				(3, 'Charlie', 10),
				(4, 'David', 30);
			`,
		},
		{
			name: "join_departments",
			query: `
				CREATE TABLE IF NOT EXISTS join_departments (
					id INT,
					department_name TEXT,
					PRIMARY KEY (id)
				);
			`,
			populate: `
				INSERT INTO join_departments (id, department_name) VALUES
				(10, 'HR'),
				(20, 'Engineering'),
				(30, 'Sales');
			`,
		},
		{
			name: "join_salaries",
			query: `
				CREATE TABLE IF NOT EXISTS join_salaries (
					employee_id INT,
					salary FLOAT,
					PRIMARY KEY (employee_id)
				);
			`,
			populate: `
				INSERT INTO join_salaries (employee_id, salary) VALUES
				(1, 60000.0),
				(2, 80000.0),
				(3, 75000.0),
				(4, 50000.0);
			`,
		},
		{
			name: "join_locations",
			query: `
				CREATE TABLE IF NOT EXISTS join_locations (
					department_id INT,
					location TEXT,
					PRIMARY KEY (department_id)
				);
			`,
			populate: `
				INSERT INTO join_locations (department_id, location) VALUES
				(10, 'New York'),
				(20, 'San Francisco'),
				(30, 'Chicago');
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

func TestInnerJoin(t *testing.T) {

	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateJoinData(t, conn)

	tests := []struct {
		name         string
		query        string
		scanFunc     func(rows *pgx.Rows) ([]interface{}, error)
		expectedRows [][]interface{}
	}{
		{
			name: "inner_join_employees_departments",
			query: `
				SELECT e.name, d.department_name
				FROM join_employees e
				INNER JOIN join_departments d ON e.department_id = d.id
				ORDER BY e.name;
			`,
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var name, dept string
				err := rows.Scan(&name, &dept)
				return []interface{}{name, dept}, err
			},
			expectedRows: [][]interface{}{
				{"Alice", "HR"},
				{"Bob", "Engineering"},
				{"Charlie", "HR"},
				{"David", "Sales"},
			},
		},
		{
			name: "inner_join_multiple_tables",
			query: `
				SELECT e.name, d.department_name, s.salary, l.location
				FROM join_employees e
				INNER JOIN join_departments d ON e.department_id = d.id
				INNER JOIN join_salaries s ON e.id = s.employee_id
				INNER JOIN join_locations l ON d.id = l.department_id
				ORDER BY e.name;
			`,
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var name, dept, loc string
				var salary float64
				err := rows.Scan(&name, &dept, &salary, &loc)
				return []interface{}{name, dept, salary, loc}, err
			},
			expectedRows: [][]interface{}{
				{"Alice", "HR", 60000.0, "New York"},
				{"Bob", "Engineering", 80000.0, "San Francisco"},
				{"Charlie", "HR", 75000.0, "New York"},
				{"David", "Sales", 50000.0, "Chicago"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rows, err := conn.Query(tt.query)
			if err != nil {
				t.Fatalf("Query failed: %v", err)
			}
			defer rows.Close()

			var resultRows [][]interface{}
			for rows.Next() {
				row, err := tt.scanFunc(rows)
				if err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				resultRows = append(resultRows, row)
			}

			if len(resultRows) != len(tt.expectedRows) {
				t.Fatalf("Expected %d rows, got %d", len(tt.expectedRows), len(resultRows))
			}
			for i, expectedRow := range tt.expectedRows {
				for j, expectedValue := range expectedRow {
					if resultRows[i][j] != expectedValue {
						t.Errorf("Expected row %d column %d to be %v, got %v", i, j, expectedValue, resultRows[i][j])
					}
				}
			}
		})
	}
}
