package basefuncs

import (
	"context"
	"fmt"
	psdk "petacore/sdk"
	"unicode/utf8"
)

// LengthFunction implements PostgreSQL LENGTH(text) -> int
// Returns number of characters in the input string.
type LengthFunction struct {
	*psdk.BaseFunction
}

func (f *LengthFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1317,
		ProName:     "LENGTH",
		ProArgTypes: []psdk.OID{psdk.PTypeText},
		ProRetType:  psdk.PTypeInt4,
	}
}

func (f *LengthFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("LENGTH requires exactly one argument")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("LENGTH argument must be text")
	}
	// Count Unicode code points
	return utf8.RuneCountInString(s), nil
}
