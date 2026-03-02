package rparser

import (
	"fmt"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rhelpers/rmodels"
	"petacore/internal/runtime/rsql/table"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strconv"
	"strings"
	"time"
)

func ParsePrimaryLiteral(allocator pmem.Allocator, value parser.IValueContext) (rmodels.Expression, error) {
	if value.NUMBER() != nil {
		numStr := value.NUMBER().GetText()
		if strings.Contains(numStr, ".") {
			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				fields := []serializers.FieldDef{{
					Name: "?column?",
					OID:  ptypes.PTypeFloat8,
				}}
				resultSchema := serializers.NewBaseSchema(fields)
				serVal, err := serializers.Float8SerializerInstance.Serialize(allocator, val)
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize float value: %v", err)
				}
				resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack float value: %v", err)
				}
				return &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:   []*ptypes.Row{resultRow},
						Schema: resultSchema,
					},
				}, nil
			}
		} else {
			if val, err := strconv.Atoi(numStr); err == nil {
				fields := []serializers.FieldDef{{
					Name: "?column?",
					OID:  ptypes.PTypeInt4,
				}}
				resultSchema := serializers.NewBaseSchema(fields)
				serVal, err := serializers.Int4SerializerInstance.Serialize(allocator, int32(val))
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize integer value: %v", err)
				}
				resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack integer value: %v", err)
				}

				return &rmodels.ResultRowsExpression{
					Row: &table.ExecuteResult{
						Rows:   []*ptypes.Row{resultRow},
						Schema: resultSchema,
					},
				}, nil
			}
		}
	}

	if valueLiteral := value.STRING_LITERAL(); valueLiteral != nil {
		str := valueLiteral.GetText()
		// Remove quotes
		if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
			str = str[1 : len(str)-1]
		}
		fields := []serializers.FieldDef{{
			Name: "?column?",
			OID:  ptypes.PTypeText,
		}}
		resultSchema := serializers.NewBaseSchema(fields)
		serVal, err := serializers.TextSerializerInstance.Serialize(allocator, str)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize string value: %v", err)
		}
		resultRow, err := resultSchema.Pack(allocator, [][]byte{serVal})
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack string value: %v", err)
		}
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{resultRow},
				Schema: resultSchema,
			},
		}, nil
	}

	if value.TRUE() != nil {
		fields := []serializers.FieldDef{{
			Name: "?column?",
			OID:  ptypes.PTypeBool,
		}}
		resultSchema := serializers.NewBaseSchema(fields)
		val, err := serializers.BoolSerializerInstance.Serialize(allocator, true)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize boolean value: %v", err)
		}
		resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack boolean value: %v", err)
		}
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{resultRow},
				Schema: resultSchema,
			},
		}, nil
	}

	if value.FALSE() != nil {
		fields := []serializers.FieldDef{{
			Name: "?column?",
			OID:  ptypes.PTypeBool,
		}}
		resultSchema := serializers.NewBaseSchema(fields)
		val, err := serializers.BoolSerializerInstance.Serialize(allocator, false)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize boolean value: %v", err)
		}
		resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack boolean value: %v", err)
		}
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{resultRow},
				Schema: resultSchema,
			},
		}, nil
	}

	if value.CURRENT_TIMESTAMP() != nil {
		fields := []serializers.FieldDef{{
			Name: "?column?",
			OID:  ptypes.PTypeTimestamp,
		}}
		resultSchema := serializers.NewBaseSchema(fields)
		timeVal := time.Now()
		val, err := serializers.TimestampSerializerInstance.Serialize(allocator, &timeVal)
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize timestamp value: %v", err)
		}
		resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
		if err != nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack timestamp value: %v", err)
		}
		return &rmodels.ResultRowsExpression{
			Row: &table.ExecuteResult{
				Rows:   []*ptypes.Row{resultRow},
				Schema: resultSchema,
			},
		}, nil
	}

	// Check for typed literals (DATE, TIMESTAMP, INTERVAL)
	if typedLit := value.TypedLiteral(); typedLit != nil {
		// Get the string literal value
		strLit := typedLit.STRING_LITERAL()
		if strLit == nil {
			return nil, fmt.Errorf("[ParsePrimaryLiteral] typed literal missing string value")
		}

		str := strLit.GetText()
		// Remove quotes
		if len(str) >= 2 && str[0] == '\'' && str[len(str)-1] == '\'' {
			str = str[1 : len(str)-1]
		}

		// Check which type of literal
		// TODO not supported in typesystem
		// if typedLit.INTERVAL_TYPE() != nil {
		// 	// Parse interval string
		// 	intervalMicros, err := utils.ParseInterval(str)
		// 	if err != nil {
		// 		return nil, fmt.Errorf("invalid interval: %v", err)
		// 	}
		// 	return &rmodels.ResultRowsExpression{
		// 		Row: &table.ExecuteResult{
		// 			Rows: [][]interface{}{{intervalMicros}},
		// 			Columns: []table.TableColumn{
		// 				{Idx: 0, Name: "?column?", Type: table.ColTypeInterval},
		// 			},
		// 		},
		// 	}, nil
		// } else
		if typedLit.DATE_TYPE() != nil || typedLit.TIMESTAMP_TYPE() != nil {
			// Parse date/timestamp string
			// TODO: Implement proper date/timestamp parsing
			// For now, try parsing ISO 8601 format
			t, err := time.Parse("2006-01-02", str)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", str)
				if err != nil {
					return nil, fmt.Errorf("[ParsePrimaryLiteral] invalid date/timestamp format: %s", str)
				}
			}
			colType := ptypes.PTypeTimestamp
			if typedLit.TIMESTAMP_TYPE() != nil {
				colType = ptypes.PTypeTimestampz
			}
			fields := []serializers.FieldDef{{
				Name: "?column?",
				OID:  colType,
			}}
			resultSchema := serializers.NewBaseSchema(fields)
			val, err := serializers.TimestampSerializerInstance.Serialize(allocator, &t)
			if err != nil {
				return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to serialize timestamp value: %v", err)
			}
			resultRow, err := resultSchema.Pack(allocator, [][]byte{val})
			if err != nil {
				return nil, fmt.Errorf("[ParsePrimaryLiteral] failed to pack timestamp value: %v", err)
			}

			return &rmodels.ResultRowsExpression{
				Row: &table.ExecuteResult{
					Rows:   []*ptypes.Row{resultRow},
					Schema: resultSchema,
				},
			}, nil
		}
	}
	return nil, nil
}
