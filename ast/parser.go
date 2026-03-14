package ast

import (
	"slices"

	"github.com/Moorad/carapace/lexer"
	"github.com/Moorad/carapace/utils"
)

type Parser[T any] struct {
	Tokens  []lexer.Token
	current int
}

func NewParser(tokens []lexer.Token) Parser[any] {
	return Parser[any]{
		Tokens: tokens,
	}
}

func (p *Parser[T]) match(tokens ...lexer.TokenType) bool {
	if slices.ContainsFunc(tokens, p.check) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser[T]) assertMatch(token lexer.TokenType, errorMsg string) lexer.Token {
	if p.check(token) {
		return p.advance()
	}

	utils.Error(p.peek().Line, errorMsg)
	return p.peek()
}

func (p *Parser[T]) check(_type lexer.TokenType) bool {
	if p.isEndReached() {
		return false
	}

	return p.peek().Type == _type
}

func (p *Parser[T]) advance() lexer.Token {
	if !p.isEndReached() {
		p.current++
	}

	return p.peekBack()
}

func (p *Parser[T]) peek() lexer.Token {
	return p.Tokens[p.current]
}

func (p *Parser[T]) peekBack() lexer.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser[T]) isEndReached() bool {
	return p.peek().Type == lexer.EOF
}

func (p *Parser[T]) Parse() Expr[T] {
	return p.expression()
}

func (p *Parser[T]) expression() Expr[T] {
	return p.equality()
}

func (p *Parser[T]) equality() Expr[T] {
	expr := p.comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.peekBack()
		right := p.comparison()
		expr = &Binary[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}

	}

	return expr
}

func (p *Parser[T]) comparison() Expr[T] {
	expr := p.term()

	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.peekBack()
		right := p.term()
		expr = &Binary[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser[T]) term() Expr[T] {
	expr := p.factor()

	for p.match(lexer.PLUS, lexer.MINUS) {
		operator := p.peekBack()
		right := p.factor()
		expr = &Binary[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser[T]) factor() Expr[T] {
	expr := p.unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.peekBack()
		right := p.unary()
		expr = &Binary[T]{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser[T]) unary() Expr[T] {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.peekBack()
		right := p.unary()
		return &Unary[T]{
			Operator: operator,
			Right:    right,
		}
	}

	return p.primary()
}

func (p *Parser[T]) primary() Expr[T] {
	if p.match(lexer.FALSE) {
		return &Literal[T]{
			Value: false,
		}
	}

	if p.match(lexer.TRUE) {
		return &Literal[T]{
			Value: true,
		}
	}

	if p.match(lexer.NULL) {
		return &Literal[T]{
			Value: nil,
		}
	}

	if p.match(lexer.NUMBER, lexer.STRING) {
		return &Literal[T]{
			Value: p.peekBack().Literal,
		}
	}

	if p.match(lexer.LEFT_BRACKET) {
		expr := p.expression()

		p.assertMatch(lexer.RIGHT_BRACKET, "Expected closing bracket ')' but was not found")

		return &Grouping[T]{
			Expression: expr,
		}
	}

	utils.Error(p.peek().Line, "Expected expression but was not found")
	return &Literal[T]{}
}
