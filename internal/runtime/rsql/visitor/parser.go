package visitor

import (
	"fmt"
	"log"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"

	"github.com/antlr4-go/antlr/v4"
)

// ParseSQL парсит SQL запрос с помощью ANTLR
func ParseSQL(query string) (statements.SQLStatement, error) {
	input := antlr.NewInputStream(query)
	lexer := parser.NewsqlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewsqlParser(stream)

	errorListener := &sqlErrorListener{}
	p.AddErrorListener(errorListener)
	p.BuildParseTrees = true

	tree := p.Query()
	listener := &sqlListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	if errorListener.err != nil {
		return nil, errorListener.err
	}

	if listener.stmt == nil {
		return nil, fmt.Errorf("failed to parse SQL query")
	}

	log.Printf("DEBUG: Parsed statement type: %s\n", listener.stmt.Type())
	return listener.stmt, nil
}
