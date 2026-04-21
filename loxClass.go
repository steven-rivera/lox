package main

var _ LoxCallable = (*LoxClass)(nil)

type LoxClass struct {
	Name       string
	SuperClass *LoxClass
	Methods    map[string]*LoxFunction
}

func NewLoxClass(name string, superClass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		Name:    name,
		SuperClass: superClass,
		Methods: methods,
	}
}

func (c *LoxClass) Call(interpreter *Interpreter, arguments []any) any {
	instance := NewLoxInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).Call(interpreter, arguments)
	}
	return instance
}

func (c *LoxClass) Arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.Arity()
}

func (c *LoxClass) toString() string {
	return c.Name
}

func (c *LoxClass) findMethod(name string) *LoxFunction {
	if method, ok := c.Methods[name]; ok {
		return method
	}

	if c.SuperClass != nil {
		return c.SuperClass.findMethod(name)
	}
	
	return nil
}
