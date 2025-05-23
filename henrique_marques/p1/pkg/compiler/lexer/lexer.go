package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType string

type Token struct {
	Tipo  TokenType
	Valor string
}

const (
	TOKEN_PROGRAMA  TokenType = "PROGRAMA"
	TOKEN_LABEL     TokenType = "LABEL"
	TOKEN_INICIO    TokenType = "INICIO"
	TOKEN_FIM       TokenType = "FIM"
	TOKEN_VAR       TokenType = "VAR"
	TOKEN_NUM       TokenType = "NUM"
	TOKEN_OP        TokenType = "OP"
	TOKEN_ATRIB     TokenType = "="
	TOKEN_ABREPAR   TokenType = "("
	TOKEN_FECHAPAR  TokenType = ")"
	TOKEN_NEWLINE   TokenType = "\n"
	TOKEN_EOF       TokenType = "EOF"
)

var operadores = "+-*/"

func isLetter(r rune) bool {
	return unicode.IsLetter(r)
}

func isHexDigit(r rune) bool {
	return unicode.IsDigit(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isVarChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func Lex(code string) ([]Token, error) {
	tokens := []Token{}
	i := 0
	runes := []rune(code)
	for i < len(runes) {
		c := runes[i]
		
		if unicode.IsSpace(c) && c != '\n' {
			i++
			continue
		}

		if c == '\n' {
			tokens = append(tokens, Token{TOKEN_NEWLINE, "\\n"})
			i++
			continue
		}

		if c == '=' {
			tokens = append(tokens, Token{TOKEN_ATRIB, "="})
			i++
			continue
		}

		if strings.ContainsRune(operadores, c) {
			tokens = append(tokens, Token{TOKEN_OP, string(c)})
			i++
			continue
		}

		if c == '(' {
			tokens = append(tokens, Token{TOKEN_ABREPAR, "("})
			i++
			continue
		}

		if c == ')' {
			tokens = append(tokens, Token{TOKEN_FECHAPAR, ")"})
			i++
			continue
		}

		if c == '"' {
			j := i + 1
			for j < len(runes) && runes[j] != '"' {
				j++
			}
			if j >= len(runes) {
				return nil, fmt.Errorf("string não terminada")
			}
			valor := string(runes[i+1 : j])
			tokens = append(tokens, Token{TOKEN_LABEL, valor})
			i = j + 1
			continue
		}

		if isLetter(c) {
			j := i
			for j < len(runes) && isVarChar(runes[j]) {
				j++
			}
			palavra := string(runes[i:j])
			switch palavra {
			case "PROGRAMA":
				tokens = append(tokens, Token{TOKEN_PROGRAMA, palavra})
			case "INICIO":
				tokens = append(tokens, Token{TOKEN_INICIO, palavra})
			case "FIM":
				tokens = append(tokens, Token{TOKEN_FIM, palavra})
			default:
				tokens = append(tokens, Token{TOKEN_VAR, palavra})
			}
			i = j
			continue
		}

		if isHexDigit(c) {
			j := i
			for j < len(runes) && isHexDigit(runes[j]) {
				j++
			}
			num := string(runes[i:j])
			tokens = append(tokens, Token{TOKEN_NUM, num})
			i = j
			continue
		}

		return nil, fmt.Errorf("caractere inesperado: %c", c)
	}

	tokens = append(tokens, Token{TOKEN_EOF, ""})
	return tokens, nil
}
