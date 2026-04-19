package main

type ExprVisitor interface {
	VisitBinaryExpr(*BinaryExpr) any
	VisitLogicalExpr(*LogicalExpr) any
	VisitLiteralExpr(*LiteralExpr) any
	VisitGroupingExpr(*GroupingExpr) any
	VisitUnaryExpr(*UnaryExpr) any
	VisitVariableExpr(*VariableExpr) any
	VisitAssignExpr(*AssignExpr) any
	VisitCallExpr(*CallExpr) any
}

type Expr interface {
	Accept(v ExprVisitor) any
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *BinaryExpr) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(b)
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (l *LogicalExpr) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(l)
}

type GroupingExpr struct {
	Expression Expr
}

func (g *GroupingExpr) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(g)
}

type LiteralExpr struct {
	Value any
}

func (l *LiteralExpr) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(l)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (u *UnaryExpr) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(u)
}

type VariableExpr struct {
	Name Token
}

func (ve *VariableExpr) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(ve)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (ae *AssignExpr) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(ae)
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (c *CallExpr) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(c)
}


