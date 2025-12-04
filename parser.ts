type TokenType =
        | 'NUMBER'
        | 'PLUS'
        | 'MINUS'
        | 'MULTIPLY'
        | 'DIVIDE'
        | 'LPAREN'
        | 'RPAREN'
        | 'EOF'

interface Token {
        type: TokenType,
        value: string | number
}

function tokenize(input: string) {
        const tokens: Token[] = []
        let i = 0

        while (i < input.length) {
                // skipt whitespaces
                const char = input[i]
                if (/\s/.test(char)) {
                        i++
                        continue
                }

                if (/\d/.test(char)) {
                        let num = ''
                        while (i < input.length && /[\d.]/.test(input[i])) {
                                num += input[i]
                                i++
                        }

                        tokens.push({ type: 'NUMBER', value: parseFloat(num) })
                        continue
                }

        }
}




