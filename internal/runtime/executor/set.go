package executor

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/storage"
	"strings"
)

// ExecuteSet executes the SET statement to set session variables.
// TODO: Implement comprehensive variable setting logic with validation
func ExecuteSet(stmt *statements.SetStatement, store *storage.DistributedStorageVClock, sessionParams map[string]string, exCtx ExecutorContext) error {
	// Set the session parameter
	sessionParams[strings.ToLower(stmt.Variable)] = fmt.Sprintf("%v", stmt.Value)
	logger.Debugf("SET %s = %v\n", stmt.Variable, stmt.Value)
	return nil
}
