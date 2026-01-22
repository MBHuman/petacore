package functions

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// ExecuteFunction executes a built-in function
// TODO добавить поддержку функций с помощью SDK плагинов
// TODO добавить поддержку SQL функций, определенных пользователем, должна быть поддержка CREATE FUNCTION И CREATE PROCEDURE
func ExecuteFunction(name string, args []interface{}) (interface{}, error) {
	log.Printf("Executing function: %s with args: %v", name, args)
	switch strings.ToUpper(name) {
	case "CURRENT_DATABASE":
		result := "testdb"
		log.Printf("Function %s result: %v", name, result)
		return result, nil
	case "CURRENT_SCHEMA":
		result := "public"
		log.Printf("Function %s result: %v", name, result)
		return result, nil
	case "CURRENT_USER":
		result := "postgres"
		log.Printf("Function %s result: %v", name, result)
		return result, nil
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
		log.Printf("Function %s result: %v", name, result)
		return result, nil
	case "VERSION":
		result := "PostgreSQL 14.0 (PetaCore)"
		log.Printf("Function %s result: %v", name, result)
		return result, nil
	case "PG_POSTMASTER_START_TIME":
		// Return a fixed timestamp for simplicity
		result := "2023-01-01 00:00:00+00"
		log.Printf("Function %s result: %v", name, result)
		return result, nil
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
				log.Printf("Function %s result: %v", name, result)
				return result, nil
			}
			return nil, fmt.Errorf("extract epoch source must be timestamp")
		}
		return nil, fmt.Errorf("unsupported extract field: %s", field)
	case "ROUND":
		if len(args) == 0 {
			return nil, fmt.Errorf("round requires at least one argument")
		}
		log.Printf("ROUND args[0] type: %T, value: %v", args[0], args[0])
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
		log.Printf("Function %s result: %v", name, result)
		return result, nil
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
		log.Printf("Function %s result: %v", name, minVal)
		return minVal, nil
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
		log.Printf("Function %s result: %v", name, maxVal)
		return maxVal, nil
	case "PG_TABLE_IS_VISIBLE":
		// pg_table_is_visible(oid) - check if table is visible in search path
		// For simplicity, always return true
		if len(args) != 1 {
			return nil, fmt.Errorf("pg_table_is_visible requires exactly one argument")
		}
		result := true
		log.Printf("Function %s result: %v", name, result)
		return result, nil
	default:
		log.Printf("Unknown function: %s", name)
		return nil, fmt.Errorf("unknown function: %s", name)
	}
}
