# PRATT PARSING

Pratt parsing is a parsing approach that builds the AST by leveraing operator, precedence and binding powers

Compared to conventional recursive descent parsers, Pratt parsers offer a more declarative and intuitive way to define and parse expressions

## Binding power

Binding power is how tightly a token binds to it's neighboring tokens. Another way to think of this is attraction. <br/>
A `+` token will have a lower binding power then `*` token

## Null Denotation (NUD)

```go
type nud_handler func(p *parser) ast.Expr
```

A token which has a LUD handler, means it expects nothing to it's left.
Common examples of this type of token are `prefix` & unary expression

## Left Denotation (LED):

```go
type led_handler func (p *parser, left ast.Expr, bp binding_power) ast.Expr
```

Tokens which have a LED expect to be between or after some other expression to it's left. Examples of this type of handler include binary expression and all `infix` expressions.

`Postfix` expressions also fall under the LED handler.

## Lookup Tables
By using lookup tables along with NUD/LED handler functions, we can create the parser almost entirely without having to manage the recursion ourselves.
Here is an example of what our lookup tables look line:

```go

type stmt_handler func (p *parser) ast.Stmt
type nud_handler 	func (p *parser) ast.Expr
type led_handler 	func (p *parser, left ast.Expr, bp bindind_power) ast.Expr

type stmt_lookup map[lexer.TokendKind] stmt_handler
type nud_lookup  map[lexer.TokendKind] nud_handler
type led_lookup  map[lexer.TokendKind] led_handler
type bp_lookup 	 map[lexer.TokendKind] binding_power

var bp_lu 	= bp_lookup{}
var nud_lu 	= nud_lookup{}
var led_lu 	= led_lookup{}
var stmt_lu = stmt_lookup{}
```
