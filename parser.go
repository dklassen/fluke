package main

import (
	"fluke/parser"

	"github.com/antlr4-go/antlr/v4"
)

type SQLParser struct {
	parser   *parser.ClickHouseParser
	listener *ClickHouseSQLListener
}

func NewSQLParser(input string) *SQLParser {
	antlrInput := antlr.NewInputStream(input)
	lexer := parser.NewClickHouseLexer(antlrInput)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewClickHouseParser(stream)

	return &SQLParser{
		parser:   p,
		listener: NewClickHouseSQLListener(),
	}
}

func (sp *SQLParser) Parse() map[string]Schema {
	for sp.parser.GetTokenStream().LA(1) != antlr.TokenEOF {
		stmt := sp.parser.QueryStmt()
		antlr.ParseTreeWalkerDefault.Walk(sp.listener, stmt)
	}
	return sp.listener.Tables
}
