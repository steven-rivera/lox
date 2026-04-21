package main

type LoxInstance struct {
	Class *LoxClass
	Fields map[string]any
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		Class: class,
		Fields: map[string]any{},
	}
}

func (i *LoxInstance) toString() string {
	return i.Class.Name + " instance"
}

func (i *LoxInstance) get(name Token) any {
	if field, ok := i.Fields[name.Lexeme]; ok {
		return field
	}

	method := i.Class.findMethod(name.Lexeme);
	if method != nil {
		return method.bind(i)
	}

	return NewRunTimeError(name, "Undefined property '" + name.Lexeme + "'.")
}

func (i *LoxInstance) set(name Token, value any) {
	i.Fields[name.Lexeme] = value
}