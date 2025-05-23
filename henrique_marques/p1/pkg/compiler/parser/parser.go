package parser

import (
	"fmt"
	"p1/pkg/compiler/lexer"
)

type Instrucao struct {
	Var  string
	Expr []lexer.Token
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) current() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Tipo: lexer.TOKEN_EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	tok := p.current()
	p.pos++
	return tok
}

func (p *Parser) match(t lexer.TokenType) bool {
	if p.current().Tipo == t {
		p.advance()
		return true
	}
	return false
}

var precedencia = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func (p *Parser) ParsePrograma() ([]Instrucao, error) {
	instrucoes := []Instrucao{}

	if !p.match(lexer.TOKEN_PROGRAMA) {
		return nil, fmt.Errorf("Esperado 'PROGRAMA'")
	}
	if !p.match(lexer.TOKEN_LABEL) {
		return nil, fmt.Errorf("Esperado nome do programa")
	}
	if !p.match(lexer.TOKEN_NEWLINE) {
		return nil, fmt.Errorf("Esperado quebra de linha após label")
	}
	if !p.match(lexer.TOKEN_INICIO) || !p.match(lexer.TOKEN_NEWLINE) {
		return nil, fmt.Errorf("Esperado 'INICIO' na linha seguinte")
	}

	for p.current().Tipo != lexer.TOKEN_FIM && p.current().Tipo != lexer.TOKEN_EOF {
		inst, err := p.parseInstrucao()
		if err != nil {
			return nil, err
		}
		instrucoes = append(instrucoes, inst)
	}

	if !p.match(lexer.TOKEN_FIM) {
		return nil, fmt.Errorf("Esperado 'FIM'")
	}
	return instrucoes, nil
}

func (p *Parser) parseInstrucao() (Instrucao, error) {
	var nome string
	if p.current().Tipo != lexer.TOKEN_VAR {
		return Instrucao{}, fmt.Errorf("Esperado nome da variável")
	}
	nome = p.advance().Valor

	if !p.match(lexer.TOKEN_ATRIB) {
		return Instrucao{}, fmt.Errorf("Esperado '=' após variável")
	}

	expr, err := p.parseExp()
	if err != nil {
		return Instrucao{}, err
	}

	if !p.match(lexer.TOKEN_NEWLINE) {
		return Instrucao{}, fmt.Errorf("Esperado quebra de linha após expressão")
	}

	return Instrucao{Var: nome, Expr: expr}, nil
}

func (p *Parser) parseExp() ([]lexer.Token, error) {
	saida := []lexer.Token{}
	pilha := []lexer.Token{}

	for {
		tok := p.current()
		if tok.Tipo == lexer.TOKEN_NUM || tok.Tipo == lexer.TOKEN_VAR {
			saida = append(saida, tok)
			p.advance()
		} else if tok.Tipo == lexer.TOKEN_OP {
			for len(pilha) > 0 {
				top := pilha[len(pilha)-1]
				if top.Tipo == lexer.TOKEN_OP && precedencia[top.Valor] >= precedencia[tok.Valor] {
					saida = append(saida, top)
					pilha = pilha[:len(pilha)-1]
				} else {
					break
				}
			}
			pilha = append(pilha, tok)
			p.advance()
		} else if tok.Tipo == lexer.TOKEN_ABREPAR {
			pilha = append(pilha, tok)
			p.advance()
		} else if tok.Tipo == lexer.TOKEN_FECHAPAR {
			for len(pilha) > 0 && pilha[len(pilha)-1].Tipo != lexer.TOKEN_ABREPAR {
				saida = append(saida, pilha[len(pilha)-1])
				pilha = pilha[:len(pilha)-1]
			}
			if len(pilha) == 0 || pilha[len(pilha)-1].Tipo != lexer.TOKEN_ABREPAR {
				return nil, fmt.Errorf("Parêntese não balanceado")
			}
			pilha = pilha[:len(pilha)-1] // descarta o "("
			p.advance()
		} else {
			break
		}
	}

	for len(pilha) > 0 {
		if pilha[len(pilha)-1].Tipo == lexer.TOKEN_ABREPAR {
			return nil, fmt.Errorf("Parêntese não fechado")
		}
		saida = append(saida, pilha[len(pilha)-1])
		pilha = pilha[:len(pilha)-1]
	}

	return saida, nil
}