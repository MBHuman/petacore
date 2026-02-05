package baseplugin

import (
	"context"
	"fmt"
	psdk "petacore/sdk"
	"strings"
)

type UpperFunction struct{}

func (f *UpperFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:     0, // OID will be assigned by the registry
		ProName: "UPPER",
		ProArgTypes: []psdk.OID{
			psdk.PTypeText,
		},
		ProRetType: psdk.PTypeText,
	}
}

func (f *UpperFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper function requires exactly 1 argument")
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("upper function argument must be string")
	}
	return strings.ToUpper(str), nil
}
