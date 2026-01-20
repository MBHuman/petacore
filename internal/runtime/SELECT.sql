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