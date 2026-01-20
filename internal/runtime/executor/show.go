package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
	"strings"
)

// ExecuteShow выполняет SHOW statement
func ExecuteShow(stmt *statements.ShowStatement, storage *storage.DistributedStorageVClock, sessionParams map[string]string, exCtx ExecutorContext) ([]map[string]interface{}, error) {
	param := strings.ToLower(stmt.Parameter)
	switch param {
	case "transaction isolation level":
		// Return default isolation level
		return []map[string]interface{}{
			{"transaction_isolation": "read committed"},
		}, nil
	default:
		if val, ok := sessionParams[param]; ok {
			return []map[string]interface{}{
				{param: val},
			}, nil
		}
		return nil, fmt.Errorf("unknown parameter: %s", stmt.Parameter)
	}
}
