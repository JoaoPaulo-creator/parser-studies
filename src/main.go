package main

import (
	"custom_parser/src/lexer"
	// "fmt"
	"os"
)

func main() {
	bytes, _ := os.ReadFile("./examples/00.lang")
	source := string(bytes)

	tokens := lexer.Tokenize(source)
	for _, token := range tokens {
		token.Debug()
	}
	// fmt.Println(source)
}
