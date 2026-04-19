package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type Interpreter struct {
	Globals     *Environment
	Environment *Environment
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.define("clock", &LoxClock{})

	return &Interpreter{
		Globals:     globals,
		Environment: globals,
	}
}

func (i *Interpreter) VisitBinaryExpr(expr *BinaryExpr) any {
	left := i.evaluate(expr.Left)
	if err, ok := left.(error); ok {
		return err
	}
	right := i.evaluate(expr.Right)
	if err, ok := right.(error); ok {
		return err
	}

	switch expr.Operator.Type {
	case GREATER:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) >= right.(float64)
	case LESS:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) <= right.(float64)
	case SLASH:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) / right.(float64)
	case STAR:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) * right.(float64)
	case MINUS:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return err
		}
		return left.(float64) - right.(float64)
	case PLUS:
		switch l := left.(type) {
		case float64:
			if r, ok := right.(float64); ok {
				return l + r
			}
		case string:
			if r, ok := right.(string); ok {
				return l + r
			}
		}

		return NewRunTimeError(expr.Operator, "Operands must be two numbers or two strings.")

	case BANG_EQUAL:
		return !i.isEqual(left, right)
	case EQUAL_EQUAL:
		return i.isEqual(left, right)

	}

	// Unreachable.
	return nil
}

func (i *Interpreter) VisitLogicalExpr(expr *LogicalExpr) any {
	left := i.evaluate(expr.Left)
	if err, ok := left.(error); ok {
		return err
	}

	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	right := i.evaluate(expr.Right)
	if err, ok := left.(error); ok {
		return err
	}
	return right
}

func (i *Interpreter) VisitGroupingExpr(expr *GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *LiteralExpr) any {
	return expr.Value
}

func (i *Interpreter) VisitUnaryExpr(expr *UnaryExpr) any {
	right := i.evaluate(expr.Right)
	if err, ok := right.(error); ok {
		return err
	}
	switch expr.Operator.Type {
	case MINUS:
		if err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return err
		}
		return -right.(float64)
	case BANG:
		return !i.isTruthy(right)
	}

	// Unreachable.
	return nil
}

func (i *Interpreter) VisitVariableExpr(expr *VariableExpr) any {
	return i.Environment.get(expr.Name)
}

func (i *Interpreter) VisitAssignExpr(expr *AssignExpr) any {
	value := i.evaluate(expr.Value)
	if err := i.Environment.assign(expr.Name, value); err != nil {
		return err
	}
	return value
}

func (i *Interpreter) VisitCallExpr(expr *CallExpr) any {
	callee := i.evaluate(expr.Callee)

	var arguments []any
	for _, argument := range expr.Arguments {
		value := i.evaluate(argument)
		if err, ok := value.(error); ok {
			return err
		}
		arguments = append(arguments, value)
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		return NewRunTimeError(expr.Paren, "Can only call functions and classes.")
	}
	if len(arguments) != function.Arity() {
		return NewRunTimeError(expr.Paren,
			fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments)),
		)
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ExprStmt) any {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitIfStmt(stmt *IfStmt) any {
	value := i.evaluate(stmt.Condition)
	if err, ok := value.(error); ok {
		return err
	}

	if i.isTruthy(value) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *PrintStmt) any {
	value := i.evaluate(stmt.Expression)
	if err, ok := value.(error); ok {
		return err
	}
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ReturnStmt) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
		if err, ok := value.(error); ok {
			return err
		}
	}

	return NewReturnError(value)
}

func (i *Interpreter) VisitVarStmt(stmt *VarStmt) any {
	var value any = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
		if err, ok := value.(error); ok {
			return err
		}
	}
	i.Environment.define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *WhileStmt) any {
	for {
		value := i.evaluate(stmt.Condition)
		if err, ok := value.(error); ok {
			return err
		}

		if !i.isTruthy(value) {
			break
		}

		if err := i.execute(stmt.Body); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *FunctionStmt) any {
	function := NewLoxFuncntion(stmt, i.Environment)
	i.Environment.define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *BlockStmt) any {
	return i.executeBlock(stmt.Statements, NewEnvironment(i.Environment))
}

func (i *Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) any {
	previous := i.Environment
	i.Environment = environment

	defer func() {
		i.Environment = previous
	}()

	for _, stmt := range statements {
		if err, ok := i.execute(stmt).(error); ok {
			return err
		}
	}

	return nil
}

func (i *Interpreter) isTruthy(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}

func (i *Interpreter) isEqual(a, b any) bool {
	return reflect.DeepEqual(a, b)
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return NewRunTimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any) error {
	_, lok := left.(float64)
	_, rok := right.(float64)
	if lok && rok {
		return nil
	}
	return NewRunTimeError(operator, "Operands must be numbers.")
}

func (i *Interpreter) interpret(statements []Stmt) error {
	for _, statement := range statements {
		err := i.execute(statement)
		if err, ok := err.(error); ok {
			return err
		}
	}
	return nil
}

func (i *Interpreter) interpretExpr(expr Expr) error {
	value := i.evaluate(expr)
	if err, ok := value.(error); ok {
		return err
	}
	fmt.Print(i.stringify(value))
	return nil
}

func (i *Interpreter) stringify(object any) string {
	switch object := object.(type) {
	case nil:
		return "nil"
	case float64:
		return strconv.FormatFloat(object, 'f', -1, 64)
	case LoxCallable:
		return object.toString()
	default:
		return fmt.Sprint(object)
	}
}
