package ast

import "github.com/Moorad/carapace/lexer"

type Expr[T any] interface {
	Accept(visitor[T]) T
}

type Binary[T any] struct {
	Left     Expr[T]
	Operator lexer.Token
	Right    Expr[T]
}

func (b *Binary[T]) Accept(v visitor[T]) T {
	return v.visitBinary(b)
}

type Unary[T any] struct {
	Operator lexer.Token
	Right    Expr[T]
}

func (u *Unary[T]) Accept(v visitor[T]) T {
	return v.visitUnary(u)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (g *Grouping[T]) Accept(v visitor[T]) T {
	return v.visitGrouping(g)
}

type Literal[T any] struct {
	Value any
}

func (l *Literal[T]) Accept(v visitor[T]) T {
	return v.visitLiteral(l)
}
