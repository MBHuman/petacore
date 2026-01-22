-- pg_catalog global system schema
-- --- Инициализация системных таблиц для базы  --- IGNORE ---

CREATE TABLE IF NOT EXISTS pg_catalog.pg_tables (
    schemaname TEXT,
    schematable TEXT,
    tableowner TEXT,
    tablespace TEXT, -- TODO CHECK
    hasindexes BOOLEAN,
    hasrules BOOLEAN,
    hastriggers BOOLEAN,
    rowsecurity BOOLEAN,
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_columns (
    table_schema TEXT,
    table_name TEXT,
    column_name TEXT,
    original_position INT,
    column_default TEXT, -- TODO check default
    is_nullable BOOLEAN,
    data_type TEXT,
    character_maximum_length INT, -- TODO check
    numeric_precision INT, -- TODO check
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_class (
    oid INT,
    relname TEXT,
    relnamespace INT,
    reltype INT,
    reloftype INT,
    relowner INT,
    relam INT,
    relfilenode INT,
    reltablespace INT,
    relpages INT,
    reltuples FLOAT,
    relallvisible INT,
    reltoastrelid INT,
    relhasindex BOOLEAN,
    relisshared BOOLEAN,
    relpersistence TEXT,
    relkind TEXT,
    relnatts INT,
    relchecks INT,
    relhasrules BOOLEAN,
    relhastriggers BOOLEAN,
    relhassubclass BOOLEAN,
    relrowsecurity BOOLEAN,
    relforcerowsecurity BOOLEAN,
    relispopulated BOOLEAN,
    relreplident TEXT,
    relispartition BOOLEAN,
    relrewrite INT,
    relfrozenxid INT,
    relminmxid INT,
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_attribute (
    attrelid INT,
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
    attfdwoptions TEXT
);

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
    typacl TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_namespace (
    oid INT,
    nspname TEXT,
    nspowner INT,
    nspacl TEXT
);

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
    datacl TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_tablespace (
    oid INT,
    spcname TEXT,
    spcowner INT,
    spcacl TEXT,
    spcoptions TEXT
);

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
    rolconfig TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_stat_ssl (
    pid INT,
    ssl BOOLEAN,
    version TEXT,
    cipher TEXT,
    bits INT,
    compression BOOLEAN,
    clientdn TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_shdescription (
    objoid INT,
    classoid INT,
    objsubid INT,
    description TEXT
);

CREATE TABLE IF NOT EXISTS pg_catalog.pg_am (
    oid INT,
    amname TEXT,
    amhandler INT,
    amtype TEXT
);

INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (2, 'heap', 0, 't');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (403, 'btree', 0, 'i');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (405, 'hash', 0, 'i');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (783, 'gist', 0, 'i');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (2742, 'gin', 0, 'i');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (4000, 'spgist', 0, 'i');
INSERT INTO pg_catalog.pg_am (oid, amname, amhandler, amtype) VALUES (3580, 'brin', 0, 'i');
