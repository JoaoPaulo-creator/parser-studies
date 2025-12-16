package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	createTokenLookups()
	createTokenTypeLookups()

	return &parser{
		tokens: tokens,
	}
}

func Parse(tokens []lexer.Token) ast.BlockStmt {
	body := make([]ast.Stmt, 0)
	p := createParser(tokens)

	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.BlockStmt{
		Body: body,
	}
}
