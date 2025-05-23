
# Projeto de Compilador, Assembler e Encoder para a Linguagem Neander

**Autor:** Henrique Marques de Carvalho Medeiros

## Descrição

1. **Compilador**: Transforma arquivos `.ldh` (linguagem do Henrique) em código Assembly `.asm`.
2. **Assembler**: Converte o código Assembly `.asm` para formato binário `.mem` compatível com o simulador NEANDER.
3. **Encoder**: Interpreta o conteúdo do arquivo `.mem` e exibe os resultados da execução (registradores e memória).

## Instruções de Uso

Certifique-se de estar na raiz do projeto e que o Go esteja corretamente instalado em sua máquina.

### 1. Compilar um programa `.ldh` para `.asm`
```bash
go run cmd/compiler/main.go io/linguagemCriada/program.ldh
```

### 2. Montar o arquivo `.asm` em um `.mem`
```bash
go run cmd/assembler/main.go io/asm/output.asm
```

### 3. Executar o programa `.mem` no emulador
```bash
go run cmd/encoder/main.go io/build/output.mem
```

## Limitações Conhecidas

- **Divisão**: A operação de divisão ainda não está implementada.
- **Expressões compostas**: Atualmente não é possível utilizar mais de uma variável para compor uma nova variável (ex: `X = A + B` ainda não é suportado, porém `X = A + 4` funciona).
- **Sem verificação de overflow**: O sistema não detecta ou trata estouro de valores no acumulador.
