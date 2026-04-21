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
	VisitGetExpr(*GetExpr) any
	VisitSetExpr(*SetExpr) any
	VisitThisExpr(*ThisExpr) any
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

type GetExpr struct {
	Object Expr
	Name Token
}

func (g *GetExpr) Accept(v ExprVisitor) any {
	return v.VisitGetExpr(g)
}

type SetExpr struct {
	Object Expr
	Name Token
	Value Expr
}

func (s *SetExpr) Accept(v ExprVisitor) any {
	return v.VisitSetExpr(s)
}

type ThisExpr struct {
	Keyword Token
}

func (t *ThisExpr) Accept(v ExprVisitor) any {
	return v.VisitThisExpr(t)
}