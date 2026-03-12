package ast

type visitor[T any] interface {
	VisitBinary(*Binary[T]) T
	VisitGrouping(*Grouping[T]) T
	VisitUnary(*Unary[T]) T
	VisitLiteral(*Literal[T]) T
}
