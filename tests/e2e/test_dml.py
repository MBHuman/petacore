import pytest
import psycopg2
from fixtures.tables import db_connection, setup_tables

class TestDMLE2E:
    """End-to-end tests for DML operations."""

    def test_insert(self, db_connection, setup_tables):
        """Test INSERT operation."""
        cursor = db_connection.cursor()
        cursor.execute("INSERT INTO test_users (name, email) VALUES (%s, %s);", ('David', 'david@example.com'))
        cursor.execute("SELECT COUNT(*) FROM test_users;")
        result = cursor.fetchone()
        assert result[0] == 4
        cursor.close()

    def test_update(self, db_connection, setup_tables):
        """Test UPDATE operation."""
        cursor = db_connection.cursor()
        cursor.execute("UPDATE test_users SET name = %s WHERE email = %s;", ('Alice Updated', 'alice@example.com'))
        cursor.execute("SELECT name FROM test_users WHERE email = %s;", ('alice@example.com',))
        result = cursor.fetchone()
        assert result[0] == 'Alice Updated'
        cursor.close()

    def test_delete(self, db_connection, setup_tables):
        """Test DELETE operation."""
        cursor = db_connection.cursor()
        cursor.execute("DELETE FROM test_users WHERE email = %s;", ('charlie@example.com',))
        cursor.execute("SELECT COUNT(*) FROM test_users;")
        result = cursor.fetchone()
        assert result[0] == 2
        cursor.close()