SELECT
    N.oid::bigint                            AS id,
    datname                                 AS name
    D.description,
    datistemplate                           AS is_template,
    datallowconn                            AS allow_connections,
    pg_catalog.pg_get_userbyid(N.datdba)    AS "owner"
FROM pg_catalog.pg_database N
LEFT JOIN pg_catalog.pg_shdescription D
       ON N.oid = D.objoid
ORDER BY
    CASE
        WHEN datname = pg_catalog.current_database()
            THEN -1::bigint
        ELSE N.oid::bigint
    END;

CREATE TABLE IF NOT EXISTS test_nntable (
    id INT PRIMARY KEY,
    name TEXT NOT NULL
);

-- query for \d ${table_name}
SELECT c.oid,
  n.nspname,
  c.relname
FROM pg_catalog.pg_class c
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relname OPERATOR(pg_catalog.~) '^(test_nntable)$' COLLATE pg_catalog.default
  AND pg_catalog.pg_table_is_visible(c.oid)
ORDER BY 2, 3;