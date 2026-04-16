package main

type StmtVisitor interface {
	VisitExpressionStmt(*ExprStmt) any
	VisitPrintStmt(*PrintStmt) any
	VisitVarStmt(*VarStmt) any
	VisitBlockStmt(*BlockStmt) any
	VisitIfStmt(*IfStmt) any
	VisitWhileStmt(*WhileStmt) any
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
