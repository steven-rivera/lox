package main

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	if f, ok := expr.Value.(float64); ok {
		return floatString(f)
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
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
