package rops

import (
	"fmt"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
)

type numericBinOp func(pmem.Allocator, ptypes.NumericType[any], ptypes.NumericType[any]) (ptypes.NumericType[any], error)

func applyArithmetic(
	allocator pmem.Allocator,
	a, b *ptypes.Row,
	aSchema, bSchema *serializers.BaseSchema,
	opName string,
	op numericBinOp,
) (*rmodels.ResultRowsExpression, error) {
	aVal, aOid, err := aSchema.GetField(a, 0)
	if err != nil {
		return nil, err
	}
	aDesVal, err := serializers.DeserializeGeneric(aVal, aOid)
	if err != nil {
		return nil, err
	}

	bVal, bOid, err := bSchema.GetField(b, 0)
	if err != nil {
		return nil, err
	}
	bDesVal, err := serializers.DeserializeGeneric(bVal, bOid)
	if err != nil {
		return nil, err
	}

	// приводим типы если они разные (например int4 + float8)
	aDesVal, bDesVal, resultOid, err := coerceNumericTypes(allocator, aDesVal, aOid, bDesVal, bOid)
	if err != nil {
		return nil, fmt.Errorf("%s: type coercion failed: %w", opName, err)
	}

	aN, aOk := extractNumericAny(aDesVal, resultOid)
	bN, bOk := extractNumericAny(bDesVal, resultOid)
	if !aOk || !bOk {
		return nil, fmt.Errorf("%s: both values must be numeric (OID %d, %d)", opName, aOid, bOid)
	}

	res, err := op(allocator, aN, bN)
	if err != nil {
		return nil, err
	}

	// результирующая схема с правильным OID
	resultSchema := aSchema
	if resultOid != aOid {
		resultSchema = serializers.NewBaseSchema([]serializers.FieldDef{{
			Name: aSchema.Fields[0].Name,
			OID:  resultOid,
		}})
	}

	resultRow, err := resultSchema.Pack(allocator, [][]byte{res.GetBuffer()})
	if err != nil {
		return nil, err
	}

	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:   []*ptypes.Row{resultRow},
			Schema: resultSchema,
		},
	}, nil
}

// extractNumericAny извлекает NumericType[any] по известному OID
func extractNumericAny(val ptypes.BaseType[any], oid ptypes.OID) (ptypes.NumericType[any], bool) {
	switch oid {
	case ptypes.PTypeInt2:
		n, ok := ptypes.TryIntoNumeric[int16](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	case ptypes.PTypeInt4:
		n, ok := ptypes.TryIntoNumeric[int32](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	case ptypes.PTypeInt8:
		n, ok := ptypes.TryIntoNumeric[int64](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	case ptypes.PTypeFloat4:
		n, ok := ptypes.TryIntoNumeric[float32](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	case ptypes.PTypeFloat8:
		n, ok := ptypes.TryIntoNumeric[float64](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	case ptypes.PTypeNumeric:
		n, ok := ptypes.TryIntoNumeric[[]byte](val)
		if !ok {
			return nil, false
		}
		return ptypes.NewAnyWrapper(n), true
	}
	return nil, false
}

// coerceNumericTypes приводит два значения к общему типу
// правила: int2 < int4 < int8 < float4 < float8 < numeric
func coerceNumericTypes(
	allocator pmem.Allocator,
	a ptypes.BaseType[any], aOid ptypes.OID,
	b ptypes.BaseType[any], bOid ptypes.OID,
) (ptypes.BaseType[any], ptypes.BaseType[any], ptypes.OID, error) {
	if aOid == bOid {
		return a, b, aOid, nil
	}

	targetOid := widenOID(aOid, bOid)

	if aOid != targetOid {
		casted, err := castToOID(allocator, a, aOid, targetOid)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cast a from %d to %d: %w", aOid, targetOid, err)
		}
		a = casted
	}
	if bOid != targetOid {
		casted, err := castToOID(allocator, b, bOid, targetOid)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("cast b from %d to %d: %w", bOid, targetOid, err)
		}
		b = casted
	}

	return a, b, targetOid, nil
}

// widenOID возвращает более широкий из двух числовых OID
func widenOID(a, b ptypes.OID) ptypes.OID {
	order := map[ptypes.OID]int{
		ptypes.PTypeInt2:    1,
		ptypes.PTypeInt4:    2,
		ptypes.PTypeInt8:    3,
		ptypes.PTypeFloat4:  4,
		ptypes.PTypeFloat8:  5,
		ptypes.PTypeNumeric: 6,
	}
	if order[a] >= order[b] {
		return a
	}
	return b
}

// castToOID кастует значение к целевому OID через CastableType
func castToOID(allocator pmem.Allocator, val ptypes.BaseType[any], srcOid, dstOid ptypes.OID) (ptypes.BaseType[any], error) {
	switch srcOid {
	case ptypes.PTypeInt2:
		w, ok := val.(ptypes.AnyWrapper[int16])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[int16]")
		}
		c, ok := w.Inner().(ptypes.CastableType[int16])
		if !ok {
			return nil, fmt.Errorf("int2 does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	case ptypes.PTypeInt4:
		w, ok := val.(ptypes.AnyWrapper[int32])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[int32]")
		}
		c, ok := w.Inner().(ptypes.CastableType[int32])
		if !ok {
			return nil, fmt.Errorf("int4 does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	case ptypes.PTypeInt8:
		w, ok := val.(ptypes.AnyWrapper[int64])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[int64]")
		}
		c, ok := w.Inner().(ptypes.CastableType[int64])
		if !ok {
			return nil, fmt.Errorf("int8 does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	case ptypes.PTypeFloat4:
		w, ok := val.(ptypes.AnyWrapper[float32])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[float32]")
		}
		c, ok := w.Inner().(ptypes.CastableType[float32])
		if !ok {
			return nil, fmt.Errorf("float4 does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	case ptypes.PTypeFloat8:
		w, ok := val.(ptypes.AnyWrapper[float64])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[float64]")
		}
		c, ok := w.Inner().(ptypes.CastableType[float64])
		if !ok {
			return nil, fmt.Errorf("float8 does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	case ptypes.PTypeNumeric:
		w, ok := val.(ptypes.AnyWrapper[[]byte])
		if !ok {
			return nil, fmt.Errorf("expected AnyWrapper[[]byte]")
		}
		c, ok := w.Inner().(ptypes.CastableType[[]byte])
		if !ok {
			return nil, fmt.Errorf("numeric does not implement CastableType")
		}
		return c.CastTo(allocator, dstOid)
	}
	return nil, fmt.Errorf("castToOID: unsupported src OID %d", srcOid)
}

func AddValues(allocator pmem.Allocator, a, b *ptypes.Row, aSchema, bSchema *serializers.BaseSchema) (*rmodels.ResultRowsExpression, error) {
	return applyArithmetic(allocator, a, b, aSchema, bSchema, "AddValues",
		func(alloc pmem.Allocator, x, y ptypes.NumericType[any]) (ptypes.NumericType[any], error) {
			return x.Add(alloc, y)
		},
	)
}

func SubtractValues(allocator pmem.Allocator, a, b *ptypes.Row, aSchema, bSchema *serializers.BaseSchema) (*rmodels.ResultRowsExpression, error) {
	return applyArithmetic(allocator, a, b, aSchema, bSchema, "SubtractValues",
		func(alloc pmem.Allocator, x, y ptypes.NumericType[any]) (ptypes.NumericType[any], error) {
			return x.Sub(alloc, y)
		},
	)
}

func MultiplyValues(allocator pmem.Allocator, a, b *ptypes.Row, aSchema, bSchema *serializers.BaseSchema) (*rmodels.ResultRowsExpression, error) {
	return applyArithmetic(allocator, a, b, aSchema, bSchema, "MultiplyValues",
		func(alloc pmem.Allocator, x, y ptypes.NumericType[any]) (ptypes.NumericType[any], error) {
			return x.Mul(alloc, y)
		},
	)
}

func DivideValues(allocator pmem.Allocator, a, b *ptypes.Row, aSchema, bSchema *serializers.BaseSchema) (*rmodels.ResultRowsExpression, error) {
	return applyArithmetic(allocator, a, b, aSchema, bSchema, "DivideValues",
		func(alloc pmem.Allocator, x, y ptypes.NumericType[any]) (ptypes.NumericType[any], error) {
			return x.Div(alloc, y)
		},
	)
}

func ModuloValues(allocator pmem.Allocator, a, b *ptypes.Row, aSchema, bSchema *serializers.BaseSchema) (*rmodels.ResultRowsExpression, error) {
	return applyArithmetic(allocator, a, b, aSchema, bSchema, "ModuloValues",
		func(alloc pmem.Allocator, x, y ptypes.NumericType[any]) (ptypes.NumericType[any], error) {
			return x.Mod(alloc, y)
		},
	)
}
