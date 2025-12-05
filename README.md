# PRATT PARSING

Pratt parsing is a parsing approach that builds the AST by leveraing operator, precedence and binding powers

Compared to conventional recursive descent parsers, Pratt parsers offer a more declarative and intuitive way to define and parse expressions

## Binding power

Binding power is how tightly a token binds to it's neighboring tokens. Another way to think of this is attraction. <br/>
`A + token` will have a lower binding power then `A * token`
