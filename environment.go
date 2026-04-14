package main

type Environment struct {
	Values map[string]any
	Enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Values: make(map[string]any),
		Enclosing: enclosing,
	}
}

func (e *Environment) define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) get(name Token) any {
	if value, ok := e.Values[name.Lexeme]; ok {
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.get(name)
	}

	return NewRunTimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}

func (e *Environment) assign(name Token, value any) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.assign(name, value)
	}

	return NewRunTimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
