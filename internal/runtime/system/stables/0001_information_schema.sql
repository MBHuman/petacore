-- information_schema global system schema

CREATE TABLE IF NOT EXISTS information_schema.tables (
    table_schema TEXT,
    table_name TEXT,
    table_type TEXT,
    PRIMARY KEY (table_schema, table_name)
);

INSERT INTO information_schema.tables (table_schema, table_name, table_type) VALUES
    ('information_schema', 'tables', 'BASE TABLE'),
    ('information_schema', 'columns', 'BASE TABLE'),
    ('information_schema', 'schemata', 'BASE TABLE');

CREATE TABLE IF NOT EXISTS information_schema.columns (
    table_schema TEXT,
    table_name TEXT,
    column_name TEXT,
    original_position INT,
    column_default TEXT,
    is_nullable BOOLEAN,
    data_type TEXT,
    PRIMARY KEY (table_schema, table_name, column_name)
);

INSERT INTO information_schema.columns (table_schema, table_name, column_name, original_position, is_nullable, data_type) VALUES
    ('information_schema', 'tables', 'table_schema', 1, FALSE, 'TEXT'),
    ('information_schema', 'tables', 'table_name', 2, FALSE, 'TEXT'),
    ('information_schema', 'tables', 'table_type', 3, FALSE, 'TEXT'),
    ('information_schema', 'columns', 'table_schema', 1, FALSE, 'TEXT'),
    ('information_schema', 'columns', 'table_name', 2, FALSE, 'TEXT'),
    ('information_schema', 'columns', 'column_name', 3, FALSE, 'TEXT'),
    ('information_schema', 'columns', 'original_position', 4, FALSE, 'INT'),
    ('information_schema', 'columns', 'column_default', 5, TRUE, 'TEXT'),
    ('information_schema', 'columns', 'is_nullable', 6, FALSE, 'BOOLEAN'),
    ('information_schema', 'columns', 'data_type', 7, FALSE, 'TEXT'),
    ('information_schema', 'schemata', 'catalog_name', 1, FALSE, 'TEXT'),
    ('information_schema', 'schemata', 'schema_name', 2, FALSE, 'TEXT'),
    ('information_schema', 'schemata', 'schema_owner', 3, FALSE, 'TEXT'),
    ('information_schema', 'schemata', 'default_character_set_catalog', 4, TRUE, 'TEXT'),
    ('information_schema', 'schemata', 'default_character_set_schema', 5, TRUE, 'TEXT'),
    ('information_schema', 'schemata', 'default_character_set_name', 6, TRUE, 'TEXT'),
    ('information_schema', 'schemata', 'sql_path', 7, TRUE, 'TEXT');

CREATE TABLE IF NOT EXISTS information_schema.schemata (
    catalog_name TEXT,
    schema_name TEXT,
    schema_owner TEXT,
    default_character_set_catalog TEXT,
    default_character_set_schema TEXT,
    default_character_set_name TEXT,
    sql_path TEXT,
    PRIMARY KEY (catalog_name, schema_name)
);

INSERT INTO information_schema.schemata (catalog_name, schema_name, schema_owner) VALUES
    ('petacore', 'information_schema', 'system'),
    ('petacore', 'pg_catalog', 'system');
