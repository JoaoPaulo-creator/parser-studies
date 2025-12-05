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
