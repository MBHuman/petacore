-- information_schema global system schema
--- Инициализация системных таблиц для базы  --- IGNORE ---

CREATE TABLE IF NOT EXISTS information_schema.tables (
    table_schema TEXT,
    table_name TEXT,
    table_type TEXT
);

CREATE TABLE IF NOT EXISTS information_schema.columns (
    table_schema TEXT,
    table_name TEXT,
    column_name TEXT,
    original_position INT,
    column_default TEXT, -- TODO check default
    is_nullable BOOLEAN,
    data_type TEXT,
);

CREATE TABLE IF NOT EXISTS information_schema.schemata (
    catalog_name TEXT,
    schema_name TEXT,
    schema_owner TEXT,
    default_character_set_catalog TEXT,
    default_character_set_schema TEXT,
    default_character_set_name TEXT,
    sql_path TEXT
);
