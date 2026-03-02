package table

import ptypes "petacore/sdk/types"

type TableColumn struct {
	Idx               int
	Name              string
	Type              ptypes.OID
	TableIdentifier   string // table alias or name (used in query)
	OriginalTableName string // original table name (for error messages)
}

// ColumnDef определяет колонку
type ColumnDef struct {
	Idx  int
	Name string
	Type ptypes.OID
	// IsPrimaryKey bool
	IsNullable   bool
	IsUnique     bool
	IsSerial     bool
	DefaultValue interface{}
}
