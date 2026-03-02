package basefuncs

import (
	"context"
	"fmt"
	psdk "petacore/sdk"
	ptypes "petacore/sdk/types"
)

// SubstringFunction implements PostgreSQL substring(text, int, int)
// (the two-index variant). It returns the portion of the input string
// starting from the given 1-based index and spanning the specified length.
// If length argument is omitted or negative, behavior is simplified: when
// length is negative we treat it as zero.
//
// For our catalog query we only need the 3-argument form.

type SubstringFunction struct {
	*psdk.BaseFunction
}

func (f *SubstringFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1009,
		ProName:     "SUBSTRING",
		ProArgTypes: []ptypes.OID{ptypes.PTypeText, ptypes.PTypeInt4, ptypes.PTypeInt4},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *SubstringFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("SUBSTRING requires three arguments")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("SUBSTRING first argument must be text")
	}
	start, ok := args[1].(int)
	if !ok {
		// sometimes numbers arrive as int32/float64 depending on parser
		if v, ok2 := args[1].(int32); ok2 {
			start = int(v)
		} else if v, ok2 := args[1].(float64); ok2 {
			start = int(v)
		} else {
			return nil, fmt.Errorf("SUBSTRING second argument must be integer")
		}
	}
	length, ok := args[2].(int)
	if !ok {
		if v, ok2 := args[2].(int32); ok2 {
			length = int(v)
		} else if v, ok2 := args[2].(float64); ok2 {
			length = int(v)
		} else {
			return nil, fmt.Errorf("SUBSTRING third argument must be integer")
		}
	}

	// PostgreSQL uses 1-based indexing; convert to 0-based
	if start < 1 {
		start = 1
	}
	start--
	if start > len(s) {
		return "", nil
	}
	if length < 0 {
		return "", nil
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end], nil
}
