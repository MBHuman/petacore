package pgcatalog

import (
	"petacore/internal/runtime/rhelpers"
	"petacore/internal/runtime/rsql/statements"
)

func ExecutePgDatabase(stmt *statements.SelectStatement) ([]map[string]interface{}, error) {
	allRows := []map[string]interface{}{
		{
			"oid":           1,
			"datname":       "postgres",
			"datdba":        10,
			"encoding":      6,
			"datcollate":    "en_US.UTF-8",
			"datctype":      "en_US.UTF-8",
			"datistemplate": false,
			"datallowconn":  true,
			"datconnlimit":  -1,
			"datlastsysoid": 0,
			"datfrozenxid":  0,
			"datminmxid":    1,
			"dattablespace": 1663,
			"datacl":        nil,
		},
		{
			"oid":           13207,
			"datname":       "template0",
			"datdba":        10,
			"encoding":      6,
			"datcollate":    "en_US.UTF-8",
			"datctype":      "en_US.UTF-8",
			"datistemplate": true,
			"datallowconn":  false,
			"datconnlimit":  -1,
			"datlastsysoid": 0,
			"datfrozenxid":  0,
			"datminmxid":    1,
			"dattablespace": 1663,
			"datacl":        nil,
		},
		{
			"oid":           1,
			"datname":       "template1",
			"datdba":        10,
			"encoding":      6,
			"datcollate":    "en_US.UTF-8",
			"datctype":      "en_US.UTF-8",
			"datistemplate": true,
			"datallowconn":  true,
			"datconnlimit":  -1,
			"datlastsysoid": 0,
			"datfrozenxid":  0,
			"datminmxid":    1,
			"dattablespace": 1663,
			"datacl":        nil,
		},
		{
			"oid":           16384,
			"datname":       "testdb",
			"datdba":        10,
			"encoding":      6,
			"datcollate":    "en_US.UTF-8",
			"datctype":      "en_US.UTF-8",
			"datistemplate": false,
			"datallowconn":  true,
			"datconnlimit":  -1,
			"datlastsysoid": 0,
			"datfrozenxid":  0,
			"datminmxid":    1,
			"dattablespace": 1663,
			"datacl":        nil,
		},
	}
	return rhelpers.FilterColumns(allRows, stmt), nil
}
