package ast

type visitor[T any] interface {
	visitBinary(*Binary[T]) T
	visitGrouping(*Grouping[T]) T
	visitUnary(*Unary[T]) T
	visitLiteral(*Literal[T]) T
}
