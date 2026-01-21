package wire

import (
	"petacore/internal/runtime/rsql/statements"

	"github.com/jackc/pgproto3/v2"
)

// PreparedStatement represents a prepared statement
type PreparedStatement struct {
	Query             string
	Stmt              statements.SQLStatement
	Params            []interface{}
	Columns           []pgproto3.FieldDescription
	ParamOIDs         []uint32
	ResultFormatCodes []int16
}
