package petacore_test

import "testing"

func BenchmarkSimpleSelect(b *testing.B) {
	conn := NewPgConnection(b)
	defer conn.Close()

	// Prepare a statement
	_, err := conn.Prepare("select_stmt", "SELECT id, value FROM test_table WHERE id = $1;")
	if err != nil {
		b.Fatalf("Failed to prepare statement: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var id int
		var value string

		// Execute the prepared statement
		err = conn.QueryRow("select_stmt", 1).Scan(&id, &value)
		if err != nil {
			b.Fatalf("Failed to execute prepared statement: %v", err)
		}
	}
}
