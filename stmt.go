package main

type StmtVisitor interface {
	VisitExpressionStmt(*ExprStmt) any
	VisitPrintStmt(*PrintStmt) any
	VisitVarStmt(*VarStmt) any
	VisitBlockStmt(*BlockStmt) any 
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
	Name Token
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
