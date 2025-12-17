package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/helpers"
	"custom_parser/src/lexer"
	"fmt"
	"strconv"
)

func parseExpr(p *parser, bp bindinPower) ast.Expr {
	// first parse the nud
	tokenKind := p.currentTokenKind()
	nudFn, exists := nudLu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("NUD HANDLER EXPECTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
	}

	// while we have a led and the current bp is less than bp of current token
	// continue parsing the left hand side
	left := nudFn(p)
	for bpLu[p.currentTokenKind()] > bp {
		tokenKind = p.currentTokenKind()
		ledFn, exists := ledLu[tokenKind]

		if !exists {
			panic(fmt.Sprintf("LED HANDLER EXPECTED FOR TOKEN %s\n", lexer.TokenKindString(tokenKind)))
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

func parsePrefixExpr(p *parser) ast.Expr {
	operatorToken := p.advance()
	rhs := parseExpr(p, default_bp)

	return ast.PrefixExpr{
		Operator:  operatorToken,
		RightExpr: rhs,
	}
}

func parseGroupingExpr(p *parser) ast.Expr {
	p.advance() //advance past grouping start
	expr := parseExpr(p, default_bp)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parseAssignmentExpr(p *parser, left ast.Expr, bp bindinPower) ast.Expr {
	operatorToken := p.advance()
	rhs := parseExpr(p, bp)

	return ast.AssignmentExpr{
		Operator: operatorToken,
		Value:    rhs,
		Assignee: left,
	}
}

func parseStructInstantiationExpr(p *parser, left ast.Expr, bp bindinPower) ast.Expr {
	var structName = helpers.ExpectType[ast.SymbolExpr](left).Value
	var properties = map[string]ast.Expr{}

	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		propertyName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		expr := parseExpr(p, logical)

		properties[propertyName] = expr
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.StructInstantiationExpr{
		StructName: structName,
		Properties: properties,
	}
}

func parseArrayInstantiationExpr(p *parser) ast.Expr {
	var underlyingType ast.Type
	var contents = []ast.Expr{}
	p.expect(lexer.OPEN_BRACKET)
	p.expect(lexer.CLOSE_BRACKET)

	underlyingType = parseType(p, default_bp)
	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		contents = append(contents, parseExpr(p, logical))
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.ArrayInstantiationExpr{
		Underlying: underlyingType,
		Contents:   contents,
	}
}
