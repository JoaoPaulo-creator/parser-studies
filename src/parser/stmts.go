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
	isConstant := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Inside variable declaration expected to find variable name").Value
	p.expect(lexer.ASSIGNMENT)
	assinedValue := parseExpr(p, assignment)
	p.expect(lexer.SEMI_COLON)

	return ast.VarDeclStmt{
		IsConstant:    isConstant,
		VariableName:  varName,
		AssignedValue: assinedValue,
	}
}
