package functions

import (
	"context"
	"fmt"
	"math"
	"petacore/internal/logger"
	"petacore/internal/runtime/rsql/table"
	psdk "petacore/sdk"
	"strconv"
	"strings"
)

var functionRegistry *psdk.FunctionRegistry

// SetFunctionRegistry sets the global function registry for user-defined functions
func SetFunctionRegistry(registry *psdk.FunctionRegistry) {
	functionRegistry = registry
}

// ExecuteFunction executes a built-in or user-defined function
// TODO добавить поддержку функций с помощью SDK плагинов
// TODO добавить поддержку SQL функций, определенных пользователем, должна быть поддержка CREATE FUNCTION И CREATE PROCEDURE
func ExecuteFunction(name string, args []interface{}) (*table.ExecuteResult, error) {
	logger.Debugf("Executing function: %s with args: %v", name, args)

	// Check for user-defined functions first
	if functionRegistry != nil {
		if fn, exists := functionRegistry.GetByName(strings.ToUpper(name)); exists {
			result, err := fn.Execute(context.Background(), args...)
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
	case "CURRENT_DATABASE":
		result := "testdb"
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]any{{result}},
			Columns: []table.TableColumn{
				{Name: "current_database", Type: table.ColTypeString},
			},
		}, nil
	case "CURRENT_SCHEMA":
		result := "public"
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]any{{result}},
			Columns: []table.TableColumn{
				{Name: "current_schema", Type: table.ColTypeString},
			},
		}, nil
	case "CURRENT_USER":
		result := "postgres"
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]any{{result}},
			Columns: []table.TableColumn{
				{Name: "current_user", Type: table.ColTypeString},
			},
		}, nil
	case "CURRENT_SCHEMAS":
		// current_schemas(false) returns {public}
		var result interface{}
		if len(args) > 0 {
			if includeImplicit, ok := args[0].(bool); ok && !includeImplicit {
				result = []string{"public"}
			} else {
				result = []string{"public"}
			}
		} else {
			result = []string{"public"}
		}
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "current_schemas", Type: table.ColTypeString},
			},
		}, nil
	case "VERSION":
		result := "PostgreSQL 14.0 (PetaCore)"
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "version", Type: table.ColTypeString},
			},
		}, nil
	case "PG_POSTMASTER_START_TIME":
		// Return a fixed timestamp for simplicity
		result := "2023-01-01 00:00:00+00"
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "pg_postmaster_start_time", Type: table.ColTypeString},
			},
		}, nil
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
	case "ROUND":
		if len(args) == 0 {
			return nil, fmt.Errorf("round requires at least one argument")
		}
		logger.Debugf("ROUND args[0] type: %T, value: %v", args[0], args[0])
		var val float64
		switch v := args[0].(type) {
		case float64:
			val = v
		case int:
			val = float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				val = f
			} else {
				return nil, fmt.Errorf("round argument must be numeric")
			}
		default:
			return nil, fmt.Errorf("round argument must be numeric")
		}
		// For simplicity, round to nearest integer
		result := int(math.Round(val))
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "round", Type: table.ColTypeInt},
			},
		}, nil
	case "MIN":
		if len(args) == 0 {
			return nil, fmt.Errorf("min requires at least one argument")
		}
		minVal := args[0]
		for _, arg := range args[1:] {
			switch v := arg.(type) {
			case int:
				if vi, ok := minVal.(int); ok && v < vi {
					minVal = v
				}
			case float64:
				if vf, ok := minVal.(float64); ok && v < vf {
					minVal = v
				} else if vi, ok := minVal.(int); ok && float64(vi) > v {
					minVal = v
				}
			case string:
				if vs, ok := minVal.(string); ok && v < vs {
					minVal = v
				}
			}
		}
		logger.Debugf("Function %s result: %v", name, minVal)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{minVal}},
			Columns: []table.TableColumn{
				{Name: "min", Type: table.ColTypeString},
			},
		}, nil
	case "MAX":
		if len(args) == 0 {
			return nil, fmt.Errorf("max requires at least one argument")
		}
		maxVal := args[0]
		for _, arg := range args[1:] {
			switch v := arg.(type) {
			case int:
				if vi, ok := maxVal.(int); ok && v > vi {
					maxVal = v
				}
			case float64:
				if vf, ok := maxVal.(float64); ok && v > vf {
					maxVal = v
				} else if vi, ok := maxVal.(int); ok && float64(vi) < v {
					maxVal = v
				}
			case string:
				if vs, ok := maxVal.(string); ok && v > vs {
					maxVal = v
				}
			}
		}
		logger.Debugf("Function %s result: %v", name, maxVal)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{maxVal}},
			Columns: []table.TableColumn{
				{Name: "max", Type: table.ColTypeString},
			},
		}, nil
	case "PG_TABLE_IS_VISIBLE":
		// pg_table_is_visible(oid) - check if table is visible in search path
		// For simplicity, always return true
		if len(args) != 1 {
			return nil, fmt.Errorf("pg_table_is_visible requires exactly one argument")
		}
		result := true
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "pg_table_is_visible", Type: table.ColTypeBool},
			},
		}, nil

	case "PG_GET_USERBYID":
		// pg_get_userbyid(uid) - get username by user ID
		if len(args) != 1 {
			return nil, fmt.Errorf("pg_get_userbyid requires exactly one argument")
		}
		var uid int
		switch v := args[0].(type) {
		case int:
			uid = v
		case int32:
			uid = int(v)
		case int64:
			uid = int(v)
		case float64:
			uid = int(v)
		default:
			return nil, fmt.Errorf("pg_get_userbyid argument must be integer")
		}
		// For simplicity, return "postgres" for uid 10, else "user<uid>"
		var result string
		if uid == 10 {
			result = "postgres"
		} else {
			result = fmt.Sprintf("user%d", uid)
		}
		logger.Debugf("Function %s result: %v", name, result)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{result}},
			Columns: []table.TableColumn{
				{Name: "pg_get_userbyid", Type: table.ColTypeString},
			},
		}, nil
	default:
		logger.Debugf("Unknown function: %s", name)
		return nil, fmt.Errorf("unknown function: %s", name)
	}
}

// ExecuteAggregateFunction executes aggregate functions over groups of values
func ExecuteAggregateFunction(name string, args []interface{}) (*table.ExecuteResult, error) {
	logger.Debugf("Executing aggregate function: %s with args: %v", name, args)
	switch strings.ToUpper(name) {
	case "MAX":
		if len(args) == 0 {
			return nil, fmt.Errorf("max requires at least one argument")
		}
		values, ok := args[0].([]interface{})
		if !ok {
			return nil, fmt.Errorf("max argument must be a slice of values")
		}
		if len(values) == 0 {
			return nil, nil
		}
		maxVal := values[0]
		colType := table.ColTypeString
		// Determine type based on first non-nil value
		for _, val := range values {
			if val == nil {
				continue
			}
			switch val.(type) {
			case int:
				colType = table.ColTypeInt
			case float64:
				colType = table.ColTypeFloat
			case string:
				colType = table.ColTypeString
			}
			break // Only check first non-nil value for type
		}
		// Now find the actual max
		for _, val := range values {
			if val == nil {
				continue
			}
			switch v := val.(type) {
			case int:
				if maxVal == nil {
					maxVal = v
				} else if vi, ok := maxVal.(int); ok && v > vi {
					maxVal = v
				}
			case float64:
				if maxVal == nil {
					maxVal = v
				} else if vf, ok := maxVal.(float64); ok && v > vf {
					maxVal = v
				} else if vi, ok := maxVal.(int); ok && float64(vi) < v {
					maxVal = v
				}
			case string:
				if maxVal == nil {
					maxVal = v
				} else if vs, ok := maxVal.(string); ok && v > vs {
					maxVal = v
				}
			}
		}
		logger.Debugf("Aggregate function %s result: %v", name, maxVal)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{maxVal}},
			Columns: []table.TableColumn{
				{Name: "max", Type: colType},
			},
		}, nil
	case "MIN":
		if len(args) == 0 {
			return nil, fmt.Errorf("min requires at least one argument")
		}
		values, ok := args[0].([]interface{})
		if !ok {
			return nil, fmt.Errorf("min argument must be a slice of values")
		}
		if len(values) == 0 {
			return nil, nil
		}
		minVal := values[0]
		colType := table.ColTypeString
		// Determine type based on first non-nil value
		for _, val := range values {
			if val == nil {
				continue
			}
			switch val.(type) {
			case int:
				colType = table.ColTypeInt
			case float64:
				colType = table.ColTypeFloat
			case string:
				colType = table.ColTypeString
			}
			break // Only check first non-nil value for type
		}
		// Now find the actual min
		for _, val := range values {
			if val == nil {
				continue
			}
			switch v := val.(type) {
			case int:
				if minVal == nil {
					minVal = v
				} else if vi, ok := minVal.(int); ok && v < vi {
					minVal = v
				}
			case float64:
				if minVal == nil {
					minVal = v
				} else if vf, ok := minVal.(float64); ok && v < vf {
					minVal = v
				} else if vi, ok := minVal.(int); ok && float64(vi) > v {
					minVal = v
				}
			case string:
				if minVal == nil {
					minVal = v
				} else if vs, ok := minVal.(string); ok && v < vs {
					minVal = v
				}
			}
		}
		logger.Debugf("Aggregate function %s result: %v", name, minVal)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{minVal}},
			Columns: []table.TableColumn{
				{Name: "min", Type: colType},
			},
		}, nil
	case "COUNT":
		if len(args) == 0 {
			return nil, fmt.Errorf("count requires at least one argument")
		}
		values, ok := args[0].([]interface{})
		if !ok {
			return nil, fmt.Errorf("count argument must be a slice of values")
		}
		count := 0
		for _, val := range values {
			if val != nil {
				count++
			}
		}
		colType := table.ColTypeFloat
		logger.Debugf("Aggregate function %s result: %v", name, count)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{float64(count)}},
			Columns: []table.TableColumn{
				{Name: "count", Type: colType},
			},
		}, nil
	case "SUM":
		if len(args) == 0 {
			return nil, fmt.Errorf("sum requires at least one argument")
		}
		values, ok := args[0].([]interface{})
		if !ok {
			return nil, fmt.Errorf("sum argument must be a slice of values")
		}
		var sum float64
		for _, val := range values {
			if val == nil {
				continue
			}
			switch v := val.(type) {
			case int:
				sum += float64(v)
			case float64:
				sum += v
			default:
				return nil, fmt.Errorf("sum can only operate on numeric values")
			}
		}
		logger.Debugf("Aggregate function %s result: %v", name, sum)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{sum}},
			Columns: []table.TableColumn{
				{Name: "sum", Type: table.ColTypeFloat},
			},
		}, nil
	case "AVG":
		if len(args) == 0 {
			return nil, fmt.Errorf("avg requires at least one argument")
		}
		values, ok := args[0].([]interface{})
		if !ok {
			return nil, fmt.Errorf("avg argument must be a slice of values")
		}
		var sum float64
		count := 0
		colType := table.ColTypeFloat
		for _, val := range values {
			if val == nil {
				continue
			}
			switch v := val.(type) {
			case int:
				sum += float64(v)
			case float64:
				sum += v
			default:
				return nil, fmt.Errorf("avg can only operate on numeric values")
			}
			count++
		}
		if count == 0 {
			return nil, nil
		}
		avg := sum / float64(count)
		logger.Debugf("Aggregate function %s result: %v", name, avg)
		return &table.ExecuteResult{
			Rows: [][]interface{}{{avg}},
			Columns: []table.TableColumn{
				{Name: "avg", Type: colType},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown aggregate function: %s", name)
	}
}
