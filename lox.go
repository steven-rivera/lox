package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError bool
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)

	tokens, errs := scanner.scanTokens()
	if errs != nil {
		l.hadError = true
	}

	for _, token := range tokens {
		fmt.Println(token.toString())
	}

}
func (l *Lox) runFile(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l.run(string(fileContents))

	if l.hadError {
		os.Exit(65)
	}
}

func (l *Lox) runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		l.run(scanner.Text())
		l.hadError = false
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func Error(line int, message string) error {
	Report(line, "", message)
	return fmt.Errorf("%s", message)
}

func Report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
}
