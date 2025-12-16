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

func parseExpressionStmt(p *parser) ast.ExpressionStmt {
	expression := parseExpr(p, default_bp)
	p.expect(lexer.SEMI_COLON)
	return ast.ExpressionStmt{
		Expression: expression,
	}
}

func parseBlockStmt(p *parser) ast.Stmt {
	p.expect(lexer.OPEN_CURLY)
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseStmt(p))
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.BlockStmt{
		Body: body,
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

func parseFnParamsAndBody(p *parser) ([]ast.Parameter, ast.Type, []ast.Stmt) {
	functionParams := make([]ast.Parameter, 0)
	p.expect(lexer.OPEN_PAREN)

	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		paramName := p.expect(lexer.IDENTIFIER).Value
		p.expect(lexer.COLON)
		paramType := parseType(p, default_bp)

		functionParams = append(functionParams, ast.Parameter{
			Name: paramName,
			Type: paramType,
		})

		if !p.currentToken().IsOneOfMany(lexer.CLOSE_PAREN, lexer.EOF) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.CLOSE_PAREN)
	var returnType ast.Type

	if p.currentTokenKind() == lexer.COLON {
		p.advance()
		returnType = parseType(p, default_bp)
	}

	functionBody := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p)).Body
	return functionParams, returnType, functionBody
}

func parseFnDeclaration(p *parser) ast.Stmt {
	p.advance()
	functionName := p.expect(lexer.IDENTIFIER).Value
	functionParams, returnType, functionBody := parseFnParamsAndBody(p)

	return ast.FunctionDeclarationStmt{
		Parameters: functionParams,
		ReturnType: returnType,
		Body:       functionBody,
		Name:       functionName,
	}
}

func parseIfStmt(p *parser) ast.Stmt {
	p.advance()
	condition := parseExpr(p, assignment)
	consequent := parseBlockStmt(p)

	var alternate ast.Stmt
	if p.currentTokenKind() == lexer.ELSE {
		p.advance()

		if p.currentTokenKind() == lexer.IF {
			alternate = parseIfStmt(p)
		} else {
			alternate = parseBlockStmt(p)
		}
	}

	return ast.IfStmt{
		Condition:  condition,
		Consequent: consequent,
		Alternate:  alternate,
	}
}

func parseImportStmt(p *parser) ast.Stmt {
	p.advance()
	var importFrom string
	importName := p.expect(lexer.IDENTIFIER).Value

	if p.currentTokenKind() == lexer.FROM {
		p.advance()
		importFrom = p.expect(lexer.STRING).Value
	} else {
		importFrom = importName
	}

	p.expect(lexer.SEMI_COLON)
	return ast.ImportStmt{
		Name: importName,
		From: importFrom,
	}
}

func parseForeachStmt(p *parser) ast.Stmt {
	p.advance()
	valueName := p.expect(lexer.IDENTIFIER).Value

	var index bool
	if p.currentTokenKind() == lexer.COMMA {
		p.expect(lexer.COMMA)
		p.expect(lexer.IDENTIFIER)
		index = true
	}

	p.expect(lexer.IN)
	iterable := parseExpr(p, default_bp)
	body := ast.ExpectStmt[ast.BlockStmt](parseBlockStmt(p)).Body

	return ast.ForeachStmt{
		Value:    valueName,
		Index:    index,
		Iterable: iterable,
		Body:     body,
	}
}

func parseClassDeclarationStmt(p *parser) ast.Stmt {
	p.advance()
	className := p.expect(lexer.IDENTIFIER).Value
	classBody := parseBlockStmt(p)

	return ast.ClassDeclarationStmt{
		Name: className,
		Body: ast.ExpectStmt[ast.BlockStmt](classBody).Body,
	}
}
