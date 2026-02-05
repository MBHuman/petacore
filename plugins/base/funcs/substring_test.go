package basefuncs

import (
    "context"
    "testing"
)

func TestSubstringFunction(t *testing.T) {
    f := &SubstringFunction{}
    ctx := context.Background()

    cases := []struct{
        s string
        start int
        length int
        out string
    }{
        {"abcdef",1,3,"abc"},
        {"abcdef",2,2,"bc"},
        {"abcdef",1,10,"abcdef"},
        {"abcdef",10,5,""},
        {"",1,1,""},
    }

    for _, c := range cases {
        res, err := f.Execute(ctx, c.s, c.start, c.length)
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if res != c.out {
            t.Errorf("expected %q got %q", c.out, res)
        }
    }
}
