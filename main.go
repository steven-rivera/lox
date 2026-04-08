package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize [<filename>]")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// TODO: Uncomment the code below to pass the first stage

	lox := Lox{}

	if len(os.Args) == 3 {
		lox.runFile(os.Args[2])
	} else {
		lox.runPrompt()
	}
}
