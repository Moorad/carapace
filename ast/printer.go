package ast

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr[string]) string {
	var str strings.Builder

	str.WriteString("(")
	str.WriteString(name)
	for _, expr := range exprs {
		str.WriteString(" ")
		str.WriteString(expr.Accept(p))
	}
	str.WriteString(")")

	return str.String()
}

func (p *AstPrinter) VisitBinary(expr *Binary[string]) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitUnary(expr *Unary[string]) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitGrouping(expr *Grouping[string]) string {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteral(expr *Literal[string]) string {
	if expr.Value == nil {
		return "null"
	}

	return fmt.Sprintf("%v", expr.Value)
}
