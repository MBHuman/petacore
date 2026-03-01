package basefuncs

import (
	"context"
	"fmt"
	"math"
	psdk "petacore/sdk"
	ptypes "petacore/sdk/types"
	"time"
)

type CurrentDatabaseFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentDatabaseFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1000,
		ProName:     "CURRENT_DATABASE",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *CurrentDatabaseFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		// return nil, psdk.NewErrorf(psdk.ErrInvalidArgument, "CURRENT_DATABASE does not take any arguments")
		return nil, fmt.Errorf("CURRENT_DATABASE does not take any arguments")
	}
	return "testdb", nil
}

//

type CurrentSchemaFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentSchemaFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1001,
		ProName:     "CURRENT_SCHEMA",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *CurrentSchemaFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		// return nil, psdk.NewErrorf(psdk.ErrInvalidArgument, "CURRENT_SCHEMA does not take any arguments")
		return nil, fmt.Errorf("CURRENT_SCHEMA does not take any arguments")
	}
	return "public", nil
}

//

type CurrentSchemasFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentSchemasFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1003,
		ProName:     "CURRENT_SCHEMAS",
		ProArgTypes: []ptypes.OID{ptypes.PTypeBool},
		ProRetType:  ptypes.PTypeNameArray,
	}
}

func (f *CurrentSchemasFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("CURRENT_SCHEMAS requires exactly one boolean argument")
	}

	includeImplicit, ok := args[0].(bool)
	if !ok {
		return nil, fmt.Errorf("CURRENT_SCHEMAS argument must be a boolean")
	}

	// Return the current search path
	// When includeImplicit is true, we include implicit schemas like pg_catalog
	// When false, we only return explicit schemas
	if includeImplicit {
		return []string{"pg_catalog", "public"}, nil
	}
	return []string{"public"}, nil
}

//

type CurrentUserFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentUserFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1002,
		ProName:     "CURRENT_USER",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *CurrentUserFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		// return nil, psdk.NewErrorf(psdk.ErrInvalidArgument, "CURRENT_USER does not take any arguments")
		return nil, fmt.Errorf("CURRENT_USER does not take any arguments")
	}

	return "postgres", nil
}

//

type CurrentCatalogFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentCatalogFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1010,
		ProName:     "CURRENT_CATALOG",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeName,
	}
}

func (f *CurrentCatalogFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("CURRENT_CATALOG does not take any arguments")
	}
	return "testdb", nil
}

//

type NowFunction struct {
	*psdk.BaseFunction
}

func (f *NowFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1015,
		ProName:     "NOW",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeTimestampz,
	}
}

func (f *NowFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("NOW does not take any arguments")
	}
	// Check if there's a query start time in the context
	if startTime, ok := ctx.Value("queryStartTime").(int64); ok {
		return startTime, nil
	}
	// Return current time as int64 (microseconds since Unix epoch)
	nowValue := time.Now().UnixMicro()
	return nowValue, nil
}

type CurrentRoleFunction struct {
	*psdk.BaseFunction
}

func (f *CurrentRoleFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1011,
		ProName:     "CURRENT_ROLE",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeName,
	}
}

func (f *CurrentRoleFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("CURRENT_ROLE does not take any arguments")
	}
	return "postgres", nil
}

//

type SessionUserFunction struct {
	*psdk.BaseFunction
}

func (f *SessionUserFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1012,
		ProName:     "SESSION_USER",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeName,
	}
}

func (f *SessionUserFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("SESSION_USER does not take any arguments")
	}
	return "postgres", nil
}

//

type UserFunction struct {
	*psdk.BaseFunction
}

func (f *UserFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1013,
		ProName:     "USER",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeName,
	}
}

func (f *UserFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("USER does not take any arguments")
	}
	return "postgres", nil
}

//

type VersionFunction struct {
	*psdk.BaseFunction
}

func (f *VersionFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1004,
		ProName:     "VERSION",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *VersionFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("VERSION does not take any arguments")
	}

	return "PostgreSQL 14.0 (PetaCore)", nil
}

//

type RoundFunction struct {
	*psdk.BaseFunction
}

func (r *RoundFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("ROUND function requires exactly one argument")
	}

	num, ok := args[0].(float64)
	if !ok {
		return nil, fmt.Errorf("ROUND function argument must be a numeric type")
	}

	rounded := math.Round(num)

	return int16(rounded), nil
}

func (r *RoundFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1005,
		ProName:     "ROUND",
		ProArgTypes: []ptypes.OID{ptypes.PTypeNumeric},
		ProRetType:  ptypes.PTypeInt2,
	}
}

//

type PgTableIsVisibleFunction struct {
	*psdk.BaseFunction
}

func (p *PgTableIsVisibleFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("pg_table_is_visible function requires exactly one argument")
	}

	// Accept common numeric representations. Some callers (expression evaluator)
	// may produce float64 for numeric columns, so accept those too.
	switch args[0].(type) {
	case int, int32, int64, int16, uint32, uint64, float32, float64:
		// For now, treat all tables as visible. In future, consult schema/catalog.
		return true, nil
	default:
		return nil, fmt.Errorf("pg_table_is_visible function argument must be a numeric type, got %T", args[0])
	}
}

func (p *PgTableIsVisibleFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1006,
		ProName:     "PG_TABLE_IS_VISIBLE",
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt4},
		ProRetType:  ptypes.PTypeBool,
	}
}

//

type PgGetUserbyIDFunction struct {
	*psdk.BaseFunction
}

func (p *PgGetUserbyIDFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("PG_GET_USERBYID function requires exactly one argument")
	}

	_, ok := args[0].(int16)
	if !ok {
		return nil, fmt.Errorf("PG_GET_USERBYID function argument must be an integer")
	}

	return "postgres", nil
}

func (p *PgGetUserbyIDFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1007,
		ProName:     "PG_GET_USERBYID",
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt2},
		ProRetType:  ptypes.PTypeText,
	}
}

//

type PgPostmasterStartTimeFunction struct {
	*psdk.BaseFunction
}

func (p *PgPostmasterStartTimeFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("PG_POSTMASTER_START_TIME function does not take any arguments")
	}

	// Return a fixed start time as int64 (microseconds since Unix epoch)
	// For example: 2024-01-01 00:00:00 UTC
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return startTime.UnixMicro(), nil
}

func (p *PgPostmasterStartTimeFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1008,
		ProName:     "PG_POSTMASTER_START_TIME",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeTimestampz,
	}
}

//

type PgBackendPidFunction struct {
	*psdk.BaseFunction
}

func (p *PgBackendPidFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("PG_BACKEND_PID function does not take any arguments")
	}

	// Return a mock PID. In a real implementation, this would return the actual backend process ID.
	// For now, we return a consistent value to satisfy PostgreSQL clients.
	return int32(12345), nil
}

func (p *PgBackendPidFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         2026,
		ProName:     "PG_BACKEND_PID",
		ProArgTypes: []ptypes.OID{},
		ProRetType:  ptypes.PTypeInt4,
	}
}

//

type MaxFunction struct {
	*psdk.BaseFunction
}

func (m *MaxFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("MAX function requires at least one argument")
	}

	maxVal, ok := args[0].(int16)
	if !ok {
		return nil, fmt.Errorf("MAX function argument must be an integer")
	}

	for _, arg := range args[1:] {
		v, ok := arg.(int16)
		if !ok {
			return nil, fmt.Errorf("MAX function argument must be an integer")
		}
		if v > maxVal {
			maxVal = v
		}
	}

	return maxVal, nil
}

func (m *MaxFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1009,
		ProName:     "MAX",
		ProArgTypes: []ptypes.OID{ptypes.PTypeInt2},
		ProRetType:  ptypes.PTypeInt2,
	}
}
