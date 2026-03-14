package interpreter

import (
	"fmt"

	"github.com/Moorad/carapace/ast"
	"github.com/Moorad/carapace/lexer"
	"github.com/Moorad/carapace/utils"
)

type Interpreter struct{}

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if value, ok := value.(bool); ok {
		return value
	}

	return true
}

func isEqual(left any, right any) bool {
	return left == right
}

func assertCastUnaryNumber(right any, operator lexer.Token) float32 {
	if num, ok := right.(float32); ok {
		return num
	}

	utils.Error(operator.Line, fmt.Sprintf("right value for operand '%v' must be a number", operator.Lexeme))

	return 0.0
}

func assertCastBinaryNumber(left any, right any, operator lexer.Token) (float32, float32) {
	leftNum, isLeftOk := left.(float32)
	rightNum, isRightOk := right.(float32)

	if isLeftOk && isRightOk {
		return leftNum, rightNum
	}

	utils.Error(operator.Line, fmt.Sprintf("value for %s operand must be both numbers", operator.Literal))

	return 0.0, 0.0
}

func (i *Interpreter) Interpret(expr ast.Expr[any]) {
	out := i.evaluate(expr)
	println(stringify(out))
}

func (i *Interpreter) evaluate(expr ast.Expr[any]) any {
	return expr.Accept(i)
}

func (i *Interpreter) VisitLiteral(expr *ast.Literal[any]) any {
	return expr.Value
}

func (i *Interpreter) VisitGrouping(expr *ast.Grouping[any]) any {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnary(expr *ast.Unary[any]) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case lexer.BANG:
		return !isTruthy(right)
	case lexer.MINUS:
		return -assertCastUnaryNumber(right, expr.Operator)
	}

	// Unreachable
	return nil
}

func (i *Interpreter) VisitBinary(expr *ast.Binary[any]) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case lexer.GREATER:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left > right
	case lexer.GREATER_EQUAL:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left >= right
	case lexer.LESS:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left < right
	case lexer.LESS_EQUAL:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left <= right
	case lexer.EQUAL_EQUAL:
		return isEqual(left, right)
	case lexer.BANG_EQUAL:
		return !isEqual(left, right)
	case lexer.MINUS:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left - right
	case lexer.SLASH:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left / right
	case lexer.STAR:
		left, right := assertCastBinaryNumber(left, right, expr.Operator)
		return left * right
	case lexer.PLUS:
		leftNum, isLeftNum := left.(float32)
		rightNum, isRightNum := right.(float32)

		if isLeftNum && isRightNum {
			return leftNum + rightNum
		}

		leftStr, isLeftStr := left.(string)
		rightStr, isRightStr := right.(string)

		if isLeftStr && isRightStr {
			return leftStr + rightStr
		}

		utils.Error(expr.Operator.Line, "can only add two numbers or concatenate two strings but found different types")
	}

	// Unreachable
	return nil
}

func stringify(t any) string {
	if t == nil {
		return "null"
	}

	return fmt.Sprint(t)
}
