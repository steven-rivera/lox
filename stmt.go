package main

type StmtVisitor interface {
	VisitExpressionStmt(*ExprStmt) any
	VisitPrintStmt(*PrintStmt) any
	VisitVarStmt(*VarStmt) any
	VisitBlockStmt(*BlockStmt) any
	VisitIfStmt(*IfStmt) any
	VisitWhileStmt(*WhileStmt) any
	VisitFunctionStmt(*FunctionStmt) any
	VisitReturnStmt(*ReturnStmt) any
	VisitClassStmt(*ClassStmt) any
}

type Stmt interface {
	Accept(v StmtVisitor) any
}

type ExprStmt struct {
	Expression Expr
}

func (es *ExprStmt) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(es)
}

type PrintStmt struct {
	Expression Expr
}

func (p *PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(p)
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func (vs *VarStmt) Accept(v StmtVisitor) any {
	return v.VisitVarStmt(vs)
}

type BlockStmt struct {
	Statements []Stmt
}

func (bs *BlockStmt) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(bs)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (is *IfStmt) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(is)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (w *WhileStmt) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(w)
}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (f *FunctionStmt) Accept(v StmtVisitor) any {
	return v.VisitFunctionStmt(f)
}

type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

func (r *ReturnStmt) Accept(v StmtVisitor) any {
	return v.VisitReturnStmt(r)
}

type ClassStmt struct {
	Name       Token
	SuperClass *VariableExpr
	Methods    []*FunctionStmt
}

func (c *ClassStmt) Accept(v StmtVisitor) any {
	return v.VisitClassStmt(c)
}
