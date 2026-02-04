package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	"strings"
)

// ExecuteShow выполняет SHOW statement
// TODO: расширить поддержку других параметров SHOW
func ExecuteShow(stmt *statements.ShowStatement, storage *storage.DistributedStorageVClock, sessionParams map[string]string, exCtx ExecutorContext) (*table.ExecuteResult, error) {
	param := strings.ToLower(stmt.Parameter)
	switch param {
	case "transaction isolation level":
		return &table.ExecuteResult{
			Rows: [][]interface{}{
				{"read committed"},
			},
			Columns: []table.TableColumn{
				{
					Idx:  0,
					Name: "transaction_isolation",
					Type: table.ColTypeString,
				},
			},
		}, nil
	default:
		if val, ok := sessionParams[param]; ok {
			return &table.ExecuteResult{
				Rows: [][]interface{}{
					{val},
				},
				Columns: []table.TableColumn{
					{
						Idx:  0,
						Name: param,
						Type: table.ColTypeString,
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("unknown parameter: %s", stmt.Parameter)
	}
}
