package lexer

import (
	"fmt"
	"strconv"

	"github.com/Moorad/carapace/utils"
)

type TokenType int

const (
	// Single char
	LEFT_CURLY TokenType = iota
	RIGHT_CURLY
	LEFT_BRACKET
	RIGHT_BRACKET
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR
	PERCENT

	// Multi char
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// KEYWORDS
	VAR
	AND
	OR
	TRUE
	FALSE
	NULL
	IF
	ELSE
	FOR
	WHILE
	FUNCTION
	RETURN
	CLASS
	THIS
	SUPER
	PRINT
	EOF
)

var reservedKeywords = map[string]TokenType{
	"var":      VAR,
	"and":      AND,
	"or":       OR,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"function": FUNCTION,
	"return":   RETURN,
	"class":    CLASS,
	"this":     THIS,
	"super":    SUPER,
	"print":    PRINT,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

type Scanner struct {
	Source  string
	Tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		Source:  source,
		current: 0,
		start:   0,
		line:    1,
	}
}

func (s *Scanner) Scan() error {

	for !s.isEndReached() {
		s.start = s.current

		s.scanToken()

	}

	s.addToken(EOF)

	return nil
}

func (s *Scanner) scanToken() {
	char := s.consume()

	switch char {
	// Single char
	case "(":
		s.addToken(LEFT_BRACKET)
	case ")":
		s.addToken(RIGHT_BRACKET)
	case "{":
		s.addToken(LEFT_CURLY)
	case "}":
		s.addToken(RIGHT_CURLY)
	case ",":
		s.addToken(COMMA)
	case ".":
		s.addToken(DOT)
	case "-":
		s.addToken(MINUS)
	case "+":
		s.addToken(PLUS)
	case ";":
		s.addToken(SEMICOLON)
	case "*":
		s.addToken(STAR)
	case "%":
		s.addToken(PERCENT)
	// Multi char conditionals
	case "!":
		if s.consumeIf("=") {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case "=":
		if s.consumeIf("=") {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case "<":
		if s.consumeIf("=") {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case ">":
		if s.consumeIf("=") {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case "/":
		if s.consumeIf("/") {
			for !s.isEndReached() && s.peek() != "\n" {
				s.consume()
			}
		} else {
			s.addToken(SLASH)
		}
	// Literals
	case "\"":
		for s.peek() != "\"" && !s.isEndReached() {
			if s.peek() == "\n" {
				s.line++
			}

			s.consume()
		}

		if s.isEndReached() {
			utils.Error(s.line, "Unterminated string")
		}

		s.consume()

		// Removing the trailing '"'
		s.addLiteralToken(STRING, s.Source[s.start+1:s.current-1])
	// Whitespace
	case "\n":
		s.line++
	case " ":
	case "\r":
	case "\t":
		break
	default:
		if s.isDigits(char) {
			s.parseNumber()
		} else if s.isAlpha(char) {
			s.parseIdentifier()
		} else {
			println(rune('3'), s.isDigits(char))
			utils.Error(s.line, fmt.Sprintf("Unexpected character found '%s'", char))
		}

	}
}

func (s *Scanner) addToken(_type TokenType) {
	lexeme := s.Source[s.start:s.current]

	if _type == EOF {
		lexeme = string(rune(0))
	}

	s.Tokens = append(s.Tokens, Token{Type: _type, Lexeme: lexeme, Line: s.line})
}

func (s *Scanner) addLiteralToken(_type TokenType, literalValue any) {
	lexeme := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, Token{Type: _type, Lexeme: lexeme, Literal: literalValue, Line: s.line})
}

func (s *Scanner) isEndReached() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) peek() string {
	return s.peekAhead(0)
}

func (s *Scanner) peekAhead(jump int) string {
	if s.current+jump >= len(s.Source) {
		return string(rune(0))
	}

	char := string(s.Source[s.current+jump])

	return char
}

func (s *Scanner) consume() string {
	char := string(s.Source[s.current])
	s.current++

	return char
}

func (s *Scanner) consumeIf(char string) bool {
	if s.isEndReached() {
		return false
	}

	if string(s.Source[s.current]) != char {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isDigits(char string) bool {
	r := []rune(char)[0]
	return r >= '0' && r <= '9'
}

func (s *Scanner) parseNumber() {
	for s.isDigits(s.peek()) {
		s.consume()
	}

	if s.peek() == "." && s.isDigits(s.peekAhead(1)) {
		s.consume()

		for s.isDigits(s.peek()) {
			s.consume()
		}
	}

	num, _ := strconv.ParseFloat(s.Source[s.start:s.current], 32)
	s.addLiteralToken(NUMBER, float32(num))

}

func (s *Scanner) isAlpha(char string) bool {
	r := []rune(char)[0]
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == 'c'

}

func (s *Scanner) isAlphaNumeric(char string) bool {
	return s.isAlpha(char) || s.isDigits(char)
}

func (s *Scanner) parseIdentifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.consume()
	}

	identifierName := s.Source[s.start:s.current]
	tokenType, found := reservedKeywords[identifierName]

	if !found {
		tokenType = IDENTIFIER
	}

	s.addToken(tokenType)

}
