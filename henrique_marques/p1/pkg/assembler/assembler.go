package assembler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"p1/pkg/assembler/lexer"
)

const (
	TOKEN_SECTION = "SECTION"
	TOKEN_EOF     = "EOF"
	TOKEN_INSTR   = "INSTRUCTION"
	TOKEN_NUMBER  = "NUMBER"
	TOKEN_VAR     = "VARIABLE"
	TOKEN_DEFINE  = "DEFINE"
	TOKEN_UNKNOWN = "UNKNOWN"
)

var (
	Instructions = map[string]uint8{
		"NOP": 0x00, "STA": 0x10, "LDA": 0x20, "ADD": 0x30,
		"OR": 0x40, "AND": 0x50, "NOT": 0x60, "JMP": 0x80,
		"JN": 0x90, "JZ": 0xA0, "HLT": 0xF0,
	}
)

type Assembler struct {
	Tokens  []lexer.Token
	PC      uint8         
	Output  []uint8
	Labels  map[string]uint8
	StartPC uint8
}

func NewAssembler(tokens []lexer.Token) *Assembler {
	return &Assembler{
		Tokens:  tokens,
		PC:      0,
		Labels:  make(map[string]uint8),
		StartPC: 0,
	}
}

// FirstPass calcula os endereços dos rótulos e atualiza o PC conforme as diretivas.
// Nesta passagem, o PC é incrementado de forma contínua, respeitando a ordem das seções.
func (a *Assembler) FirstPass() error {
	currentSection := "CODE"
	for i := 0; i < len(a.Tokens); i++ {
		token := a.Tokens[i]
		if token.Tipo == TOKEN_SECTION {
			currentSection = strings.ToUpper(token.Valor)
			continue
		}
		if currentSection == "CODE" {
			switch token.Tipo {
			case TOKEN_INSTR, TOKEN_NUMBER, TOKEN_VAR:
				a.PC += 2
			case TOKEN_DEFINE:
				if token.Valor == "ORG" {
					i++
					if i >= len(a.Tokens) {
						return fmt.Errorf("esperado número após ORG")
					}
					value, err := parseNumber(a.Tokens[i].Valor)
					if err != nil {
						return fmt.Errorf("número inválido após ORG: %s", a.Tokens[i].Valor)
					}
					a.PC = uint8(value)
					a.StartPC = uint8(value)
				}
			}
		} else if currentSection == "DATA" {
			if token.Tipo == TOKEN_DEFINE && token.Valor == "ORG" {
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após ORG")
				}
				value, err := parseNumber(a.Tokens[i].Valor)
				if err != nil {
					return fmt.Errorf("número inválido após ORG: %s", a.Tokens[i].Valor)
				}
				a.PC = uint8(value)
				continue
			}

			if token.Tipo == TOKEN_VAR {
				a.Labels[token.Valor] = a.PC
				continue
			}

			if token.Tipo == TOKEN_DEFINE && token.Valor == "DB" {
				i++
				if i >= len(a.Tokens) {
					return fmt.Errorf("esperado número após DB")
				}
				a.PC += 2
			}
		}
	}
	return nil
}

// SecondPass gera o buffer de memória (512 bytes) com base nos tokens.
// Na seção CODE, escreve as instruções de forma sequencial (usando um PC local).
// Na seção DATA, escreve os valores nos endereços definidos pelos rótulos.
func (a *Assembler) SecondPass() error {
	mem := make([]uint8, 512)
	pcCode := a.StartPC
	currentSection := "CODE"
	var currentVar string

	for i := 0; i < len(a.Tokens); i++ {
		token := a.Tokens[i]
		if token.Tipo == TOKEN_SECTION {
			currentSection = strings.ToUpper(token.Valor)
			continue
		}

		if currentSection == "CODE" {
			switch token.Tipo {
			case TOKEN_INSTR:
				opcode, ok := Instructions[token.Valor]
				if !ok {
					return fmt.Errorf("instrução desconhecida: %s", token.Valor)
				}
				realAddr := pcCode * 2
				mem[realAddr] = opcode
				mem[realAddr+1] = 0x00
				pcCode += 1
			case TOKEN_NUMBER:
				value, err := parseNumber(token.Valor)
				if err != nil {
					return fmt.Errorf("número inválido: %s", token.Valor)
				}
				realAddr := pcCode * 2
				mem[realAddr] = uint8(value)
				mem[realAddr+1] = 0x00
				pcCode += 1
			case TOKEN_VAR:
				addr, ok := a.Labels[token.Valor]
				if !ok {
					return fmt.Errorf("label não definida: %s", token.Valor)
				}
				realAddr := pcCode * 2
				mem[realAddr] = addr
				mem[realAddr+1] = 0x00
				pcCode += 1
			case TOKEN_DEFINE:
				if token.Valor == "ORG" {
					i++
					if i >= len(a.Tokens) {
						return fmt.Errorf("esperado número após ORG")
					}
					value, err := parseNumber(a.Tokens[i].Valor)
					if err != nil {
						return fmt.Errorf("número inválido após ORG: %s", a.Tokens[i].Valor)
					}
					pcCode = uint8(value)
				}
			}
		} else if currentSection == "DATA" {
			if token.Tipo == TOKEN_VAR {
				currentVar = token.Valor
			} else if token.Tipo == TOKEN_DEFINE {
				if token.Valor == "DB" {
					i++
					if i >= len(a.Tokens) {
						return fmt.Errorf("esperado número após DB")
					}
					value, err := parseNumber(a.Tokens[i].Valor)
					if err != nil {
						return fmt.Errorf("número inválido após DB: %s", a.Tokens[i].Valor)
					}
					addr, ok := a.Labels[currentVar]
					if !ok {
						return fmt.Errorf("label não definida para variável: %s", currentVar)
					}
					realAddr := addr * 2
					mem[realAddr] = uint8(value)
					mem[realAddr+1] = 0x00
				} else if token.Valor == "ORG" {
					i++
					if i >= len(a.Tokens) {
						return fmt.Errorf("esperado número após ORG")
					}
					value, err := parseNumber(a.Tokens[i].Valor)
					if err != nil {
						return fmt.Errorf("número inválido após ORG: %s", a.Tokens[i].Valor)
					}
					_ = value
				}
			}
		}
	}

	a.Output = mem
	return nil
}


func parseNumber(s string) (uint64, error) {
	return strconv.ParseUint(s, 16, 8)
}

// WriteMEM grava o arquivo .mem com um cabeçalho fixo (4 bytes) e preenche até 516 bytes.
func (a *Assembler) WriteMEM(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := []uint8{0x03, 0x4E, 0x44, 0x52} // Cabeçalho fixo
	output := append(header, a.Output...)

	for len(output) < 516 {
		output = append(output, 0x00)
	}

	_, err = file.Write(output)
	if err != nil {
		return err
	}

	return nil
}
