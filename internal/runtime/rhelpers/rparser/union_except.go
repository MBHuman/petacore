package rparser

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

// parseUnionExceptStatement парсит UNION/EXCEPT операции (левоассоциативные)
func ParseUnionExceptStatement(ctx parser.IUnionExceptStatementContext) (*statements.SelectStatement, error) {
	if ctx == nil {
		return nil, nil
	}

	intersectCtxs := ctx.AllIntersectStatement()
	if len(intersectCtxs) == 0 {
		return nil, nil
	}

	// Парсим первый intersect statement
	result, err := ParseIntersectStatement(intersectCtxs[0])
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	// Если есть UNION/EXCEPT операции, обрабатываем их слева направо
	if len(intersectCtxs) > 1 {
		unions := ctx.AllUNION()
		excepts := ctx.AllEXCEPT()
		alls := ctx.AllALL()

		allIdx := 0
		for i := 1; i < len(intersectCtxs); i++ {
			right, err := ParseIntersectStatement(intersectCtxs[i])
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, nil
			}

			combined := &statements.CombinedSelectStatement{
				Left:  result,
				Right: right,
			}

			// Определяем тип операции (UNION или EXCEPT)
			// unions и excepts содержат токены в порядке появления
			operationIdx := i - 1
			if operationIdx < len(unions) && unions[operationIdx] != nil {
				combined.OperationType = statements.UnionOperation
			} else if operationIdx < len(excepts) && excepts[operationIdx] != nil {
				combined.OperationType = statements.ExceptOperation
			}

			// Проверяем наличие ALL
			if allIdx < len(alls) && alls[allIdx] != nil {
				combined.All = true
				allIdx++
			}

			result = &statements.SelectStatement{
				Combined:      combined,
				SubqueryCache: make(map[*statements.SelectStatement]interface{}),
			}
		}
	}

	return result, nil
}
