package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
	"fmt"
)

type (
	type_nudHandler func(p *parser) ast.Type
	type_ledHandler func(p *parser, left ast.Type, bp bindinPower) ast.Type

	type_nudLookup map[lexer.TokenKind]type_nudHandler
	type_ledLookup map[lexer.TokenKind]type_ledHandler
	type_bpLookup  map[lexer.TokenKind]bindinPower
)

var (
	type_bpLu  = type_bpLookup{}
	type_nudLu = type_nudLookup{}
	type_ledLu = type_ledLookup{}
)

func typeLed(kind lexer.TokenKind, bp bindinPower, ledFn type_ledHandler) {
	type_bpLu[kind] = bp
	type_ledLu[kind] = ledFn
}

func typeNud(kind lexer.TokenKind, nudFn type_nudHandler) {
	type_bpLu[kind] = primary
	type_nudLu[kind] = nudFn
}

func createTokenTypeLookups() {
	typeNud(lexer.IDENTIFIER, parseSymbolType)
	typeNud(lexer.OPEN_BRACKET, parseArrayType)
}

func parseSymbolType(p *parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parseArrayType(p *parser) ast.Type {
	p.advance()
	p.expect(lexer.CLOSE_BRACKET)

	var underlyingType = parseType(p, default_bp)

	return ast.ArrayType{
		Underlying: underlyingType,
	}
}

func parseType(p *parser, bp bindinPower) ast.Type {
	// first parse the nud
	tokenKind := p.currentTokenKind()
	nudFn, exists := type_nudLu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("TYPE_NUD HANDLER EXPECTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
	}

	// while we have a led and the current bp is less than bp of current token
	// continue parsing the left hand side
	left := nudFn(p)
	for type_bpLu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFn, exists := type_ledLu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED HANDLER EXPECTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = ledFn(p, left, type_bpLu[p.currentTokenKind()])
	}

	return left
}
