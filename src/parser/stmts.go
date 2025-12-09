package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, exists := stmtLu[p.currentTokenKind()]
	if exists {
		return stmtFn(p)
	}

	expression := parseExpr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parseVarDeclStmt(p *parser) ast.Stmt {
	var explicitType ast.Type
	var assinedValue ast.Expr

	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value

	if p.currentTokenKind() == lexer.COLON {
		p.advance() // eat the colon
		explicitType = parseType(p, default_bp)
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSIGNMENT)
		assinedValue = parseExpr(p, assignment)
	} else if explicitType == nil {
		panic("Missing either right-hand side in var declaration or explit type.")
	}

	p.expect(lexer.SEMI_COLON)

	if isConstant && assinedValue == nil {
		panic("Cannot define constant without providing a value!")
	}

	return ast.VarDeclStmt{
		IsConstant:    isConstant,
		VariableName:  varName,
		AssignedValue: assinedValue,
		ExplicitType:  explicitType,
	}
}
