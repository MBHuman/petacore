package executor

import (
	"fmt"
	"petacore/internal/runtime/rsql/statements"
	"petacore/internal/runtime/rsql/table"
	"petacore/internal/storage"
	"petacore/sdk/pmem"
	"petacore/sdk/serializers"
	ptypes "petacore/sdk/types"
	"strings"
)

// ExecuteShow выполняет SHOW statement
// TODO: расширить поддержку других параметров SHOW
func ExecuteShow(allocator pmem.Allocator, stmt *statements.ShowStatement, storage *storage.DistributedStorageVClock, sessionParams map[string]string, exCtx ExecutorContext) (*table.ExecuteResult, error) {
	param := strings.ToLower(stmt.Parameter)
	switch param {
	case "transaction isolation level":
		fields := []serializers.FieldDef{{
			Name: "transaction_isolation",
			OID:  ptypes.PTypeText,
		}}
		value := "read committed" // TODO: сделать динамическим, в зависимости от реального уровня изоляции
		row, err := serializers.SerializeGeneric(allocator, value, ptypes.PTypeText)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize value: %w", err)
		}
		rows := [][]byte{row}
		schema := serializers.NewBaseSchema(fields)
		outRow, err := schema.Pack(allocator, rows)
		if err != nil {
			return nil, fmt.Errorf("failed to pack row: %w", err)
		}

		return &table.ExecuteResult{
			Rows:   []*ptypes.Row{outRow},
			Schema: schema,
		}, nil
	default:
		if val, ok := sessionParams[param]; ok {
			fields := []serializers.FieldDef{{
				Name: param,
				OID:  ptypes.PTypeText,
			}}
			row, err := serializers.SerializeGeneric(allocator, val, ptypes.PTypeText)
			if err != nil {
				return nil, fmt.Errorf("failed to serialize value: %w", err)
			}
			rows := [][]byte{row}
			schema := serializers.NewBaseSchema(fields)
			outRow, err := schema.Pack(allocator, rows)
			if err != nil {
				return nil, fmt.Errorf("failed to pack row: %w", err)
			}
			return &table.ExecuteResult{
				Rows:   []*ptypes.Row{outRow},
				Schema: schema,
			}, nil
		}
		return nil, fmt.Errorf("unknown parameter: %s", stmt.Parameter)
	}
}
