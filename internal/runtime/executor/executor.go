package executor

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/planner"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"

	"go.uber.org/zap"
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
		return ExecuteSelectWithPlanner(s, storage, exCtx)
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

// ExecuteSelectWithPlanner использует новую архитектуру planner + executor
func ExecuteSelectWithPlanner(stmt *statements.SelectStatement, storage *storage.DistributedStorageVClock, exCtx ExecutorContext) (*table.ExecuteResult, error) {
	// 1. Создаем план выполнения
	plannerCtx := planner.PlannerContext{
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
	}

	queryPlan, err := planner.CreateQueryPlan(stmt, plannerCtx)
	if err != nil {
		return nil, fmt.Errorf("error creating query plan: %w", err)
	}

	// 2. Выполняем план
	executorCtx := planner.ExecutorContext{
		Database: exCtx.Database,
		Schema:   exCtx.Schema,
		Storage:  storage,
		GoCtx:    context.Background(),
	}

	logger.Debugf("Executing query plan: ", zap.Any("queryPlan", queryPlan))
	result, err := planner.ExecutePlan(queryPlan, executorCtx)
	if err != nil {
		return nil, fmt.Errorf("error executing query plan: %w", err)
	}
	logger.Debug("Planner result", zap.Any("result", result))

	return result, nil
}
