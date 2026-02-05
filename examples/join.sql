select t.oid,
	case when nsp.nspname in ('pg_catalog', 'public') then t.typname
		else nsp.nspname||'.'||t.typname
	end
from pg_type t
left join pg_type base_type on t.typelem=base_type.oid
left join pg_class base_cls ON base_type.typrelid = base_cls.oid
left join pg_namespace nsp on t.typnamespace=nsp.oid
where (
	  t.typtype in('b', 'p', 'r', 'e')
	  and (base_type.oid is null or base_type.typtype in('b', 'p', 'r'))
	);

select t.oid from pg_type t;

select base_type.oid, t.typelem, t.oid from pg_type t
left join pg_type base_type on t.typelem=base_type.oid
left join pg_class base_cls ON base_type.typrelid = base_cls.oid;


CREATE TABLE orders (
    id INT PRIMARY KEY,
    user_id INT
);

CREATE TABLE users (
    id INT PRIMARY KEY,
    name TEXT
);

INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'Bob');
INSERT INTO orders (id, user_id) VALUES (1, 1), (2, 2), (3, 1);

SELECT u.id, o.user_id, * FROM orders o
JOIN users u ON o.user_id = u.id;

SELECT pg_catalog.quote_ident(c.relname) FROM pg_catalog.pg_class c WHERE c.relkind IN ('r', 'S', 'v', 'm', 'f', 'p') AND substring(pg_catalog.quote_ident(c.relname),1,5)='pg_ty' AND pg_catalog.pg_table_is_visible(c.oid)
UNION
SELECT pg_catalog.quote_ident(n.nspname) || '.' FROM pg_catalog.pg_namespace n WHERE substring(pg_catalog.quote_ident(n.nspname) || '.',1,5)='pg_ty' AND (SELECT pg_catalog.count(*) FROM pg_catalog.pg_namespace WHERE substring(pg_catalog.quote_ident(nspname) || '.',1,5) = substring('pg_ty',1,pg_catalog.length(pg_catalog.quote_ident(nspname))+1)) > 1
UNION
SELECT pg_catalog.quote_ident(n.nspname) || '.' || pg_catalog.quote_ident(c.relname) FROM pg_catalog.pg_class c, pg_catalog.pg_namespace n WHERE c.relnamespace = n.oid AND c.relkind IN ('r', 'S', 'v', 'm', 'f', 'p') AND substring(pg_catalog.quote_ident(n.nspname) || '.' || pg_catalog.quote_ident(c.relname),1,5)='pg_ty' AND substring(pg_catalog.quote_ident(n.nspname) || '.',1,5) = substring('pg_ty',1,pg_catalog.length(pg_catalog.quote_ident(n.nspname))+1) AND (SELECT pg_catalog.count(*) FROM pg_catalog.pg_namespace WHERE substring(pg_catalog.quote_ident(nspname) || '.',1,5) = substring('pg_ty',1,pg_catalog.length(pg_catalog.quote_ident(nspname))+1)) = 1
LIMIT 1000
