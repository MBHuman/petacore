package functions

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/rsql/table"
	psdk "petacore/sdk"
	ptypes "petacore/sdk/types"
	"strings"
)

var functionRegistry *psdk.FunctionRegistry

// SetFunctionRegistry sets the global function registry for user-defined functions
func SetFunctionRegistry(registry *psdk.FunctionRegistry) {
	functionRegistry = registry
}

// IsAggregateFunction проверяет, является ли функция агрегатной
func IsAggregateFunction(name string) bool {
	if functionRegistry == nil {
		return false
	}

	if fn, exists := functionRegistry.GetByName(strings.ToUpper(name)); exists {
		return fn.GetFunction().IsAggregate
	}

	return false
}

// ExecuteFunction executes a built-in or user-defined function
// TODO добавить поддержку функций с помощью SDK плагинов
// TODO добавить поддержку SQL функций, определенных пользователем, должна быть поддержка CREATE FUNCTION И CREATE PROCEDURE
func ExecuteFunction(name string, args []interface{}) (*table.ExecuteResult, error) {
	return ExecuteFunctionWithContext(context.Background(), name, args)
}

// ExecuteFunctionWithContext executes a function with a context
func ExecuteFunctionWithContext(ctx context.Context, name string, args []interface{}) (*table.ExecuteResult, error) {
	logger.Debugf("Executing function: %s with args: %v", name, args)

	// Check for user-defined functions first
	if functionRegistry != nil {
		if fn, exists := functionRegistry.GetByName(strings.ToUpper(name)); exists {
			result, err := fn.Execute(ctx, args...)
			if err != nil {
				return nil, err
			}
			logger.Debugf("User-defined function %s result: %v", name, result)
			return &table.ExecuteResult{
				Rows: [][]interface{}{{result}},
				Columns: []table.TableColumn{
					{Name: name, Type: fn.GetFunction().ProRetType.ToColType()},
				},
			}, nil
		} else {
			logger.Debugf("User-defined function %s not found", strings.ToUpper(name))
		}
	} else {
		logger.Debugf("Function registry is nil")
	}

	// Fall back to built-in functions
	switch strings.ToUpper(name) {
	// TODO добавить поддержку array, чтобы можно было вернуть name[] а не просто name или varchar
	// case "CURRENT_SCHEMAS":
	// 	// current_schemas(false) returns {public}
	// 	var result interface{}
	// 	if len(args) > 0 {
	// 		if includeImplicit, ok := args[0].(bool); ok && !includeImplicit {
	// 			result = []string{"public"}
	// 		} else {
	// 			result = []string{"public"}
	// 		}
	// 	} else {
	// 		result = []string{"public"}
	// 	}
	// 	logger.Debugf("Function %s result: %v", name, result)
	// 	return &table.ExecuteResult{
	// 		Rows: [][]interface{}{{result}},
	// 		Columns: []table.TableColumn{
	// 			{Name: "current_schemas", Type: table.ColTypeString},
	// 		},
	// 	}, nil
	// TODO лучше захардкодить EXTRACT не как обычная функция работает
	case "EXTRACT":
		if len(args) != 2 {
			return nil, fmt.Errorf("extract requires two arguments")
		}
		field, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("extract field must be string")
		}
		field = strings.ToUpper(field)
		source := args[1]
		if field == "EPOCH" {
			// Try to parse source as timestamp
			if _, ok := source.(string); ok {
				// Parse timestamp, for simplicity assume format
				// Return fixed unix for now
				result := float64(1672531200) // 2023-01-01 00:00:00 UTC
				logger.Debugf("Function %s result: %v", name, result)
				return &table.ExecuteResult{
					Rows: [][]interface{}{{result}},
					Columns: []table.TableColumn{
						{Name: "extract", Type: table.ColTypeFloat},
					},
				}, nil
			}
			return nil, fmt.Errorf("extract epoch source must be timestamp")
		}
		return nil, fmt.Errorf("unsupported extract field: %s", field)
	default:
		logger.Debugf("Unknown function: %s", name)
		return nil, fmt.Errorf("unknown function: %s", name)
	}
}

// ExecuteAggregateFunction executes aggregate functions over groups of values
func ExecuteAggregateFunction(name string, args []interface{}) (*table.ExecuteResult, error) {
	logger.Debugf("Executing aggregate function: %s with args: %v", name, args)

	// Проверяем SDK registry сначала
	if functionRegistry != nil {
		// Try to pick correct overload based on argument element types (for aggregates args are slices)
		upper := strings.ToUpper(name)
		// If there is an overload matching argument types, prefer it
		// Build argTypes: for aggregates, args are slices of values (one slice per argument)
		argOIDs := make([]ptypes.OID, 0, len(args))
		for _, a := range args {
			// each a is expected to be []interface{}
			var detected ptypes.OID = ptypes.PTypeText
			if slice, ok := a.([]interface{}); ok {
				// find first non-nil value to detect type
				for _, v := range slice {
					if v == nil {
						continue
					}
					switch v.(type) {
					case int, int32, int64:
						detected = ptypes.PTypeInt4
					case float32, float64:
						detected = ptypes.PTypeFloat8
					case bool:
						detected = ptypes.PTypeBool
					case string:
						detected = ptypes.PTypeText
					default:
						detected = ptypes.PTypeText
					}
					break
				}
			} else {
				// Not a slice - fall back to inspect value directly
				if a != nil {
					switch a.(type) {
					case int, int32, int64:
						detected = ptypes.PTypeInt4
					case float32, float64:
						detected = ptypes.PTypeFloat8
					case bool:
						detected = ptypes.PTypeBool
					case string:
						detected = ptypes.PTypeText
					default:
						detected = ptypes.PTypeText
					}
				}
			}
			argOIDs = append(argOIDs, detected)
		}

		// Try to find exact overload by arg types
		if len(argOIDs) > 0 {
			if fnMatch, ok := functionRegistry.GetByNameAndArgTypes(upper, argOIDs); ok {
				if !fnMatch.GetFunction().IsAggregate {
					return nil, fmt.Errorf("function %s is not an aggregate function", name)
				}
				result, err := fnMatch.Execute(context.Background(), args...)
				if err != nil {
					return nil, err
				}
				colType := fnMatch.GetFunction().ProRetType.ToColType()
				return &table.ExecuteResult{
					Rows:    [][]interface{}{{result}},
					Columns: []table.TableColumn{{Name: strings.ToLower(name), Type: colType}},
				}, nil
			}
		}

		// Fallback: return first function by name
		if fn, exists := functionRegistry.GetByName(upper); exists {
			if !fn.GetFunction().IsAggregate {
				return nil, fmt.Errorf("function %s is not an aggregate function", name)
			}
			result, err := fn.Execute(context.Background(), args...)
			if err != nil {
				return nil, err
			}
			logger.Debugf("Aggregate function %s result: %v", name, result)
			colType := fn.GetFunction().ProRetType.ToColType()
			return &table.ExecuteResult{
				Rows:    [][]interface{}{{result}},
				Columns: []table.TableColumn{{Name: strings.ToLower(name), Type: colType}},
			}, nil
		}
	}

	// Если функция не найдена в registry
	return nil, fmt.Errorf("unknown aggregate function: %s", name)
}

// GetRegisteredFunction returns a registered function by name (for wire/type inference)
func GetRegisteredFunction(name string) (psdk.IFunction, bool) {
	if functionRegistry == nil {
		return nil, false
	}
	return functionRegistry.GetByName(strings.ToUpper(name))
}
