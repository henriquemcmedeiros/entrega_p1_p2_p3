package encoder

import (
	"fmt"
	"log"
	"os"
)

const (
	TOTAL_SIZE = 516

	NOP = 0x00
	STA = 0x10
	LDA = 0x20
	ADD = 0x30
	OR  = 0x40
	AND = 0x50
	NOT = 0x60
	JMP = 0x80
	JN  = 0x90
	JZ  = 0xA0
	HLT = 0xF0
)

func flagZero(AC int) bool {
	return AC == 0x00
}

func flagNeg(AC int) bool {
	return AC & 0x80 != 0
}

func RunBinary(caminhoArquivo string) {
	AC := 0
	PC := 0x04

	memory, err := os.ReadFile(caminhoArquivo)
	if err != nil {
		log.Fatalf("Não foi possível ler o arquivo!")
		return
	}

	posicao := 0

	for memory[PC] != HLT && PC <= 0xFF {
		fmt.Printf("AC: %2x PC: %2x FZ: %5t FN: %5t INSTRUCAO: %2x CONTEUDO: %2x\n", AC & 0xFF, PC, flagZero(AC), flagNeg(AC), memory[PC], memory[PC+2])

		switch memory[PC] {
			case STA:
				PC += 2
				posicao = int(memory[PC]) * 2 + 4
				memory[posicao] = byte(AC)
				PC += 2
			case LDA:
				PC += 2
				posicao = int(memory[PC]) * 2 + 4
				AC = int(memory[posicao])
				PC += 2
			case ADD:
				PC += 2
				posicao = int(memory[PC]) * 2 + 4
				AC += int(memory[posicao])
				PC += 2
			case OR:
				PC += 2
				posicao = int(memory[PC]) * 2 + 4
				AC |= int(memory[posicao])
				PC += 2
			case AND:
				PC += 2
				posicao = int(memory[PC]) * 2 + 4
				AC &= int(memory[posicao])
				PC += 2
			case NOT:
				AC = ^AC
				PC += 2
			case JMP:
				PC += 2
				PC = int(memory[PC]) * 2 + 4
			case JN:
				PC += 2
				if flagNeg(AC) {
					PC = int(memory[PC]) * 2 + 4
				} else {
					PC += 2
				}
			case JZ:
				PC += 2
				if flagZero(AC) {
					PC = int(memory[PC]) * 2 + 4
				} else {
					PC += 2
				}
			default:
				PC += 2
		}
	}

	fmt.Println("========== Retorno de Memória ===========")
	for i := 0; i < TOTAL_SIZE; i++ {
		fmt.Printf("%3x:%3x ", i, memory[i])
		if i%16 == 15 {
			fmt.Println()
		}
	}
}