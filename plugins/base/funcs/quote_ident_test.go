package basefuncs

import (
	"context"
	"testing"
)

func TestQuoteIdentFunction(t *testing.T) {
	f := &QuoteIdentFunction{}
	ctx := context.Background()

	cases := []struct {
		in  string
		out string
	}{
		{"abc", "abc"},
		{"a\"b", "\"a\"\"b\""},
		{"", "\"\""},
		{"simple name", "\"simple name\""},
	}

	for _, c := range cases {
		res, err := f.Execute(ctx, c.in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != c.out {
			t.Errorf("expected %q got %q", c.out, res)
		}
	}
}
