import pytest
import psycopg2
from fixtures.tables import db_connection, setup_tables

class TestDQLE2E:
    """End-to-end tests for DQL operations."""

    def test_select_all(self, db_connection, setup_tables):
        """Test SELECT all records."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT * FROM test_users;")
        results = cursor.fetchall()
        assert len(results) == 3
        cursor.close()

    def test_select_where(self, db_connection, setup_tables):
        """Test SELECT with WHERE clause."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT name FROM test_users WHERE email = %s;", ('alice@example.com',))
        result = cursor.fetchone()
        assert result[0] == 'Alice'
        cursor.close()

    def test_select_count(self, db_connection, setup_tables):
        """Test SELECT COUNT."""
        cursor = db_connection.cursor()
        cursor.execute("SELECT COUNT(*) FROM test_users;")
        result = cursor.fetchone()
        assert result[0] == 3
        cursor.close()