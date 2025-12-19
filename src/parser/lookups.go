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

func nud(kind lexer.TokenKind, nudFn nudHandler) {
	nudLu[kind] = nudFn
}

func stmt(kind lexer.TokenKind, stmtFn stmtHandler) {
	bpLu[kind] = default_bp
	stmtLu[kind] = stmtFn
}

func createTokenLookups() {
	led(lexer.ASSIGNMENT, assignment, parseAssignmentExpr)
	led(lexer.PLUS_EQUALS, assignment, parseAssignmentExpr)
	led(lexer.MINUS_EQUALS, assignment, parseAssignmentExpr)

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
	nud(lexer.NUMBER, parsePrimaryExpr)
	nud(lexer.STRING, parsePrimaryExpr)
	nud(lexer.IDENTIFIER, parsePrimaryExpr)
	nud(lexer.OPEN_PAREN, parseGroupingExpr)
	nud(lexer.DASH, parsePrefixExpr)

	// Call/Member/Arrays expressions
	led(lexer.OPEN_CURLY, call, parseStructInstantiationExpr)
	nud(lexer.OPEN_BRACKET, parseArrayInstantiationExpr)

	// Grouping Expr
	// nud(lexer.FN, default_bp, parseFnExpr)

	// Statements
	stmt(lexer.CONST, parseVarDeclStmt)
	stmt(lexer.LET, parseVarDeclStmt)
	stmt(lexer.STRUCT, parseStructDeclStmt)
	stmt(lexer.FN, parseFnDeclStmt)
	stmt(lexer.CLASS, parseClassDeclStmt)
}
