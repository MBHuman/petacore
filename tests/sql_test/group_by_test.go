package petacore_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/jackc/pgx"
)

func PopulateGroupByData(t testing.TB, conn *pgx.ConnPool) {
	t.Helper()

	tables := []struct {
		name     string
		query    string
		populate string
	}{
		{
			name: "group_by_orders",
			query: `
				CREATE TABLE IF NOT EXISTS group_by_orders (
					id INT,
					customer_id INT,
					amount FLOAT,
					PRIMARY KEY (id)
				);
			`,
			populate: `
				INSERT INTO group_by_orders (id, customer_id, amount) VALUES
				(1, 1, 100.0),
				(2, 1, 150.0),
				(3, 2, 200.0),
				(4, 2, 50.0),
				(5, 3, 300.0);
			`,
		},
		{
			name: "group_by_customers",
			query: `
				CREATE TABLE IF NOT EXISTS group_by_customers (
					id INT,
					name TEXT,
					PRIMARY KEY (id)
				);
				`,
			populate: `
				INSERT INTO group_by_customers (id, name) VALUES
				(1, 'Alice'),
				(2, 'Bob'),
				(3, 'Charlie');
			`,
		},
		{
			name: "group_by_products",
			query: `
				CREATE TABLE IF NOT EXISTS group_by_products (
					id INT,
					category_id INT,
					customer_id INT,
					price FLOAT,
					PRIMARY KEY (id, category_id, customer_id)
				);
			`,
			populate: `
				INSERT INTO group_by_products (id, category_id, customer_id, price) VALUES
				(1, 10, 1, 25.0),
				(1, 10, 2, 30.0),
				(2, 10, 2, 30.0),
				(3, 20, 1, 45.0),
				(4, 20, 3, 50.0);	
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

func TestAggregateFunctions(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name    string
		query   string
		checkFn func(t *testing.T, row *pgx.Row) error
	}{
		{
			name:  "SUM test",
			query: "SELECT SUM(amount) FROM group_by_orders;",
			checkFn: func(t *testing.T, row *pgx.Row) error {
				var sum float64
				err := row.Scan(&sum)
				if err != nil {
					return fmt.Errorf("failed to scan result: %v", err)
				}
				if sum != 800.0 {
					return fmt.Errorf("expected sum 800.0, got %.2f", sum)
				}
				return nil
			},
		},
		{
			name:  "COUNT test",
			query: "SELECT COUNT(*) FROM group_by_orders;",
			// expected: 5,
			checkFn: func(t *testing.T, row *pgx.Row) error {
				var cnt int
				err := row.Scan(&cnt)
				if err != nil {
					return fmt.Errorf("failed to scan result: %v", err)
				}
				if cnt != 5 {
					return fmt.Errorf("expected count 5, got %d", cnt)
				}
				return nil
			},
		},
		{
			name:  "AVG test",
			query: "SELECT AVG(amount) FROM group_by_orders;",
			checkFn: func(t *testing.T, row *pgx.Row) error {
				var avg float64
				err := row.Scan(&avg)
				if err != nil {
					return fmt.Errorf("failed to scan result: %v", err)
				}
				if avg != 160.0 {
					return fmt.Errorf("expected avg 160.0, got %.2f", avg)
				}
				return nil
			},
		},
		{
			name:  "MAX test",
			query: "SELECT MAX(amount) FROM group_by_orders;",
			checkFn: func(t *testing.T, row *pgx.Row) error {
				var max float64
				err := row.Scan(&max)
				if err != nil {
					return fmt.Errorf("failed to scan result: %v", err)
				}
				if max != 300.0 {
					return fmt.Errorf("expected max 300.0, got %.2f", max)
				}
				return nil
			},
		},
		{
			name:  "MIN test",
			query: "SELECT MIN(amount) FROM group_by_orders;",
			checkFn: func(t *testing.T, row *pgx.Row) error {
				var min float64
				err := row.Scan(&min)
				if err != nil {
					return fmt.Errorf("failed to scan result: %v", err)
				}
				if min != 50.0 {
					return fmt.Errorf("expected min 50.0, got %.2f", min)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := conn.QueryRow(tt.query)
			if tt.checkFn != nil {
				if err := tt.checkFn(t, row); err != nil {
					t.Errorf("Check failed: %v", err)
				}
			}
		})
	}
}

func TestMultipleAggregates(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name     string
		query    string
		scanFunc func(rows *pgx.Rows) ([]interface{}, error)
		expected [][]interface{}
	}{
		{
			name:  "SUM and COUNT test",
			query: "SELECT SUM(amount), COUNT(*) FROM group_by_orders;",
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var sum float64
				var count int
				err := rows.Scan(&sum, &count)
				if err != nil {
					return nil, err
				}
				return []interface{}{sum, count}, nil
			},
			expected: [][]interface{}{
				{800.0, 5},
			},
		},
		{
			name:  "AVG, MAX and MIN test",
			query: "SELECT AVG(amount), MAX(amount), MIN(amount) FROM group_by_orders;",
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var avg, max, min float64
				err := rows.Scan(&avg, &max, &min)
				if err != nil {
					return nil, err
				}
				return []interface{}{avg, max, min}, nil
			},
			expected: [][]interface{}{
				{160.0, 300.0, 50.0},
			},
		},
		{
			name:  "Multiple SUM test",
			query: "SELECT SUM(amount) AS a1, SUM(amount + 1) AS a2, SUM(amount + 2) AS a3, SUM(amount + 3) AS a4 FROM group_by_orders;",
			scanFunc: func(rows *pgx.Rows) ([]interface{}, error) {
				var a1, a2, a3, a4 float64
				err := rows.Scan(&a1, &a2, &a3, &a4)
				if err != nil {
					return nil, err
				}
				return []interface{}{a1, a2, a3, a4}, nil
			},
			expected: [][]interface{}{
				{800.0, 805.0, 810.0, 815.0},
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

			// проверяем ошибку после итерации
			if err := rows.Err(); err != nil {
				t.Fatalf("Rows iteration error: %v", err)
			}

			if len(resultRows) != len(tt.expected) {
				t.Fatalf("Expected %d rows, got %d", len(tt.expected), len(resultRows))
			}

			for i, expectedRow := range tt.expected {
				resultRow := resultRows[i]
				for j, expectedValue := range expectedRow {
					// float сравнение через epsilon
					switch exp := expectedValue.(type) {
					case float64:
						got, ok := resultRow[j].(float64)
						if !ok {
							t.Errorf("row %d col %d: expected float64, got %T", i+1, j+1, resultRow[j])
							continue
						}
						if math.Abs(got-exp) > 1e-9 {
							t.Errorf("row %d col %d: expected %v, got %v", i+1, j+1, exp, got)
						}
					default:
						if resultRow[j] != expectedValue {
							t.Errorf("row %d col %d: expected %v, got %v", i+1, j+1, expectedValue, resultRow[j])
						}
					}
				}
			}
		})
	}
}

func TestOrdersGroupBy(t *testing.T) {

	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name    string
		query   string
		checkFn func(t *testing.T, rows *pgx.Rows) error
	}{
		{
			name:  "Total amount per customer",
			query: "SELECT customer_id, SUM(amount) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]float64{
					1: 250.0,
					2: 250.0,
					3: 300.0,
				}
				results := make(map[int]float64)
				for rows.Next() {
					var customerID int
					var totalAmount float64
					err := rows.Scan(&customerID, &totalAmount)
					if err != nil {
						return err
					}
					results[customerID] = totalAmount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedAmount := range expected {
					if results[custID] != expectedAmount {
						return fmt.Errorf("for customer_id %d, expected total amount %.2f, got %.2f", custID, expectedAmount, results[custID])
					}
				}
				return nil
			},
		},
		{
			name:  "Count * orders per customer",
			query: "SELECT customer_id, COUNT(*) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]int{
					1: 2,
					2: 2,
					3: 1,
				}
				results := make(map[int]int)
				for rows.Next() {
					var customerID int
					var orderCount int
					err := rows.Scan(&customerID, &orderCount)
					if err != nil {
						return err
					}
					results[customerID] = orderCount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedCount := range expected {
					if results[custID] != expectedCount {
						return fmt.Errorf("for customer_id %d, expected order count %d, got %d", custID, expectedCount, results[custID])
					}
				}
				return nil
			},
		},
		{
			name:  "Count orders per customer",
			query: "SELECT customer_id, COUNT(1) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]int{
					1: 2,
					2: 2,
					3: 1,
				}
				results := make(map[int]int)
				for rows.Next() {
					var customerID int
					var orderCount int
					err := rows.Scan(&customerID, &orderCount)
					if err != nil {
						return err
					}
					results[customerID] = orderCount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedCount := range expected {
					if results[custID] != expectedCount {
						return fmt.Errorf("for customer_id %d, expected order count %d, got %d", custID, expectedCount, results[custID])
					}
				}
				return nil
			},
		},
		{
			name:  "Average order amount per customer",
			query: "SELECT customer_id, AVG(amount) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]float64{
					1: 125.0,
					2: 125.0,
					3: 300.0,
				}
				results := make(map[int]float64)
				for rows.Next() {
					var customerID int
					var avgAmount float64
					err := rows.Scan(&customerID, &avgAmount)
					if err != nil {
						return err
					}
					results[customerID] = avgAmount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedAmount := range expected {
					if results[custID] != expectedAmount {
						return fmt.Errorf("for customer_id %d, expected average amount %.2f, got %.2f", custID, expectedAmount, results[custID])
					}
				}
				return nil
			},
		},
		{
			name:  "Max order amount per customer",
			query: "SELECT customer_id, MAX(amount) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]float64{
					1: 150.0,
					2: 200.0,
					3: 300.0,
				}
				results := make(map[int]float64)
				for rows.Next() {
					var customerID int
					var maxAmount float64
					err := rows.Scan(&customerID, &maxAmount)
					if err != nil {
						return err
					}
					results[customerID] = maxAmount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedAmount := range expected {
					if results[custID] != expectedAmount {
						return fmt.Errorf("for customer_id %d, expected max amount %.2f, got %.2f", custID, expectedAmount, results[custID])
					}
				}
				return nil
			},
		},
		{
			name:  "Min order amount per customer",
			query: "SELECT customer_id, MIN(amount) FROM group_by_orders GROUP BY customer_id;",
			checkFn: func(t *testing.T, rows *pgx.Rows) error {
				expected := map[int]float64{
					1: 100.0,
					2: 50.0,
					3: 300.0,
				}
				results := make(map[int]float64)
				for rows.Next() {
					var customerID int
					var minAmount float64
					err := rows.Scan(&customerID, &minAmount)
					if err != nil {
						return err
					}
					results[customerID] = minAmount
				}

				if len(results) != len(expected) {
					return fmt.Errorf("expected %d results, got %d", len(expected), len(results))
				}
				for custID, expectedAmount := range expected {
					if results[custID] != expectedAmount {
						return fmt.Errorf("for customer_id %d, expected min amount %.2f, got %.2f", custID, expectedAmount, results[custID])
					}
				}
				return nil
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
			if tt.checkFn != nil {
				if err := tt.checkFn(t, rows); err != nil {
					t.Errorf("Check failed: %v", err)
				}
			}
			// If no checkFn provided, just ensure the query executes without error.
			// Detailed result validation is handled in other tests.
			if err := rows.Err(); err != nil {
				t.Fatalf("Error iterating rows: %v", err)
			}
		})
	}
}

func TestOrdersCombinedGroupBy(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name     string
		query    string
		expected map[int]map[int]float64 // customer_id -> category_id -> total_price
	}{
		{
			name: "Total price per customer and category",
			query: `
				SELECT customer_id, category_id, SUM(price)
				FROM group_by_products
				GROUP BY customer_id, category_id;
			`,
			expected: map[int]map[int]float64{
				1: {
					10: 25.0,
					20: 45.0,
				},
				2: {
					10: 60.0,
				},
				3: {
					20: 50.0,
				},
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

			results := make(map[int]map[int]float64)
			for rows.Next() {
				var customerID, categoryID int
				var totalPrice float64
				err := rows.Scan(&customerID, &categoryID, &totalPrice)
				if err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				if _, exists := results[customerID]; !exists {
					results[customerID] = make(map[int]float64)
				}
				results[customerID][categoryID] = totalPrice
			}

			if len(results) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(results))
			}

			for custID, expectedCategories := range tt.expected {
				actualCategories, exists := results[custID]
				if !exists {
					t.Errorf("Missing results for customer_id %d", custID)
					continue
				}
				if len(actualCategories) != len(expectedCategories) {
					t.Errorf("For customer_id %d, expected %d categories, got %d", custID, len(expectedCategories), len(actualCategories))
					continue
				}
				for catID, expectedPrice := range expectedCategories {
					if actualCategories[catID] != expectedPrice {
						t.Errorf("For customer_id %d and category_id %d, expected total price %.2f, got %.2f", custID, catID, expectedPrice, actualCategories[catID])
					}
				}
			}
		})
	}
}

func TestJoinGroupBy(t *testing.T) {
	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name     string
		query    string
		expected map[string]float64 // customer_name -> total_amount
	}{
		{
			name: "Total amount per customer name",
			query: `
				SELECT c.name, SUM(o.amount)
				FROM group_by_orders o
				JOIN group_by_customers c ON o.customer_id = c.id
				GROUP BY c.name;
			`,
			expected: map[string]float64{
				"Alice":   250.0,
				"Bob":     250.0,
				"Charlie": 300.0,
			},
		},
		{
			name: "Total amount per customer name 2",
			query: `
				SELECT name, SUM(o.amount)
				FROM group_by_customers c
				JOIN group_by_orders o ON o.customer_id = c.id
				GROUP BY name;
			`,
			expected: map[string]float64{
				"Alice":   250.0,
				"Bob":     250.0,
				"Charlie": 300.0,
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

			results := make(map[string]float64)
			for rows.Next() {
				var customerName string
				var totalAmount float64
				err := rows.Scan(&customerName, &totalAmount)
				if err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				results[customerName] = totalAmount
			}

			if len(results) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(results))
			}

			for custName, expectedAmount := range tt.expected {
				if results[custName] != expectedAmount {
					t.Errorf("For customer name %s, expected total amount %.2f, got %.2f", custName, expectedAmount, results[custName])
				}
			}
		})
	}
}
