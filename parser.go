package main

import (
	"fmt"
	"slices"
)

type Parser struct {
	Tokens  []Token
	Current int
	Lox     Lox
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) peek() Token {
	return p.Tokens[p.Current]
}

func (p *Parser) previous() Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) parse() ([]Stmt, []error) {
	var errs []error = nil
	statements := []Stmt{}

	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			errs = append(errs, err)
		} else {
			statements = append(statements, statement)
		}
	}
	return statements, errs
}

func (p *Parser) parseExpr() (Expr, error) {
	return p.expression()
}

func (p *Parser) declaration() (Stmt, error) {
	var stmt Stmt
	var err error

	if p.match(FUN) {
		stmt, err = p.function("function")
	} else if p.match(VAR) {
		stmt, err = p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}

	if err != nil {
		p.synchronize()
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) function(kind string) (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(LEFT_PAREN, "Expect '(' after "+kind+" name."); err != nil {
		return nil, err
	}

	var parameters []Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.Error(p.peek(), "Can't have more than 255 parameters.")
			}

			param, err := p.consume(IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}

			parameters = append(parameters, param)

			if !p.match(COMMA) {
				break
			}
		}
	}

	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after parameters."); err != nil {
		return nil, err
	}
	if _, err := p.consume(LEFT_BRACE, "Expect '{' before "+kind+" body."); err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &FunctionStmt{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer Expr = nil
	if p.match(EQUAL) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		initializer = expr
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after variable declaration."); err != nil {
		return nil, err
	}

	return &VarStmt{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(RETURN) {
		return p.returnStatement()
	}

	if p.match(WHILE) {
		return p.whileStatement()
	}

	if p.match(FOR) {
		return p.forStatement()
	}

	if p.match(LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return &BlockStmt{
			Statements: block,
		}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Stmt, error) {
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after if condition."); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after value."); err != nil {
		return nil, err
	}

	return &PrintStmt{
		Expression: expr,
	}, nil

}

func (p *Parser) returnStatement() (Stmt, error) {
	keyword := p.previous()
	
	var value Expr = nil
	var err error
	if !p.check(SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after return value."); err != nil{
		return nil, err
	}

	return &ReturnStmt{
		Keyword: keyword,
		Value: value,
	}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after condition."); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) forStatement() (Stmt, error) {
	if _, err := p.consume(LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
		return nil, err
	}

	var err error = nil
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}

	if err != nil {
		return nil, err
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after loop condition."); err != nil {
		return nil, err
	}

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.consume(RIGHT_PAREN, "Expect ')' after for clauses."); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &BlockStmt{
			Statements: []Stmt{body, &ExprStmt{Expression: increment}},
		}
	}

	if condition == nil {
		condition = &LiteralExpr{Value: true}
	}

	body = &WhileStmt{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &BlockStmt{
			Statements: []Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(SEMICOLON, "Expect ';' after expression."); err != nil {
		return nil, err
	}

	return &ExprStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) block() ([]Stmt, error) {
	var statements []Stmt

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	if _, err := p.consume(RIGHT_BRACE, "Expect '}', after block."); err != nil {
		return nil, err
	}

	return statements, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(*VariableExpr); ok {
			return &AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		}

		p.Error(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &LogicalExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()

		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	var expr, err = p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var arguments []Expr
	if !p.check(RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.Error(p.peek(), "Can't have more than 255 arguments.")
			}

			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expr)

			if !p.match(COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &LiteralExpr{
			Value: false,
		}, nil
	}
	if p.match(TRUE) {
		return &LiteralExpr{
			Value: true,
		}, nil
	}
	if p.match(NIL) {
		return &LiteralExpr{
			Value: nil,
		}, nil
	}

	if p.match(NUMBER, STRING) {
		return &LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(IDENTIFIER) {
		return &VariableExpr{
			Name: p.previous(),
		}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}

		return &GroupingExpr{
			Expression: expr,
		}, nil
	}

	return nil, p.Error(p.peek(), "Expect expression.")
}

func (p *Parser) consume(t TokenType, message string) (Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	return Token{}, p.Error(p.peek(), message)
}

func (p *Parser) Error(token Token, message string) error {
	ParseError(token, message)
	return fmt.Errorf("%s", message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}
		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
