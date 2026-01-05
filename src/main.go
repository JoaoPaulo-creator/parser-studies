package main

import (
	"custom_parser/src/lexer"
	"custom_parser/src/parser"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	bytes, _ := os.ReadFile("./examples/07.lang")
	source := string(bytes)

	tokens := lexer.Tokenize(source)
	ast := parser.Parse(tokens)
	litter.Dump(ast)
}
