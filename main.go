package main

import (
	"fmt"
	"os"

	"github.com/Moorad/carapace/ast"
	"github.com/Moorad/carapace/lexer"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage:\n\tcarapace [file_path]")
		os.Exit(1)
	}

	filePath := os.Args[1]

	content, err := runFile(filePath)

	if err != nil {
		panic(fmt.Sprintf("Failed to read %s:\n%v", filePath, err))
	}

	source := string(content)
	scanner := lexer.NewScanner(source)

	err = scanner.Scan()

	if err != nil {
		panic(fmt.Sprintf("Syntax Error:\n%v", err))
	}

	printer := ast.AstPrinter{}
	parser := ast.Parser[string]{
		Tokens: scanner.Tokens,
	}

	syntaxTree := parser.Parse()

	println(syntaxTree.Accept(&printer))

}

func runFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
