package main

import (
    "fmt"
    "log"
    "os"

    "app/lexer"
)

func main() {
    // Lê todo o conteúdo de code.ldh
    data, err := os.ReadFile("code.ldh")
    if err != nil {
        log.Fatalf("erro ao ler o arquivo code.ldh: %v", err)
    }
    sourceCode := string(data)

    // Inicializa o lexer com o conteúdo do arquivo
    l := lexer.New(sourceCode)

    // Itera sobre os tokens até EOF
    for tok := l.NextToken(); tok.Type != lexer.EOF; tok = l.NextToken() {
        fmt.Printf("%+v\n", tok)
    }
}

