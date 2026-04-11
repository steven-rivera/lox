package main

import (
	"bufio"
	"fmt"
	"os"
)

type Lox struct {
	hadError     bool
	justTokenize bool
}

func (l *Lox) run(source string) {
	scanner := NewScanner(source)

	tokens, errs := scanner.scanTokens()
	if errs != nil {
		l.hadError = true
	}

	if l.justTokenize {
		for _, token := range tokens {
			fmt.Println(token.toString())
		}
		return 
	}
	
	parser := NewParser(tokens)
	expr, err := parser.parse()
	if err != nil {
		l.hadError = true
		return
	}

	printer := AstPrinter{}
	fmt.Println(printer.print(expr))

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

func ScanError(line int, message string) error {
	Report(line, "", message)
	return fmt.Errorf("%s", message)
}

func Report(line int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
}

func ParseError(token Token, message string) {
	if token.Type == EOF {
		Report(token.Line, " at end", message)
	} else {
		Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}
