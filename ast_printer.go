package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *BinaryExpr) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitLogicalExpr(expr *LogicalExpr) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *GroupingExpr) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *LiteralExpr) any {
	if expr.Value == nil {
		return "nil"
	}
	if f, ok := expr.Value.(float64); ok {
		return floatString(f)
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitUnaryExpr(expr *UnaryExpr) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *VariableExpr) any {
	return p.parenthesize(expr.Name.Lexeme)
}

func (p *AstPrinter) VisitAssignExpr(expr *AssignExpr) any {
	return p.parenthesize("var " + expr.Name.Lexeme + " = " + p.print(expr.Value))
}

func (p *AstPrinter) VisitCallExpr(expr *CallExpr) any {
	return p.parenthesize("func " + p.print(expr.Callee), expr.Arguments...)
}

func (p *AstPrinter) VisitGetExpr(expr *GetExpr) any {
	return p.parenthesize("get ", expr.Object)
}

func (p *AstPrinter) VisitSetExpr(expr *SetExpr) any {
	return p.parenthesize("set ", expr.Object, expr.Value)
}

func (p *AstPrinter) VisitThisExpr(expr *ThisExpr) any {
	return "this"
}

func (p *AstPrinter) VisitSuperExpr(expr *SuperExpr) any {
	return "super " + expr.Method.Lexeme
}

func (p *AstPrinter) print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}
	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.Accept(p).(string))
	}
	builder.WriteString(")")

	return builder.String()
}

// func Eval(expr Expr) any {
//     switch e := expr.(type) {
//     case *Binary:
//         return Eval(e.Left) + Eval(e.Right)
//     ...
// }

// func Print(expr Expr) string {
//     switch e := expr.(type) {
//     case *Binary:
//         return "(" + Print(e.Left) + " + " + Print(e.Right) + ")"
//     ...
// }
