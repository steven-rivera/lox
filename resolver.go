package main

type FunctionType = int

const (
	NONE FunctionType = iota
	FUNCTION
)

var _ ExprVisitor = (*Resolver)(nil)
var _ StmtVisitor = (*Resolver)(nil)
type Resolver struct {
	Interpreter     *Interpreter
	Scopes          Stack[map[string]bool]
	currentFunction FunctionType
	hadError        bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		Interpreter: interpreter,
		Scopes:      Stack[map[string]bool]{},
		currentFunction: NONE,
		hadError:    false,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) any {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExprStmt) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUNCTION)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *IfStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *PrintStmt) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *ReturnStmt) any {
	if r.currentFunction == NONE {
		LoxError(stmt.Keyword, "Can't return from top-level code.")
		r.hadError = true
	}
	
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *VarStmt) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *WhileStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) any {
	r.resolveExpr(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolveExpr(argument)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) any {
	if !r.Scopes.IsEmpty() {
		if isReady, ok := r.Scopes.Peek()[expr.Name.Lexeme]; ok && !isReady {
			LoxError(expr.Name, "Can't read local variable in its own initializer.")
			r.hadError = true
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveStmts(statements []Stmt) any {
	for _, statement := range statements {
		r.resolveStmt(statement)
	}
	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveFunction(function *FunctionStmt, funcType FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveExpr(expr Expr) {
	expr.Accept(r)
}

func (r *Resolver) beginScope() {
	r.Scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.Scopes.Pop()
}

func (r *Resolver) declare(name Token) {
	if r.Scopes.IsEmpty() {
		return
	}

	scope := r.Scopes.Peek()
	if _, ok := scope[name.Lexeme]; ok {
		LoxError(name, "Already a variable with this name in this scope.")
		r.hadError = true
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.Scopes.IsEmpty() {
		return
	}

	scope := r.Scopes.Peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.Scopes.Size() - 1; i >= 0; i-- {
		if _, ok := r.Scopes.Get(i)[name.Lexeme]; ok {
			r.Interpreter.resolve(expr, r.Scopes.Size()-1-i)
			return
		}
	}
}
