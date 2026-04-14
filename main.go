package main

import (
	"fmt"
	"os"
)

var commands map[string]bool = map[string]bool{
	"tokenize": true,
	"parse":    true,
	"evaluate": true,
	"run": true,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize [<filename>]")
		os.Exit(1)
	}

	command := os.Args[1]

	if _, ok := commands[command]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	lox := Lox{
		command: command,
		interpreter: NewInterpreter(),
	}

	if len(os.Args) == 3 {
		lox.runFile(os.Args[2])
	} else {
		lox.runPrompt()
	}
}
