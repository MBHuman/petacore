package rparser

import (
	"context"
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rhelpers/subquery"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strconv"
	"strings"
)

func ParseConcatExpression(allocator pmem.Allocator, ctx context.Context, concatExpr parser.IConcatExpressionContext, row *table.ResultRow, subExec subquery.SubqueryExecutor) (rmodels.Expression, error) {
	if concatExpr == nil {
		return nil, nil
	}

	exprs := concatExpr.AllAdditiveExpression()
	if len(exprs) == 0 {
		return nil, nil
	}

	if len(exprs) == 1 && len(concatExpr.AllCONCAT()) == 0 {
		return ParseAdditiveExpression(allocator, ctx, exprs[0], row, subExec)
	}

	fields := []serializers.FieldDef{{
		Name: "concat_result",
		OID:  ptypes.PTypeText,
	}}
	resultSchema := serializers.NewBaseSchema(fields)

	// накапливаем результат как string — проще и надёжнее
	var accumulated strings.Builder

	for _, e := range exprs {
		val, err := ParseAdditiveExpression(allocator, ctx, e, row, subExec)
		if err != nil {
			return nil, err
		}

		valExpr, ok := val.(*rmodels.ResultRowsExpression)
		if !ok {
			return nil, fmt.Errorf("[ParseConcatExpression] concatenation supports only result row expressions")
		}

		if len(valExpr.Row.Rows) != 1 || len(valExpr.Row.Schema.Fields) != 1 {
			return nil, fmt.Errorf("[ParseConcatExpression] expected single value for concatenation, got %d rows and %d fields",
				len(valExpr.Row.Rows), len(valExpr.Row.Schema.Fields))
		}

		valField, oid, err := valExpr.Row.Schema.GetField(valExpr.Row.Rows[0], 0)
		if err != nil {
			return nil, fmt.Errorf("[ParseConcatExpression] failed to get field: %v", err)
		}

		// конвертируем любой тип в строку через каст
		str, err := coerceToString(allocator, valField, oid)
		if err != nil {
			return nil, fmt.Errorf("[ParseConcatExpression] failed to coerce value to string: %v", err)
		}

		accumulated.WriteString(str)
	}

	finalStr := accumulated.String()
	logger.Debugf("Concatenated result: %s\n", finalStr)

	resultBuf, err := serializers.TextSerializerInstance.Serialize(allocator, finalStr)
	if err != nil {
		return nil, fmt.Errorf("[ParseConcatExpression] failed to serialize result: %v", err)
	}

	resultRow, err := resultSchema.Pack(allocator, [][]byte{resultBuf})
	if err != nil {
		return nil, fmt.Errorf("[ParseConcatExpression] failed to pack result: %v", err)
	}

	return &rmodels.ResultRowsExpression{
		Row: &table.ExecuteResult{
			Rows:   []*ptypes.Row{resultRow},
			Schema: resultSchema,
		},
	}, nil
}

// coerceToString конвертирует любое значение в строку для конкатенации
func coerceToString(allocator pmem.Allocator, buf []byte, oid ptypes.OID) (string, error) {
	switch oid {
	case ptypes.PTypeText, ptypes.PTypeVarchar:
		return string(buf), nil

	case ptypes.PTypeInt2:
		v, err := serializers.Int2SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt4:
		v, err := serializers.Int4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(v.IntoGo()), 10), nil

	case ptypes.PTypeInt8:
		v, err := serializers.Int8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(v.IntoGo(), 10), nil

	case ptypes.PTypeFloat4:
		v, err := serializers.Float4SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(float64(v.IntoGo()), 'f', -1, 32), nil

	case ptypes.PTypeFloat8:
		v, err := serializers.Float8SerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(v.IntoGo(), 'f', -1, 64), nil

	// TODO добавить поддержку, надо чтобы в строку переводился
	// case ptypes.PTypeNumeric:
	// 	v, err := serializers.NumericSerializerInstance.Deserialize(buf)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	f, err := v.
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	return f.Text('f', int(v.Meta.Scale)), nil

	case ptypes.PTypeBool:
		v, err := serializers.BoolSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		if v.IntoGo() {
			return "true", nil
		}
		return "false", nil

	case ptypes.PTypeDate:
		v, err := serializers.DateSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil date")
		}
		return tm.Format("2006-01-02"), nil

	case ptypes.PTypeTimestamp:
		v, err := serializers.TimestampSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil timestamp")
		}
		return tm.Format("2006-01-02 15:04:05"), nil

	case ptypes.PTypeTimestampz:
		v, err := serializers.TimestampzSerializerInstance.Deserialize(buf)
		if err != nil {
			return "", err
		}
		tm := v.IntoGo()
		if tm == nil {
			return "", fmt.Errorf("coerceToString: nil timestampz")
		}
		return tm.UTC().Format("2006-01-02 15:04:05+00"), nil

	default:
		return "", fmt.Errorf("coerceToString: unsupported OID %d", oid)
	}
}
