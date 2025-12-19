// lexer.go
package lexer

import (
	"fmt"
	"unicode"
)

type lexer struct {
	source string
	pos    int
	line   int
	Tokens []Token
}

func Tokenize(source string) []Token {
	lex := &lexer{
		source: source,
		pos:    0,
		line:   1,
		Tokens: make([]Token, 0),
	}

	for !lex.atEOF() {
		lex.scanToken()
	}

	lex.push(newUniqueToken(EOF, "EOF"))
	return lex.Tokens
}

func (lex *lexer) scanToken() {
	ch := lex.peek()

	// Skip whitespace
	if unicode.IsSpace(rune(ch)) {
		lex.skipWhitespace()
		return
	}

	// Comments
	if ch == '/' && lex.peekNext() == '/' {
		lex.skipComment()
		return
	}

	// String literals
	if ch == '"' {
		lex.scanString()
		return
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) {
		lex.scanNumber()
		return
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		lex.scanIdentifier()
		return
	}

	// Multi-character operators (check these before single-char)
	switch ch {
	case '=':
		if lex.peekNext() == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(EQUALS, "=="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(ASSIGNMENT, "="))
		return

	case '!':
		if lex.peekNext() == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(NOT_EQUALS, "!="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(NOT, "!"))
		return

	case '<':
		if lex.peekNext() == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(LESS_EQUALS, "<="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(LESS, "<"))
		return

	case '>':
		if lex.peekNext() == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(GREATER_EQUALS, ">="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(GREATER, ">"))
		return

	case '|':
		if lex.peekNext() == '|' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(OR, "||"))
			return
		}

	case '&':
		if lex.peekNext() == '&' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(AND, "&&"))
			return
		}

	case '.':
		if lex.peekNext() == '.' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(DOT_DOT, ".."))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(DOT, "."))
		return

	case '+':
		next := lex.peekNext()
		if next == '+' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(PLUS_PLUS, "++"))
			return
		}
		if next == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(PLUS_EQUALS, "+="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(PLUS, "+"))
		return

	case '-':
		next := lex.peekNext()
		if next == '-' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(MINUS_MINUS, "--"))
			return
		}
		if next == '=' {
			lex.advance()
			lex.advance()
			lex.push(newUniqueToken(MINUS_EQUALS, "-="))
			return
		}
		lex.advance()
		lex.push(newUniqueToken(DASH, "-"))
		return

	// case '?':
	// 	if lex.peekNext() == '?' && lex.peekAhead(2) == '=' {
	// 		lex.advance()
	// 		lex.advance()
	// 		lex.advance()
	// 		lex.push(newUniqueToken(NULLISH_ASSIGNMENT, "??="))
	// 		return
	// 	}
	// 	lex.advance()
	// 	lex.push(newUniqueToken(QUESTION, "?"))
	// 	return

	// Single-character tokens
	case '[':
		lex.advance()
		lex.push(newUniqueToken(OPEN_BRACKET, "["))
		return
	case ']':
		lex.advance()
		lex.push(newUniqueToken(CLOSE_BRACKET, "]"))
		return
	case '{':
		lex.advance()
		lex.push(newUniqueToken(OPEN_CURLY, "{"))
		return
	case '}':
		lex.advance()
		lex.push(newUniqueToken(CLOSE_CURLY, "}"))
		return
	case '(':
		lex.advance()
		lex.push(newUniqueToken(OPEN_PAREN, "("))
		return
	case ')':
		lex.advance()
		lex.push(newUniqueToken(CLOSE_PAREN, ")"))
		return
	case ';':
		lex.advance()
		lex.push(newUniqueToken(SEMI_COLON, ";"))
		return
	case ':':
		lex.advance()
		lex.push(newUniqueToken(COLON, ":"))
		return
	case ',':
		lex.advance()
		lex.push(newUniqueToken(COMMA, ","))
		return
	case '/':
		lex.advance()
		lex.push(newUniqueToken(SLASH, "/"))
		return
	case '*':
		lex.advance()
		lex.push(newUniqueToken(STAR, "*"))
		return
	case '%':
		lex.advance()
		lex.push(newUniqueToken(PERCENT, "%"))
		return
	}

	panic(fmt.Sprintf("lexer error: unexpected character '%c' at position %d", ch, lex.pos))
}

func (lex *lexer) scanString() {
	start := lex.pos
	lex.advance() // Skip opening quote

	for !lex.atEOF() && lex.peek() != '"' {
		lex.advance()
	}

	if lex.atEOF() {
		panic("lexer error: unterminated string literal")
	}

	lex.advance() // Skip closing quote
	value := lex.source[start:lex.pos]
	lex.push(newUniqueToken(STRING, value))
}

func (lex *lexer) scanNumber() {
	start := lex.pos

	// Scan integer part
	for !lex.atEOF() && unicode.IsDigit(rune(lex.peek())) {
		lex.advance()
	}

	// Check for decimal part
	if !lex.atEOF() && lex.peek() == '.' && unicode.IsDigit(rune(lex.peekNext())) {
		lex.advance() // Skip '.'
		for !lex.atEOF() && unicode.IsDigit(rune(lex.peek())) {
			lex.advance()
		}
	}

	value := lex.source[start:lex.pos]
	lex.push(newUniqueToken(NUMBER, value))
}

func (lex *lexer) scanIdentifier() {
	start := lex.pos

	for !lex.atEOF() {
		ch := lex.peek()
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			lex.advance()
		} else {
			break
		}
	}

	value := lex.source[start:lex.pos]

	// Check if it's a reserved keyword
	if kind, found := reserved_lu[value]; found {
		lex.push(newUniqueToken(kind, value))
	} else {
		lex.push(newUniqueToken(IDENTIFIER, value))
	}
}

func (lex *lexer) skipWhitespace() {
	for !lex.atEOF() && unicode.IsSpace(rune(lex.peek())) {
		if lex.peek() == '\n' {
			lex.line++
		}
		lex.advance()
	}
}

func (lex *lexer) skipComment() {
	// Skip until end of line
	for !lex.atEOF() && lex.peek() != '\n' {
		lex.advance()
	}
	if !lex.atEOF() {
		lex.advance() // Skip the newline
		lex.line++
	}
}

// Helper methods
func (lex *lexer) peek() byte {
	if lex.atEOF() {
		return 0
	}
	return lex.source[lex.pos]
}

func (lex *lexer) peekNext() byte {
	if lex.pos+1 >= len(lex.source) {
		return 0
	}
	return lex.source[lex.pos+1]
}

func (lex *lexer) peekAhead(n int) byte {
	if lex.pos+n >= len(lex.source) {
		return 0
	}
	return lex.source[lex.pos+n]
}

func (lex *lexer) advance() {
	lex.pos++
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) atEOF() bool {
	return lex.pos >= len(lex.source)
}
