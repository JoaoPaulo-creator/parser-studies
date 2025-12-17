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

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

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
