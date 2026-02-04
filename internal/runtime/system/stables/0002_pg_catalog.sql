-- pg_catalog global system schema
-- Инициализация системных таблиц для базы 

CREATE TABLE IF NOT EXISTS pg_catalog.pg_tables (
    schemaname TEXT,
    schematable TEXT,
    tableowner TEXT,
    tablespace TEXT, -- TODO CHECK
    hasindexes BOOLEAN,
    hasrules BOOLEAN,
    hastriggers BOOLEAN,
    rowsecurity BOOLEAN,
    PRIMARY KEY (schemaname, schematable)
);

INSERT INTO pg_catalog.pg_tables (schemaname, schematable, tableowner, tablespace, hasindexes, hasrules, hastriggers, rowsecurity) VALUES
('pg_catalog', 'pg_tables', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_columns', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_class', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_attribute', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_proc', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_type', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_namespace', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_database', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_tablespace', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_roles', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_stat_ssl', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_shdescription', 'system', NULL, FALSE, FALSE, FALSE, FALSE),
('pg_catalog', 'pg_am', 'system', NULL, FALSE, FALSE, FALSE, FALSE);


CREATE TABLE IF NOT EXISTS pg_catalog.pg_columns (
    table_schema TEXT,
    table_name TEXT,
    column_name TEXT,
    original_position INTEGER,
    column_default TEXT, -- TODO check default
    is_nullable BOOLEAN,
    data_type TEXT,
    character_maximum_length INTEGER, -- TODO check
    numeric_precision INTEGER, -- TODO check
    PRIMARY KEY (table_schema, table_name, column_name)
);

INSERT INTO pg_catalog.pg_columns (table_schema, table_name, column_name, original_position, column_default, is_nullable, data_type, character_maximum_length, numeric_precision) VALUES
('pg_catalog', 'pg_tables',  'schemaname', 1, NULL, FALSE, 'TEXT', NULL, NULL),
('pg_catalog', 'pg_tables',  'schematable', 2, NULL, FALSE, 'TEXT', NULL, NULL),
('pg_catalog', 'pg_tables',  'tableowner', 3, NULL, FALSE, 'TEXT', NULL, NULL),
('pg_catalog', 'pg_tables',  'tablespace', 4, NULL, TRUE, 'TEXT', NULL, NULL),
('pg_catalog', 'pg_tables',  'hasindexes', 5, NULL, FALSE, 'BOOLEAN', NULL, NULL),
('pg_catalog', 'pg_tables',  'hasrules', 6, NULL, FALSE, 'BOOLEAN', NULL, NULL),
('pg_catalog', 'pg_tables',  'hastriggers', 7, NULL, FALSE, 'BOOLEAN', NULL, NULL),
('pg_catalog', 'pg_tables',  'rowsecurity', 8, NULL, FALSE, 'BOOLEAN', NULL, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_class (
    oid INTEGER,
    relname TEXT,
    relnamespace INTEGER,
    reltype INTEGER,
    reloftype INTEGER,
    relowner INTEGER,
    relam INTEGER,
    relfilenode INTEGER,
    reltablespace INTEGER,
    relpages INTEGER,
    reltuples FLOAT,
    relallvisible INTEGER,
    reltoastrelid INTEGER,
    relhasindex BOOLEAN,
    relisshared BOOLEAN,
    relpersistence TEXT,
    relkind TEXT,
    relnatts INTEGER,
    relchecks INTEGER,
    relhasrules BOOLEAN,
    relhastriggers BOOLEAN,
    relhassubclass BOOLEAN,
    relrowsecurity BOOLEAN,
    relforcerowsecurity BOOLEAN,
    relispopulated BOOLEAN,
    relreplident TEXT,
    relispartition BOOLEAN,
    relrewrite INTEGER,
    relfrozenxid INTEGER,
    relminmxid INTEGER,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_class (
    oid, relname, relnamespace, reltype, reloftype, relowner,
    relam, relfilenode, reltablespace, relpages, reltuples,
    relallvisible, reltoastrelid, relhasindex, relisshared,
    relpersistence, relkind, relnatts, relchecks, relhasrules,
    relhastriggers, relhassubclass, relrowsecurity,
    relforcerowsecurity, relispopulated, relreplident,
    relispartition, relrewrite, relfrozenxid, relminmxid
) VALUES
(1259, 'pg_tables', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 8, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1260, 'pg_columns', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 9, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1261, 'pg_class', 11, 83, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 30, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1249, 'pg_attribute', 11, 75, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 23, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1255, 'pg_proc', 11, 81, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 29, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1247, 'pg_type', 11, 71, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 31, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(2615, 'pg_namespace', 11, 2615, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 4, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1262, 'pg_database', 11, 1248, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 14, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1213, 'pg_tablespace', 11, 1213, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 5, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1260, 'pg_roles', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 13, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1262, 'pg_stat_ssl', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 7, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(1263, 'pg_shdescription', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 4, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0),
(2601, 'pg_am', 11, 0, 0, 10, 0, 0, 0, 0, 0.0, 0, 0, FALSE, FALSE, 'p', 'r', 4, 0, FALSE, FALSE, FALSE, FALSE, FALSE, TRUE, 'n', FALSE, 0, 0, 0);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_attribute (
    attrelid INTEGER,
    attname TEXT,
    atttypid INT,
    attstattarget INT,
    attlen INT,
    attnum INT,
    attndims INT,
    attcacheoff INT,
    atttypmod INT,
    attbyval BOOLEAN,
    attstorage CHAR(1),
    attalign CHAR(1),
    attnotnull BOOLEAN,
    atthasdef BOOLEAN,
    attidentity CHAR(1),
    attgenerated CHAR(1),
    attisdropped BOOLEAN,
    attislocal BOOLEAN,
    attinhcount INT,
    attcollation INT,
    attacl TEXT,
    attoptions TEXT,
    attfdwoptions TEXT,
    PRIMARY KEY (attrelid, attname)
);

INSERT INTO pg_catalog.pg_attribute (attrelid, attname, atttypid, attstattarget, attlen, attnum, attndims, attcacheoff, atttypmod, attbyval, attstorage, attalign, attnotnull, atthasdef, attidentity, attgenerated, attisdropped, attislocal, attinhcount, attcollation, attacl, attoptions, attfdwoptions) VALUES
(1, 'schemaname', 25, -1, -1, 1, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(1, 'schematable', 25, -1, -1, 2, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(1, 'tableowner', 25, -1, -1, 3, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(1, 'tablespace', 25, -1, -1, 4, 0, -1, -1, FALSE, 'p', 'i', TRUE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(1, 'hasindexes', 16, -1, 1, 5, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0,NULL , NULL , NULL),
(1,'hasrules',16,-1 ,1 ,6 ,0 ,-1 ,-1 ,FALSE ,'p' ,'i' ,FALSE ,FALSE ,'' ,'' ,FALSE ,TRUE ,0 ,0 ,NULL ,NULL ,NULL),
(1,'hastriggers',16 ,-1 ,1 ,7 ,0 ,-1 ,-1 ,FALSE ,'p' ,'i' ,FALSE ,FALSE ,'' ,'' ,FALSE ,TRUE ,0 ,0 ,NULL ,NULL ,NULL),
(1,'rowsecurity',16 ,-1 ,1 ,8 ,0 ,-1 ,-1 ,FALSE ,'p' ,'i' ,FALSE ,FALSE ,'' ,'' ,FALSE ,TRUE ,0 ,0 ,NULL ,NULL ,NULL),
(2, 'table_schema', 25, -1, -1, 1, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'table_name', 25, -1, -1, 2, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'column_name', 25, -1, -1, 3, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'original_position', 23, -1, 4, 4, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'column_default', 25, -1, -1, 5, 0, -1, -1, TRUE, 'p', 'i', TRUE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'is_nullable', 16, -1, 1, 6, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'data_type', 25, -1, -1, 7, 0, -1, -1, FALSE, 'p', 'i', FALSE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'character_maximum_length', 23, -1, 4, 8, 0, -1, -1, TRUE, 'p', 'i', TRUE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL),
(2, 'numeric_precision', 23, -1, 4, 9, 0, -1, -1, TRUE, 'p', 'i', TRUE, FALSE, '', '', FALSE, TRUE, 0, 0, NULL, NULL, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_proc (
    oid INT,
    proname TEXT,
    pronamespace INT,
    proowner INT,
    prolang INT,
    procost FLOAT,
    prorows FLOAT,
    provariadic INT,
    protransform TEXT,
    proisagg BOOLEAN,
    proiswindow BOOLEAN,
    prosecdef BOOLEAN,
    proleakproof BOOLEAN,
    proisstrict BOOLEAN,
    proretset BOOLEAN,
    provolatile CHAR(1),
    proparallel CHAR(1),
    pronargs INT,
    pronargdefaults INT,
    prorettype INT,
    proargtypes TEXT,
    proallargtypes TEXT,
    proargmodes TEXT,
    proargnames TEXT,
    proargdefaults TEXT,
    protrftypes TEXT,
    prosrc TEXT,
    probin TEXT,
    proconfig TEXT,
    proacl TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_type (
    oid INT,
    typname TEXT,
    typnamespace INT,
    typowner INT,
    typlen INT,
    typbyval BOOLEAN,
    typtype CHAR(1),
    typcategory CHAR(1),
    typispreferred BOOLEAN,
    typisdefined BOOLEAN,
    typdelim CHAR(1),
    typrelid INT,
    typelem INT,
    typarray INT,
    typinput TEXT,
    typoutput TEXT,
    typreceive TEXT,
    typsend TEXT,
    typmodin TEXT,
    typmodout TEXT,
    typanalyze TEXT,
    typalign CHAR(1),
    typstorage CHAR(1),
    typnotnull BOOLEAN,
    typbasetype INT,
    typtypmod INT,
    typndims INT,
    typcollation INT,
    typdefaultbin TEXT,
    typdefault TEXT,
    typacl TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_type (oid, typname, typnamespace, typowner, typlen, typbyval, typtype, typcategory, typispreferred, typisdefined, typdelim, typrelid, typelem, typarray, typinput, typoutput, typreceive, typsend, typmodin, typmodout, typanalyze, typalign, typstorage, typnotnull, typbasetype, typtypmod, typndims, typcollation, typdefaultbin, typdefault, typacl) VALUES
(16, 'bool', 11, 10, 1, TRUE, 'b', 'B', TRUE, TRUE, ',', 0, 0, 1000, 'boolin', 'boolout', 'boolrecv', 'boolsend', 'booltypmodin', 'booltypmodout', 'boolanalyze', 'c', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(17, 'bytea', 11, 10, -1, FALSE, 'b', 'U', FALSE, TRUE, ',', 0, 0, 1001, 'byteain', 'byteaout', 'bytearecv', 'byteasend', '-', '-', '-', 'i', 'x', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(18, 'char', 11, 10, 1, TRUE, 'b', 'Z', FALSE, TRUE, ',', 0, 0, 1002, 'charin', 'charout', 'charrecv', 'charsend', '-', '-', '-', 'c', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(19, 'name', 11, 10, 64, FALSE, 'b', 'S', FALSE, TRUE, ',', 0, 0, 1003, 'namein', 'nameout', 'namerecv', 'namesend', '-', '-', '-', 'c', 'p', FALSE, 0, -1, 0, 950, NULL, NULL, NULL),
(20, 'int8', 11, 10, 8, TRUE, 'b', 'N', FALSE, TRUE, ',', 0, 0, 1016, 'int8in', 'int8out', 'int8recv', 'int8send', '-', '-', '-', 'd', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(21, 'int2', 11, 10, 2, TRUE, 'b', 'N', FALSE, TRUE, ',', 0, 0, 1005, 'int2in', 'int2out', 'int2recv', 'int2send', '-', '-', '-', 's', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(23, 'int4', 11, 10, 4, TRUE, 'b', 'N', TRUE, TRUE, ',', 0, 0, 1007, 'int4in', 'int4out', 'int4recv', 'int4send', 'int4typmodin', 'int4typmodout', 'int4analyze', 'i', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(24, 'regproc', 11, 10, 4, TRUE, 'b', 'N', FALSE, TRUE, ',', 0, 0, 1008, 'regprocin', 'regprocout', 'regprocrecv', 'regprocsend', '-', '-', '-', 'i', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(25, 'text', 11, 10, -1, FALSE, 'b', 'S', TRUE, TRUE, ',', 0, 0, 1009, 'textin', 'textout', 'textrecv', 'textsend', 'texttypmodin', 'texttypmodout', 'textanalyze', 'i', 'e', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(26, 'oid', 11, 10, 4, TRUE, 'b', 'N', TRUE, TRUE, ',', 0, 0, 1028, 'oidin', 'oidout', 'oidrecv', 'oidsend', '-', '-', '-', 'i', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(700, 'float4', 11, 10, 4, TRUE, 'b', 'N', FALSE, TRUE, ',', 0, 0, 1021, 'float4in', 'float4out', 'float4recv', 'float4send', '-', '-', '-', 'i', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(701, 'float8', 11, 10, 8, TRUE, 'b', 'N', TRUE, TRUE, ',', 0, 0, 1022, 'float8in', 'float8out', 'float8recv', 'float8send', '-', '-', '-', 'd', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1043, 'varchar', 11, 10, -1, FALSE, 'b', 'S', FALSE, TRUE, ',', 0, 0, 1015, 'varcharin', 'varcharout', 'varcharrecv', 'varcharsend', 'varchartypmodin', 'varchartypmodout', '-', 'i', 'e', FALSE, 0, -1, 0, 100, NULL, NULL, NULL),
(1082, 'date', 11, 10, 4, TRUE, 'b', 'D', FALSE, TRUE, ',', 0, 0, 1182, 'date_in', 'date_out', 'date_recv', 'date_send', '-', '-', '-', 'i', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1083, 'time', 11, 10, 8, TRUE, 'b', 'D', FALSE, TRUE, ',', 0, 0, 1183, 'time_in', 'time_out', 'time_recv', 'time_send', 'timetypmodin', 'timetypmodout', '-', 'd', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1114, 'timestamp', 11, 10, 8, TRUE, 'b', 'D', FALSE, TRUE, ',', 0, 0, 1115, 'timestamp_in', 'timestamp_out', 'timestamp_recv', 'timestamp_send', 'timestamptypmodin', 'timestamptypmodout', '-', 'd', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1184, 'timestamptz', 11, 10, 8, TRUE, 'b', 'D', TRUE, TRUE, ',', 0, 0, 1185, 'timestamptz_in', 'timestamptz_out', 'timestamptz_recv', 'timestamptz_send', 'timestamptztypmodin', 'timestamptztypmodout', '-', 'd', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1700, 'numeric', 11, 10, -1, FALSE, 'b', 'N', FALSE, TRUE, ',', 0, 0, 1231, 'numeric_in', 'numeric_out', 'numeric_recv', 'numeric_send', 'numerictypmodin', 'numerictypmodout', '-', 'i', 'm', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(2950, 'uuid', 11, 10, 16, FALSE, 'b', 'U', FALSE, TRUE, ',', 0, 0, 2951, 'uuid_in', 'uuid_out', 'uuid_recv', 'uuid_send', '-', '-', '-', 'c', 'p', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(114, 'json', 11, 10, -1, FALSE, 'b', 'U', FALSE, TRUE, ',', 0, 0, 199, 'json_in', 'json_out', 'json_recv', 'json_send', '-', '-', '-', 'i', 'e', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(3802, 'jsonb', 11, 10, -1, FALSE, 'b', 'U', FALSE, TRUE, ',', 0, 0, 3807, 'jsonb_in', 'jsonb_out', 'jsonb_recv', 'jsonb_send', '-', '-', '-', 'i', 'e', FALSE, 0, -1, 0, 0, NULL, NULL, NULL),
(1042, 'bpchar', 11, 10, -1, FALSE, 'b', 'S', FALSE, TRUE, ',', 0, 0, 1014, 'bpcharin', 'bpcharout', 'bpcharrecv', 'bpcharsend', 'bpchartypmodin', 'bpchartypmodout', '-', 'i', 'e', FALSE, 0, -1, 0, 100, NULL, NULL, NULL),
-- Domain types examples (typtype='d' based on base types)
(100001, 'positive_int', 11, 10, 4, TRUE, 'd', 'N', FALSE, TRUE, ',', 0, 0, 0, 'int4in', 'int4out', 'int4recv', 'int4send', '-', '-', '-', 'i', 'p', FALSE, 23, -1, 0, 0, NULL, NULL, NULL),
(100002, 'email_address', 11, 10, -1, FALSE, 'd', 'S', FALSE, TRUE, ',', 0, 0, 0, 'textin', 'textout', 'textrecv', 'textsend', '-', '-', '-', 'i', 'e', FALSE, 25, -1, 0, 0, NULL, NULL, NULL),
(100003, 'us_postal_code', 11, 10, -1, FALSE, 'd', 'S', FALSE, TRUE, ',', 0, 0, 0, 'textin', 'textout', 'textrecv', 'textsend', '-', '-', '-', 'i', 'e', FALSE, 25, -1, 0, 0, NULL, NULL, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_namespace (
    oid INT,
    nspname TEXT,
    nspowner INT,
    nspacl TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_namespace (oid, nspname, nspowner, nspacl) VALUES
(11, 'pg_catalog', 10, NULL),
(12, 'public', 10, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_database (
    oid INT,
    datname TEXT,
    datdba INT,
    encoding INT,
    datcollate TEXT,
    datctype TEXT,
    datistemplate BOOLEAN,
    datallowconn BOOLEAN,
    datconnlimit INT,
    datlastsysoid INT,
    datfrozenxid INT,
    datminmxid INT,
    dattablespace INT,
    datacl TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_database (oid, datname, datdba, encoding, datcollate, datctype, datistemplate, datallowconn, datconnlimit, datlastsysoid, datfrozenxid, datminmxid, dattablespace, datacl) VALUES
(1, 'postgres', 10, 6, 'en_US.UTF-8', 'en_US.UTF-8', FALSE, TRUE, -1, 0, 0, 0, 1663, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_tablespace (
    oid INT,
    spcname TEXT,
    spcowner INT,
    spcacl TEXT,
    spcoptions TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_tablespace (oid, spcname, spcowner, spcacl, spcoptions) VALUES
(1663, 'pg_default', 10, NULL, NULL),
(1664, 'pg_global', 10, NULL, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_roles (
    oid INT,
    rolname TEXT,
    rolsuper BOOLEAN,
    rolinherit BOOLEAN,
    rolcreaterole BOOLEAN,
    rolcreatedb BOOLEAN,
    rolcanlogin BOOLEAN,
    rolreplication BOOLEAN,
    rolconnlimit INT,
    rolpassword TEXT,
    rolvaliduntil TEXT,
    rolbypassrls BOOLEAN,
    rolconfig TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_roles (oid, rolname, rolsuper, rolinherit, rolcreaterole, rolcreatedb, rolcanlogin, rolreplication, rolconnlimit, rolpassword, rolvaliduntil, rolbypassrls, rolconfig) VALUES
(10, 'system', TRUE, TRUE, TRUE, TRUE, TRUE, TRUE, -1, NULL, NULL, FALSE, NULL),
(11, 'user1', FALSE, TRUE, FALSE, TRUE, TRUE, FALSE, -1, NULL, NULL, FALSE, NULL);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_stat_ssl (
    pid INT,
    ssl BOOLEAN,
    version TEXT,
    cipher TEXT,
    bits INT,
    compression BOOLEAN,
    clientdn TEXT,
    PRIMARY KEY (pid)
);

INSERT INTO pg_catalog.pg_stat_ssl (pid, ssl, version, cipher, bits, compression, clientdn) VALUES
(12345, TRUE, 'TLSv1.3', 'TLS_AES_256_GCM_SHA384', 256, FALSE, 'CN=client,O=example,C=US');

CREATE TABLE IF NOT EXISTS pg_catalog.pg_shdescription (
    objoid INT,
    classoid INT,
    objsubid INT,
    description TEXT,
    PRIMARY KEY (objoid, classoid, objsubid)
);

INSERT INTO pg_catalog.pg_shdescription (objoid, classoid, objsubid, description) VALUES
(1, 1259, 0, 'Description for pg_tables table'),
(2, 1259, 0, 'Description for pg_columns table');

CREATE TABLE IF NOT EXISTS pg_catalog.pg_am (
    oid INT,
    amname TEXT,
    amhandler INT,
    amtype TEXT,
    PRIMARY KEY (oid)
);

INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES
(403, 'heap', 0, 't'),
(405, 'btree', 0, 'i');
