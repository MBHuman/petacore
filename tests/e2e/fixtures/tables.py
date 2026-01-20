import pytest
import psycopg2
from psycopg2 import sql

# Database connection parameters
DB_CONFIG = {
    'host': 'localhost',
    'port': 5432,
    'database': 'petacore',
    'user': 'petacore',
    'password': 'petacore'
}

@pytest.fixture(scope="session")
def db_connection():
    """Fixture to provide a database connection."""
    conn = psycopg2.connect(**DB_CONFIG)
    conn.autocommit = True  # For DDL operations
    yield conn
    conn.close()

@pytest.fixture(scope="function")
def setup_tables(db_connection):
    """Fixture to create and drop test tables."""
    cursor = db_connection.cursor()

    # Create test table
    create_table_query = """
    CREATE TABLE IF NOT EXISTS test_users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    """
    cursor.execute(create_table_query)

    # Insert some initial data
    insert_data_query = """
    INSERT INTO test_users (name, email) VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com'),
    ('Charlie', 'charlie@example.com');
    """
    cursor.execute(insert_data_query)

    yield

    # Cleanup: drop table
    drop_table_query = "DROP TABLE IF EXISTS test_users;"
    cursor.execute(drop_table_query)
    cursor.close()
