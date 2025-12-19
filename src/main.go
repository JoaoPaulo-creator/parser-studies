package main

import (
	"custom_parser/src/lexer"
	"custom_parser/src/parser"

	// "fmt"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./examples/06.lang")
	source := string(bytes)

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	litter.Dump(ast)
}
