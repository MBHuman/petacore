package basefuncs

import (
	"context"
	"fmt"
	psdk "petacore/sdk"
	ptypes "petacore/sdk/types"
	"strings"
)

// QuoteIdentFunction implements PostgreSQL's quote_ident(text) -> text
// It wraps its argument in double quotes and escapes any embedded double
// quotes by doubling them.  This is useful for generating safe identifiers
// when the input might contain uppercase letters or special characters.
//
// Note: in the simplified environment of PetaCore we don't have a session
// catalog to check for actual names, so the function just performs the
// string transformation.
//
// SQL example: SELECT quote_ident(c.relname) FROM pg_catalog.pg_class c;

type QuoteIdentFunction struct {
	*psdk.BaseFunction
}

func (f *QuoteIdentFunction) GetFunction() *psdk.Function {
	return &psdk.Function{
		OID:         1014,
		ProName:     "QUOTE_IDENT",
		ProArgTypes: []ptypes.OID{ptypes.PTypeText},
		ProRetType:  ptypes.PTypeText,
	}
}

func (f *QuoteIdentFunction) Execute(ctx context.Context, args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("QUOTE_IDENT requires exactly one argument")
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("QUOTE_IDENT argument must be text")
	}

	// Return unquoted identifier when it's a valid unquoted SQL identifier
	// (lowercase, starts with letter or underscore, contains only
	// letters/digits/underscores). Otherwise escape and quote it.
	isUnquoted := func(t string) bool {
		if t == "" {
			return false
		}
		for i, r := range t {
			if i == 0 {
				if r != '_' && (r < 'a' || r > 'z') {
					return false
				}
			} else {
				if r != '_' && (r < 'a' || r > 'z') && (r < '0' || r > '9') {
					return false
				}
			}
		}
		return true
	}

	if isUnquoted(s) {
		return s, nil
	}

	// Escape any existing double quotes by doubling them
	escaped := strings.ReplaceAll(s, `"`, `""`)
	return `"` + escaped + `"`, nil
}
