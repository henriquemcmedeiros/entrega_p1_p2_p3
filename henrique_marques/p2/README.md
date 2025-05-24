# Linguagem LDH - Analisador Léxico em Go

Autor: Henrique Marques de Carvalho Medeiros

Este projeto é um analisador léxico escrito em Go para a linguagem definida pela gramática `bnfgramatica.txt`. O programa lê um arquivo `.ldh` contendo código fonte e imprime a sequência de tokens reconhecidos.

## Estrutura

- `main.go`: ponto de entrada do programa. Lê o arquivo `code.ldh`, inicializa o lexer e imprime os tokens.
- `lexer.go`: implementação do analisador léxico (lexer), responsável por identificar tokens válidos da linguagem.
- `bnfgramatica.txt`: define a gramática da linguagem LDH em formato BNF.

## Executando o projeto

1. Certifique-se de ter o Go instalado em seu sistema.
2. Execute o programa com:

```bash
go run main.go
```

Isso irá processar o conteúdo do arquivo `code.ldh` e imprimir os tokens reconhecidos no terminal.

## Exemplo de uso

Dado o arquivo `code.ldh` com o seguinte conteúdo:

```ldh
inicio
int x;
x = 42;
fim
```

A saída será uma lista dos tokens identificados, como:

```
{Type:INICIO Literal:inicio}
{Type:TYPE Literal:int}
{Type:IDENT Literal:x}
{Type:SEMICOLON Literal:;}
{Type:IDENT Literal:x}
{Type:ASSIGN Literal:=}
{Type:INT_LIT Literal:42}
{Type:SEMICOLON Literal:;}
{Type:FIM Literal:fim}
```