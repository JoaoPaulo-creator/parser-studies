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

type Parameter struct {
	Name string
	Type Type
}

type FunctionDeclarationStmt struct {
	Parameters []Parameter
	Name       string
	Body       []Stmt
	ReturnType Type
}

func (n FunctionDeclarationStmt) stmt() {}

type IfStmt struct {
	Condition  Expr
	Consequent Stmt
	Alternate  Stmt
}

func (n IfStmt) stmt() {}

type ImportStmt struct {
	Name string
	From string
}

func (n ImportStmt) stmt() {}

type ForeachStmt struct {
	Value    string
	Index    bool
	Iterable Expr
	Body     []Stmt
}

func (n ForeachStmt) stmt() {}

type ClassDeclarationStmt struct {
	Name string
	Body []Stmt
}

func (n ClassDeclarationStmt) stmt() {}
