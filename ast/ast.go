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
	return v.VisitBinary(b)
}

type Unary[T any] struct {
	Operator lexer.Token
	Right    Expr[T]
}

func (u *Unary[T]) Accept(v visitor[T]) T {
	return v.VisitUnary(u)
}

type Grouping[T any] struct {
	Expression Expr[T]
}

func (g *Grouping[T]) Accept(v visitor[T]) T {
	return v.VisitGrouping(g)
}

type Literal[T any] struct {
	Value any
}

func (l *Literal[T]) Accept(v visitor[T]) T {
	return v.VisitLiteral(l)
}
