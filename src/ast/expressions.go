package ast

import "custom_parser/src/lexer"

// ---------
// Literals
// ---------
type NumberExpr struct {
	Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

// ---------
// Complex expressions
// ---------

// 10 + 5 * 2
type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (n BinaryExpr) expr() {}

// examples:
// a = a + 5;
// a += 5;
// foo.bar += 5;
type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
}

func (n AssignmentExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

type MemberExpr struct {
	Member   Expr
	Property string
}

func (n MemberExpr) expr() {}

type CallExpr struct {
	Method    Expr
	Arguments []Expr
}

func (n CallExpr) expr() {}

type ComputedExpr struct {
	Member   Expr
	Property Expr
}

func (n ComputedExpr) expr() {}

type RangeExpr struct {
	Lower Expr
	Upper Expr
}

func (n RangeExpr) expr() {}

type FunctionExpr struct {
	Parameters []Parameter
	Body       []Stmt
	ReturnType Type
}

func (n FunctionExpr) expr() {}

type NewExpr struct {
	Instantiation CallExpr
}

func (n NewExpr) expr() {}

type ArrayLiteral struct {
	Contents []Expr
}

func (n ArrayLiteral) expr() {}

type StructInstantiationExpr struct {
	StructName string
	Properties map[string]Expr
}

func (n StructInstantiationExpr) expr() {}

type ArrayInstantiationExpr struct {
	Underlying Type
	Contents   []Expr
}

func (n ArrayInstantiationExpr) expr() {}
