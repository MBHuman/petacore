SELECT
    N.oid::bigint                            AS id,
    datname                                 AS name,
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
WHERE c.relname OPERATOR(pg_catalog.~) '^(pg_.*)$' COLLATE pg_catalog.default
  AND pg_catalog.pg_table_is_visible(c.oid)
ORDER BY 2, 3;

SELECT c.oid,
  n.nspname,
  c.relname
FROM pg_catalog.pg_class c
     LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
WHERE c.relname OPERATOR(pg_catalog.~) '^(pg_class)$' COLLATE pg_catalog.default
  AND n.nspname OPERATOR(pg_catalog.~) '^(pg_catalog)$' COLLATE pg_catalog.default
ORDER BY 2, 3;

CREATE TABLE orders (
    order_id INT PRIMARY KEY,
    order_date TEXT NOT NULL,
    customer_id INT NOT NULL,
    amount DECIMAL NOT NULL
);

INSERT INTO orders (order_id, order_date, customer_id, amount) VALUES
    (1, '2024-01-15', 101, 250.75),
    (2, '2024-02-20', 102, 125.00),
    (3, '2024-03-05', 103, 300.50);


CREATE TABLE IF NOT EXISTS customers (
    customer_id INT PRIMARY KEY,
    customer_name TEXT NOT NULL,
    contact_email TEXT NOT NULL
);

INSERT INTO customers (customer_id, customer_name, contact_email) VALUES
    (101, 'Alice Smith', 'alice.smith@example.com'),
    (102, 'Bob Johnson', 'bob.jognsno@example.com'),
    (103, 'Charlie Brown', 'charlize.b@example.com');

SELECT o.order_id, o.order_date, c.customer_name, o.amount
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
WHERE o.amount > 200.00
ORDER BY o.order_date;