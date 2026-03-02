package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
	ptypes "petacore/sdk/types"
)

// parseIntersectStatement парсит INTERSECT операции
func ParseIntersectStatement(ctx parser.IIntersectStatementContext) (*statements.SelectStatement, error) {
	var err error
	if ctx == nil {
		return nil, nil
	}

	primaryCtxs := ctx.AllPrimarySelectStatement()
	if len(primaryCtxs) == 0 {
		return nil, nil
	}

	// Парсим первый primary statement
	result, err := ParsePrimarySelectStatement(primaryCtxs[0])
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	// Если есть INTERSECT операции, обрабатываем их слева направо
	if len(primaryCtxs) > 1 {
		alls := ctx.AllALL()
		allIdx := 0

		for i := 1; i < len(primaryCtxs); i++ {
			right, err := ParsePrimarySelectStatement(primaryCtxs[i])
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, nil
			}

			combined := &statements.CombinedSelectStatement{
				OperationType: statements.IntersectOperation,
				Left:          result,
				Right:         right,
			}

			// Проверяем наличие ALL
			if allIdx < len(alls) && alls[allIdx] != nil {
				combined.All = true
				allIdx++
			}

			result = &statements.SelectStatement{
				Combined:      combined,
				SubqueryCache: make(map[*statements.SelectStatement]*ptypes.Row),
			}
		}
	}

	return result, nil
}
