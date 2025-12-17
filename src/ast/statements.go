package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (n ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	VariableName  string
	IsConstant    bool
	AssignedValue Expr
	ExplicitType  Type
}

func (n VarDeclStmt) stmt() {}

type StructProperty struct {
	IsStatic bool // is property static?
	Type     Type
}

type StructMethod struct {
	IsStatic bool // is property static?
	//Type     FnType
}

type StructDeclStmt struct {
	StructName string
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
}

func (n StructDeclStmt) stmt() {}
