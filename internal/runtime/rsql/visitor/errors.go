package visitor

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

// sqlErrorListener implements ANTLR error listener
type sqlErrorListener struct {
	antlr.DefaultErrorListener
	err error
}

func (el *sqlErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	el.err = fmt.Errorf("line %d:%d %s", line, column, msg)
}
