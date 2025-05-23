package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"p1/pkg/compiler/generator"
	"p1/pkg/compiler/lexer"
	"p1/pkg/compiler/parser"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Uso: go run cmd/compiler/main.go <arquivo.lfh> (exemplo: io/linguagemCriada/program.ldh)")
	}

	inputFile := os.Args[1]
	conteudo, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}

	tokens, err := lexer.Lex(string(conteudo))
	if err != nil {
		log.Fatalf("Erro léxico: %v", err)
	}

	parser := parser.NewParser(tokens)
	instrucoes, err := parser.ParsePrograma()
	if err != nil {
		log.Fatalf("Erro de parsing: %v", err)
	}

	prog := generator.GenerateASM(instrucoes)

	output := strings.Join(append(prog.Code, prog.Data...), "\n")

	err = os.WriteFile("io/asm/output.asm", []byte(output), 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar arquivo .asm: %v", err)
	}

	fmt.Println("Arquivo output.asm gerado com sucesso na pasta io/asm/output.asm")
}
