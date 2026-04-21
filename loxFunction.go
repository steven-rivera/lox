package main

type LoxFunction struct {
	Declaration   *FunctionStmt
	Closure       *Environment
	isInitializer bool
}

func NewLoxFunction(declaration *FunctionStmt, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		Declaration:   declaration,
		Closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *LoxFunction) bind(instance *LoxInstance) *LoxFunction {
	environment := NewEnvironment(f.Closure)
	environment.define("this", instance)
	return NewLoxFunction(f.Declaration, environment, f.isInitializer)
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
			if f.isInitializer {
				return f.Closure.getAt(0, "this")
			}
			return returnError.Value
		}
		return err
	}

	if f.isInitializer {
		return f.Closure.getAt(0, "this")
	}

	return nil
}

func (f *LoxFunction) Arity() int {
	return len(f.Declaration.Params)
}

func (f *LoxFunction) toString() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}
