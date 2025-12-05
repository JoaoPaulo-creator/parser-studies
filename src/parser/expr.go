package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp bindinPower) ast.Expr {
	// first parse the nud
	tokenKind := p.currentTokenKind()
	nudFn, exists := nudLu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD HANDLER EXPTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
	}

	// while we have a led and the current bp is less than bp of current token
	// continue parsing the left hand side
	left := nudFn(p)
	for bpLu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFn, exists := ledLu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED HANDLER EXPTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
		}

		left = ledFn(p, left, bpLu[p.currentTokenKind()])
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		number, _ := strconv.ParseFloat(p.advance().Value, 64)
		return ast.NumberExpr{
			Value: number,
		}
	case lexer.STRING:
		return ast.StringExpr{
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.advance().Value,
		}
	default:
		panic(fmt.Sprintf("Cannot create primary expression from %s\n", lexer.TokenKindString(p.currentTokenKind())))
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp bindinPower) ast.Expr {
	operatorToken := p.advance()
	right := parseExpr(p, bp)

	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}

}
