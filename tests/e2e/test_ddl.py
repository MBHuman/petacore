import pytest
import psycopg2
from fixtures.tables import db_connection, setup_tables

class TestDDLE2E:
    """End-to-end tests for DDL operations."""

    def test_create_table(self, db_connection):
        """Test creating a table."""
        cursor = db_connection.cursor()
        cursor.execute("""
        CREATE TABLE IF NOT EXISTS temp_table (
            id SERIAL PRIMARY KEY,
            value TEXT
        );
        """)
        # Verify table exists
        cursor.execute("""
        SELECT table_name FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'temp_table';
        """)
        result = cursor.fetchone()
        assert result is not None
        assert result[0] == 'temp_table'
        # Cleanup
        cursor.execute("DROP TABLE IF EXISTS temp_table;")
        cursor.close()

    def test_alter_table(self, db_connection, setup_tables):
        """Test altering a table."""
        cursor = db_connection.cursor()
        # Add a column
        cursor.execute("ALTER TABLE test_users ADD COLUMN age INTEGER;")
        # Verify column exists
        cursor.execute("""
        SELECT column_name FROM information_schema.columns
        WHERE table_name = 'test_users' AND column_name = 'age';
        """)
        result = cursor.fetchone()
        assert result is not None
        assert result[0] == 'age'
        cursor.close()

    def test_drop_table(self, db_connection):
        """Test dropping a table."""
        cursor = db_connection.cursor()
        # Create a table to drop
        cursor.execute("CREATE TABLE temp_drop_table (id SERIAL PRIMARY KEY);")
        # Drop it
        cursor.execute("DROP TABLE temp_drop_table;")
        # Verify it's gone
        cursor.execute("""
        SELECT table_name FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name = 'temp_drop_table';
        """)
        result = cursor.fetchone()
        assert result is None
        cursor.close()
