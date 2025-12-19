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
