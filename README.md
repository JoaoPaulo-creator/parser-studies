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
