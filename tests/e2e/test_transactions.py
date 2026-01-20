import pytest
import psycopg2
from fixtures.tables import db_connection, setup_tables

class TestTransactionsE2E:
    """End-to-end tests for transactions."""

    def test_transaction_commit(self, db_connection, setup_tables):
        """Test transaction commit."""
        conn = psycopg2.connect(**db_connection.dsn)  # New connection for transaction
        conn.autocommit = False
        cursor = conn.cursor()
        try:
            cursor.execute("INSERT INTO test_users (name, email) VALUES (%s, %s);", ('Eve', 'eve@example.com'))
            conn.commit()
            # Verify
            cursor.execute("SELECT COUNT(*) FROM test_users;")
            result = cursor.fetchone()
            assert result[0] == 4
        finally:
            cursor.close()
            conn.close()

    def test_transaction_rollback(self, db_connection, setup_tables):
        """Test transaction rollback."""
        conn = psycopg2.connect(**db_connection.dsn)
        conn.autocommit = False
        cursor = conn.cursor()
        try:
            cursor.execute("INSERT INTO test_users (name, email) VALUES (%s, %s);", ('Frank', 'frank@example.com'))
            conn.rollback()
            # Verify not inserted
            cursor.execute("SELECT COUNT(*) FROM test_users;")
            result = cursor.fetchone()
            assert result[0] == 3  # Original count
        finally:
            cursor.close()
            conn.close()

    def test_transaction_isolation(self, db_connection, setup_tables):
        """Test transaction isolation (basic check)."""
        conn1 = psycopg2.connect(**db_connection.dsn)
        conn2 = psycopg2.connect(**db_connection.dsn)
        conn1.autocommit = False
        conn2.autocommit = False
        cursor1 = conn1.cursor()
        cursor2 = conn2.cursor()
        try:
            # Insert in conn1, not committed
            cursor1.execute("INSERT INTO test_users (name, email) VALUES (%s, %s);", ('Grace', 'grace@example.com'))
            # Check in conn2 - should not see
            cursor2.execute("SELECT COUNT(*) FROM test_users;")
            result = cursor2.fetchone()
            assert result[0] == 3  # Not committed yet
            conn1.commit()
            # Now should see
            cursor2.execute("SELECT COUNT(*) FROM test_users;")
            result = cursor2.fetchone()
            assert result[0] == 4
        finally:
            cursor1.close()
            cursor2.close()
            conn1.close()
            conn2.close()