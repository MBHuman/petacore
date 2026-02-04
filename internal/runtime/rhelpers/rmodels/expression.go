package rmodels

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/table"
)

type ExpressionType int

const (
	ExpressionTypeUnknown ExpressionType = iota
	ExpressionTypeBool
	ExpressionTypeResultRows
	ExpressionTypeCase
)

func (et *ExpressionType) String() string {
	switch *et {
	case ExpressionTypeBool:
		return "bool"
	case ExpressionTypeUnknown:
		return "unknown"
	default:
		panic("unexpected rparser.ExpressionType")
	}
}

type Expression interface {
	Type() ExpressionType
}

type BoolExpression struct {
	Value bool
}

func (le *BoolExpression) Type() ExpressionType {
	return ExpressionTypeBool
}

type ResultRowsExpression struct {
	Row *table.ExecuteResult
}

func (rre *ResultRowsExpression) Type() ExpressionType {
	return ExpressionTypeResultRows
}

type CaseExpression struct {
	Context parser.ICaseExpressionContext
}

func (c *CaseExpression) Type() ExpressionType {
	return ExpressionTypeCase
}
