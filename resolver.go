package main

type FunctionType = int
type ClassType = int

const (
	FUN_TYPE_NONE FunctionType = iota
	FUN_TYPE_FUNCTION
	FUN_TYPE_INITIALIZER
	FUN_TYPE_METHOD
)

const (
	CLS_TYPE_NONE ClassType = iota
	CLS_TYPE_CLASS
	CLS_TYPE_SUBCLASS
)

var _ ExprVisitor = (*Resolver)(nil)
var _ StmtVisitor = (*Resolver)(nil)

type Resolver struct {
	Interpreter     *Interpreter
	Scopes          Stack[map[string]bool]
	currentFunction FunctionType
	currentClass    ClassType
	hadError        bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		Interpreter:     interpreter,
		Scopes:          Stack[map[string]bool]{},
		currentFunction: FUN_TYPE_NONE,
		currentClass:    CLS_TYPE_NONE,
		hadError:        false,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) any {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitClassStmt(stmt *ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = CLS_TYPE_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.SuperClass != nil {
		if stmt.Name.Lexeme == stmt.SuperClass.Name.Lexeme {
			LoxError(stmt.SuperClass.Name, "A class can't inherit from itself.")
			r.hadError = true
		}
		r.currentClass = CLS_TYPE_SUBCLASS
		r.resolveExpr(stmt.SuperClass)
		r.beginScope()
		r.Scopes.Peek()["super"] = true
	}

	r.beginScope()
	r.Scopes.Peek()["this"] = true
	for _, method := range stmt.Methods {
		declaration := FUN_TYPE_METHOD
		if method.Name.Lexeme == "init" {
			declaration = FUN_TYPE_INITIALIZER
		}
		r.resolveFunction(method, declaration)
	}
	r.endScope()

	if stmt.SuperClass != nil {
		r.endScope()
	}

	r.currentClass = enclosingClass
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *ExprStmt) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *FunctionStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FUN_TYPE_FUNCTION)
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
	if r.currentFunction == FUN_TYPE_NONE {
		LoxError(stmt.Keyword, "Can't return from top-level code.")
		r.hadError = true
	}

	if stmt.Value != nil {
		if r.currentFunction == FUN_TYPE_INITIALIZER {
			LoxError(stmt.Keyword, "Can't return a value from an initializer.")
			r.hadError = true
		}
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

func (r *Resolver) VisitGetExpr(expr *GetExpr) any {
	r.resolveExpr(expr.Object)
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

func (r *Resolver) VisitSetExpr(expr *SetExpr) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)
	return nil
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) any {
	if r.currentClass == CLS_TYPE_NONE {
		LoxError(expr.Keyword, "Can't use 'super' outside of a class.")
		r.hadError = true
	} else if r.currentClass != CLS_TYPE_SUBCLASS {
		LoxError(expr.Keyword, "Can't use 'super' in a class with no superclass.")
		r.hadError = true
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) any {
	if r.currentClass == CLS_TYPE_NONE {
		LoxError(expr.Keyword, "Can't use 'this' outside of a class.")
		r.hadError = true
		return nil
	}

	r.resolveLocal(expr, expr.Keyword)
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
