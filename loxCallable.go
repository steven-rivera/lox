package main

import "time"

type LoxCallable interface {
	Call(*Interpreter, []any) any
	Arity() int
	toString() string
}

type LoxClock struct{}

func (c *LoxClock) Call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().Unix())
}

func (c *LoxClock) Arity() int {
	return 0
}

func (c *LoxClock) toString() string {
	return "<native fn>"
}

type LoxFunction struct {
	Declaration *FunctionStmt 
	Closure *Environment
}

func NewLoxFuncntion(declaration *FunctionStmt, closure *Environment) *LoxFunction{
	return &LoxFunction{
		Declaration: declaration,
		Closure: closure,
	}
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []any) any {
	environment := NewEnvironment(f.Closure)
	for i := range len(f.Declaration.Params) {
		param := f.Declaration.Params[i].Lexeme
		arg := arguments[i]

		environment.define(param, arg)
	}

	if err := interpreter.executeBlock(f.Declaration.Body, environment); err != nil {
		if returnError, ok := err.(ReturnError); ok {
			return returnError.Value
		}
		return err
	}
	return nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) toString() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">" 
}
