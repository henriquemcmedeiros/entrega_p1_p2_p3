package generator

import (
	"fmt"
	"p1/pkg/compiler/lexer"
	"p1/pkg/compiler/parser"
	"strconv"
)

type ASMProgram struct {
	Code []string
	Data []string
}

var tmpCount = 0
var constSet = map[string]bool{}

func resetState() {
	tmpCount = 0
	constSet = map[string]bool{}
}

func newTmp() string {
	tmp := fmt.Sprintf("TMP%d", tmpCount)
	tmpCount++
	return tmp
}

func GenerateASM(instrucoes []parser.Instrucao) ASMProgram {
	resetState()
	prog := ASMProgram{
		Code: []string{".CODE", "ORG 00"},
		Data: []string{".DATA", "ORG 20"},
	}

	varsUsadas := map[string]bool{}

	for _, inst := range instrucoes {
		stack := []string{}
		for _, tok := range inst.Expr {
			switch tok.Tipo {
			case lexer.TOKEN_NUM:
				constLabel := "CONST_" + tok.Valor
				if !constSet[constLabel] {
					prog.Data = append(prog.Data, fmt.Sprintf("%s DB %s", constLabel, tok.Valor))
					constSet[constLabel] = true
				}
				stack = append(stack, constLabel)

			case lexer.TOKEN_VAR:
				varsUsadas[tok.Valor] = true
				stack = append(stack, tok.Valor)

			case lexer.TOKEN_OP:
				if len(stack) < 2 {
					panic("expressão mal formada")
				}
				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				tmp := newTmp()
				prog.Data = append(prog.Data, fmt.Sprintf("%s DB 00", tmp))

				if tok.Valor == "*" {
					prog.Code = append(prog.Code, fmt.Sprintf("LDA %s", left))
					value := 0
					if valStr := right; len(valStr) > 6 && valStr[:6] == "CONST_" {
						hex := valStr[6:]
						parsed, err := strconv.ParseUint(hex, 16, 8)
						if err == nil {
							value = int(parsed)
						}
					}
					for i := 1; i < value; i++ {
						prog.Code = append(prog.Code, fmt.Sprintf("ADD %s", left))
					}
				} else {
					switch tok.Valor {
					case "+":
						prog.Code = append(prog.Code, fmt.Sprintf("LDA %s", left))
						prog.Code = append(prog.Code, fmt.Sprintf("ADD %s", right))
					case "-":
						negTmp := newTmp()
						prog.Data = append(prog.Data, fmt.Sprintf("%s DB 00", negTmp))
						
						prog.Code = append(prog.Code, fmt.Sprintf("LDA %s", right))
						prog.Code = append(prog.Code, "NOT")
						prog.Code = append(prog.Code, "ADD CONST_01")
						prog.Code = append(prog.Code, fmt.Sprintf("STA %s", negTmp))
						
						prog.Code = append(prog.Code, fmt.Sprintf("LDA %s", left))
						prog.Code = append(prog.Code, fmt.Sprintf("ADD %s", negTmp))
						
						if !constSet["CONST_01"] {
							prog.Data = append(prog.Data, "CONST_01 DB 01")
							constSet["CONST_01"] = true
						}
					case "/":
						prog.Code = append(prog.Code, "; DIV não suportado")
					}
				}
				prog.Code = append(prog.Code, fmt.Sprintf("STA %s", tmp))
				stack = append(stack, tmp)
				
			}
		}

		if len(stack) != 1 {
			panic("erro interno: pilha final da expressão não tem 1 item")
		}
		prog.Code = append(prog.Code, fmt.Sprintf("LDA %s", stack[0]))
		prog.Code = append(prog.Code, fmt.Sprintf("STA %s", inst.Var))
	}

	for v := range varsUsadas {
		prog.Data = append(prog.Data, fmt.Sprintf("%s DB 00", v))
	}

	for _, inst := range instrucoes {
		if !varsUsadas[inst.Var] {
			prog.Data = append(prog.Data, fmt.Sprintf("%s DB 00", inst.Var))
			varsUsadas[inst.Var] = true
		}
	}

	prog.Code = append(prog.Code, "HLT")
	return prog
}