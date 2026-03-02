package rops

import (
	"fmt"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"time"
)

func CastValue(
	allocator pmem.Allocator,
	expression *rmodels.ResultRowsExpression,
	targetOID ptypes.OID,
) (rmodels.Expression, error) {
	if len(expression.Row.Rows) == 0 {
		return nil, fmt.Errorf("cast: empty result")
	}
	if len(expression.Row.Schema.Fields) == 0 {
		return nil, fmt.Errorf("cast: empty schema")
	}

	// берём первое поле первой строки
	buf, srcOID, err := expression.Row.Schema.GetField(expression.Row.Rows[0], 0)
	if err != nil {
		return nil, fmt.Errorf("cast: get field: %w", err)
	}

	// десериализуем в BaseType[any]
	srcVal, err := serializers.DeserializeGeneric(buf, srcOID)
	if err != nil {
		return nil, fmt.Errorf("cast: deserialize: %w", err)
	}

	// пробуем CastableType
	castResult, err := tryCast(allocator, srcVal, srcOID, targetOID)
	if err != nil {
		return nil, fmt.Errorf("cast to OID %d: %w", targetOID, err)
	}

	// сериализуем результат обратно в буфер
	resultBuf := castResult.GetBuffer()

	// сохраняем имя колонки
	colName := expression.Row.Schema.Fields[0].Name
	if colName == "" {
		colName = "?column?"
	}

	resultSchema := serializers.NewBaseSchema([]serializers.FieldDef{{
		Name: colName,
		OID:  targetOID,
	}})

	resultRow, err := resultSchema.Pack(allocator, [][]byte{resultBuf})
	if err != nil {
		return nil, fmt.Errorf("cast: pack result: %w", err)
	}

	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:   []*ptypes.Row{resultRow},
			Schema: resultSchema,
		},
	}, nil
}

// tryCast выполняет каст через CastableType если тип его поддерживает
func tryCast(
	allocator pmem.Allocator,
	val ptypes.BaseType[any],
	srcOID ptypes.OID,
	targetOID ptypes.OID,
) (ptypes.BaseType[any], error) {
	// извлекаем inner тип из AnyWrapper и пробуем CastableType
	switch srcOID {
	case ptypes.PTypeBool:
		if w, ok := val.(ptypes.AnyWrapper[bool]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[bool]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeInt2:
		if w, ok := val.(ptypes.AnyWrapper[int16]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[int16]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeInt4:
		if w, ok := val.(ptypes.AnyWrapper[int32]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[int32]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeInt8:
		if w, ok := val.(ptypes.AnyWrapper[int64]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[int64]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeFloat4:
		if w, ok := val.(ptypes.AnyWrapper[float32]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[float32]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeFloat8:
		if w, ok := val.(ptypes.AnyWrapper[float64]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[float64]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeNumeric:
		if w, ok := val.(ptypes.AnyWrapper[[]byte]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[[]byte]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeText:
		if w, ok := val.(ptypes.AnyWrapper[string]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[string]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeVarchar:
		if w, ok := val.(ptypes.AnyWrapper[string]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[string]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	case ptypes.PTypeDate, ptypes.PTypeTime, ptypes.PTypeTimestamp, ptypes.PTypeTimestampz:
		if w, ok := val.(ptypes.AnyWrapper[*time.Time]); ok {
			if c, ok := w.Inner().(ptypes.CastableType[*time.Time]); ok {
				return c.CastTo(allocator, targetOID)
			}
		}
	}

	return nil, fmt.Errorf("type OID %d does not support cast to OID %d", srcOID, targetOID)
}
