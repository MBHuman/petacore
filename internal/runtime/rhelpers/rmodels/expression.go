package rmodels

import (
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
)

type ExpressionType int

const (
	ExpressionTypeUnknown ExpressionType = iota
	ExpressionTypeBool
	ExpressionTypeResultRows
	ExpressionTypeCase
	ExpressionTypeSubquery
	ExpressionTypeParamRef
)

func (et *ExpressionType) String() string {
	switch *et {
	case ExpressionTypeBool:
		return "bool"
	case ExpressionTypeUnknown:
		return "unknown"
	case ExpressionTypeSubquery:
		return "subquery"
	case ExpressionTypeParamRef:
		return "paramref"
	default:
		panic("unexpected rparser.ExpressionType")
	}
}

// ParamRefExpression — ссылка на параметр (scalar subquery)
type ParamRefExpression struct {
	Index         int
	RuntimeParams map[int]interface{} // ссылка на runtime параметры для получения значения
}

func (p *ParamRefExpression) Type() ExpressionType {
	return ExpressionTypeParamRef
}

func (p *ParamRefExpression) Value() interface{} {
	if p.RuntimeParams == nil {
		return nil
	}
	return p.RuntimeParams[p.Index]
}

type SubqueryExpression struct {
	Select *statements.SelectStatement
}

func (sq *SubqueryExpression) Type() ExpressionType {
	return ExpressionTypeSubquery
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
