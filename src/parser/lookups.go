package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
)

type bindinPower int

const (
	default_bp bindinPower = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type (
	stmtHandler func(p *parser) ast.Stmt
	nudHandler  func(p *parser) ast.Expr
	ledHandler  func(p *parser, left ast.Expr, bp bindinPower) ast.Expr

	stmtLookup map[lexer.TokenKind]stmtHandler
	ledLookup  map[lexer.TokenKind]ledHandler
	bpLookup   map[lexer.TokenKind]bindinPower
	nudLookup  map[lexer.TokenKind]nudHandler
)

var (
	bpLu   = bpLookup{}
	nudLu  = nudLookup{}
	ledLu  = ledLookup{}
	stmtLu = stmtLookup{}
)

func led(kind lexer.TokenKind, bp bindinPower, ledFn ledHandler) {
	bpLu[kind] = bp
	ledLu[kind] = ledFn
}

func nud(kind lexer.TokenKind, bp bindinPower, nudFn nudHandler) {
	bpLu[kind] = bp
	nudLu[kind] = nudFn
}

func stmt(kind lexer.TokenKind, bp bindinPower, stmtFn stmtHandler) {
	bpLu[kind] = bp
	stmtLu[kind] = stmtFn
}

func createTokenLookups() {

	// Logical
	led(lexer.AND, logical, parseBinaryExpr)
	led(lexer.OR, logical, parseBinaryExpr)
	led(lexer.DOT_DOT, logical, parseBinaryExpr)

	// Relational
	led(lexer.LESS, relational, parseBinaryExpr)
	led(lexer.LESS_EQUALS, relational, parseBinaryExpr)
	led(lexer.GREATER, relational, parseBinaryExpr)
	led(lexer.GREATER_EQUALS, relational, parseBinaryExpr)
	led(lexer.EQUALS, relational, parseBinaryExpr)
	led(lexer.NOT_EQUALS, relational, parseBinaryExpr)

	// Additive & Multiplicative
	led(lexer.PLUS, additive, parseBinaryExpr)
	led(lexer.DASH, additive, parseBinaryExpr)

	led(lexer.STAR, multiplicative, parseBinaryExpr)
	led(lexer.SLASH, multiplicative, parseBinaryExpr)
	led(lexer.PERCENT, multiplicative, parseBinaryExpr)

	// Literals & Symbols
	nud(lexer.NUMBER, primary, parsePrimaryExpr)
	nud(lexer.STRING, primary, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, primary, parsePrimaryExpr)

}
