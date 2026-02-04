package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
)

type ExecutorContext struct {
	Database                string
	Schema                  string
	IsInformationSchemaInit bool
}

func ExecuteStatement(stmt statements.SQLStatement, storage *storage.DistributedStorageVClock, sessionParams map[string]string) (*table.ExecuteResult, error) {
	exCtx := ExecutorContext{
		Database: "testdb",
		Schema:   "public",
	}

	if val, ok := sessionParams["__information_schema"]; ok {
		exCtx.IsInformationSchemaInit = val == "true"
	}

	switch s := stmt.(type) {
	case *statements.EmptyStatement:
		return nil, nil
	case *statements.CreateTableStatement:
		return nil, ExecuteCreateTable(s, storage, exCtx)
	case *statements.InsertStatement:
		return nil, ExecuteInsert(s, storage, exCtx)
	case *statements.SelectStatement:
		return ExecuteSelect(s, storage, exCtx)
	case *statements.DropTableStatement:
		return nil, ExecuteDropTable(s, storage, exCtx)
	case *statements.TruncateTableStatement:
		return nil, ExecuteTruncateTable(s, storage, exCtx)
	case *statements.SetStatement:
		return nil, ExecuteSet(s, storage, sessionParams, exCtx)
	case *statements.ShowStatement:
		return ExecuteShow(s, storage, sessionParams, exCtx)
	default:
		return nil, fmt.Errorf("unsupported statement type: %s", stmt.Type())
	}
}
