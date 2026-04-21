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