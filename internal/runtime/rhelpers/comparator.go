package rhelpers

import (
	"fmt"
	"petacore/internal/runtime/rsql/items"
	"strings"
)

// matchesWhere checks if a row matches the WHERE clause condition
func matchesWhere(row map[string]interface{}, where *items.WhereClause) bool {
	fieldValue, exists := row[where.Field]
	if !exists {
		return false
	}

	switch where.Operator {
	case "=":
		return compareValues(fieldValue, where.Value) == 0
	case "!=":
		return compareValues(fieldValue, where.Value) != 0
	case "<>":
		return compareValues(fieldValue, where.Value) != 0
	case ">":
		return compareValues(fieldValue, where.Value) > 0
	case "<":
		return compareValues(fieldValue, where.Value) < 0
	case ">=":
		return compareValues(fieldValue, where.Value) >= 0
	case "<=":
		return compareValues(fieldValue, where.Value) <= 0
	default:
		return false
	}
}

// compareValues compares two values for sorting
func compareValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	switch va := a.(type) {
	case int:
		if vb, ok := b.(int); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case int64:
		if vb, ok := b.(int64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case float64:
		if vb, ok := b.(float64); ok {
			if va < vb {
				return -1
			} else if va > vb {
				return 1
			}
			return 0
		}
	case string:
		if vb, ok := b.(string); ok {
			return strings.Compare(va, vb)
		}
	}
	// Default: convert to string
	sa := fmt.Sprintf("%v", a)
	sb := fmt.Sprintf("%v", b)
	return strings.Compare(sa, sb)
}
