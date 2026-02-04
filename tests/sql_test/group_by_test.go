package petacore_test

import (
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

func TestOrdersGroupBy(t *testing.T) {

	conn := NewPgConnection(t)
	defer conn.Close()

	PopulateGroupByData(t, conn)

	tests := []struct {
		name     string
		query    string
		expected map[int]float64 // customer_id -> total_amount
	}{
		{
			name:  "Total amount per customer",
			query: "SELECT customer_id, SUM(amount) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 250.0,
				2: 250.0,
				3: 300.0,
			},
		},
		{
			name:  "Count * orders per customer",
			query: "SELECT customer_id, COUNT(*) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 2,
				2: 2,
				3: 1,
			},
		},
		{
			name:  "Count orders per customer",
			query: "SELECT customer_id, COUNT(1) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 2,
				2: 2,
				3: 1,
			},
		},
		{
			name:  "Average order amount per customer",
			query: "SELECT customer_id, AVG(amount) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 125.0,
				2: 125.0,
				3: 300.0,
			},
		},
		{
			name:  "Max order amount per customer",
			query: "SELECT customer_id, MAX(amount) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 150.0,
				2: 200.0,
				3: 300.0,
			},
		},
		{
			name:  "Min order amount per customer",
			query: "SELECT customer_id, MIN(amount) FROM group_by_orders GROUP BY customer_id;",
			expected: map[int]float64{
				1: 100.0,
				2: 50.0,
				3: 300.0,
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

			results := make(map[int]float64)
			for rows.Next() {
				var customerID int
				var totalAmount float64
				err := rows.Scan(&customerID, &totalAmount)
				if err != nil {
					t.Fatalf("Failed to scan row: %v", err)
				}
				results[customerID] = totalAmount
			}

			if len(results) != len(tt.expected) {
				t.Fatalf("Expected %d results, got %d", len(tt.expected), len(results))
			}

			for custID, expectedAmount := range tt.expected {
				if results[custID] != expectedAmount {
					t.Errorf("For customer_id %d, expected total amount %.2f, got %.2f", custID, expectedAmount, results[custID])
				}
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
