package parser

import (
	"custom_parser/src/ast"
	"custom_parser/src/lexer"
	"fmt"
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

func parseClassDeclStmt(p *parser) ast.Stmt {
	p.advance()
	className := p.expect(lexer.IDENTIFIER).Value
	classBody := parseBlockStmt(p)

	return ast.ClassDeclarationStmt{
		Name: className,
		Body: ast.ExpectStmt[ast.BlockStmt](classBody).Body,
	}
}

func parseFnDeclStmt(p *parser) ast.Stmt {
	p.advance()
	fnName := p.expect(lexer.IDENTIFIER).Value
	functionParameters, returnType, fnBody := parseFnParamsAndBody(p)

	return ast.FunctionDeclStmt{
		Parameters: functionParameters,
		ReturnType: returnType,
		Body:       fnBody,
		Name:       fnName,
	}
}

func parseStructDeclStmt(p *parser) ast.Stmt {
	p.expect(lexer.STRUCT)
	var properties = map[string]ast.StructProperty{}
	var methods = map[string]ast.StructMethod{}
	var structName = p.expect(lexer.IDENTIFIER).Value

	p.expect(lexer.OPEN_CURLY)
	for p.hasTokens() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		var isStatic bool
		var propertyName string

		if p.currentTokenKind() == lexer.STATIC {
			isStatic = true
			p.expect(lexer.STATIC)
		}

		if p.currentTokenKind() == lexer.IDENTIFIER {
			propertyName = p.expect(lexer.IDENTIFIER).Value
			p.expectError(lexer.COLON, "Expected to find colon following property name inside struct declaration")
			structType := parseType(p, default_bp)
			p.expect(lexer.SEMI_COLON)

			_, exists := properties[propertyName]
			if exists {
				panic(fmt.Sprintf("Property %s has already been defined inside struct declaration", propertyName))
			}

			properties[propertyName] = ast.StructProperty{
				IsStatic: isStatic,
				Type:     structType,
			}

			continue
		}

		panic("cannot currently handle methods inside struct declaration")
	}

	p.expect(lexer.CLOSE_CURLY)
	return ast.StructDeclStmt{
		StructName: structName,
		Properties: properties,
		Methods:    methods,
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
		Name: importFrom,
		From: importFrom,
	}
}

func parseForEarchStmt(p *parser) ast.Stmt {
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
