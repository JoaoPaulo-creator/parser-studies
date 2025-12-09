package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
	"fmt"
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

// helpers
func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	tok := p.currentToken()
	p.pos++
	return tok
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) expectError(expectKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expectKind {
		if err == nil {
			err = fmt.Sprintf("Expected %s but received %s instead\n", lexer.TokenKindString(expectKind), lexer.TokenKindString(kind))
		}
		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
