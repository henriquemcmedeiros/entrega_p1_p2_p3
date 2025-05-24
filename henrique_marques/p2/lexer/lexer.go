package lexer

import "unicode"

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL        TokenType = "ILLEGAL"
	EOF            TokenType = "EOF"
	IDENT          TokenType = "IDENT"
	INT_LITERAL    TokenType = "INT_LIT"
	FLOAT_LITERAL  TokenType = "FLOAT_LIT"
	STRING_LITERAL TokenType = "STRING_LIT"
	BOOL_LITERAL   TokenType = "BOOL_LIT"

	INICIO  = "INICIO"
	FIM     = "FIM"
	FUNC    = "FUNC"
	IF      = "IF"
	ELSE    = "ELSE"
	WHILE   = "WHILE"
	RETURN  = "RETURN"
	PRINT   = "PRINT"
	TYPE    = "TYPE"

	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="
	LTE      = "<="
	GTE      = ">="

	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
)

var keywords = map[string]TokenType{
	"inicio": INICIO,
	"fim":    FIM,
	"func":   FUNC,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"return": RETURN,
	"print":  PRINT,
	"int":    TYPE,
	"float":  TYPE,
	"string": TYPE,
	"bool":   TYPE,
	"true":   BOOL_LITERAL,
	"false":  BOOL_LITERAL,
}

type Lexer struct {
	input        string
	position     int  
	readPosition int  
	ch           byte 
}


func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}


func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}


func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}


func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(prev) + string(l.ch)}
		} else {
			tok = Token{Type: ASSIGN, Literal: string(l.ch)}
		}
	case '!':
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(prev) + string(l.ch)}
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	case '<':
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(prev) + string(l.ch)}
		} else {
			tok = Token{Type: LT, Literal: string(l.ch)}
		}
	case '>':
		if l.peekChar() == '=' {
			prev := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(prev) + string(l.ch)}
		} else {
			tok = Token{Type: GT, Literal: string(l.ch)}
		}
	case '+':
		tok = Token{Type: PLUS, Literal: string(l.ch)}
	case '-':
		tok = Token{Type: MINUS, Literal: string(l.ch)}
	case '*':
		tok = Token{Type: ASTERISK, Literal: string(l.ch)}
	case '/':
		tok = Token{Type: SLASH, Literal: string(l.ch)}
	case ',':
		tok = Token{Type: COMMA, Literal: string(l.ch)}
	case ';':
		tok = Token{Type: SEMICOLON, Literal: string(l.ch)}
	case '(':
		tok = Token{Type: LPAREN, Literal: string(l.ch)}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(l.ch)}
	case '{':
		tok = Token{Type: LBRACE, Literal: string(l.ch)}
	case '}':
		tok = Token{Type: RBRACE, Literal: string(l.ch)}
	case '[':
		tok = Token{Type: LBRACKET, Literal: string(l.ch)}
	case ']':
		tok = Token{Type: RBRACKET, Literal: string(l.ch)}
	case '"':
		tok.Type = STRING_LITERAL
		tok.Literal = l.readString()
		return tok
	case 0:
		tok = Token{Type: EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			ttype, ok := keywords[lit]
			if ok {
				tok = Token{Type: ttype, Literal: lit}
			} else {
				tok = Token{Type: IDENT, Literal: lit}
			}
			return tok
		} else if isDigit(l.ch) {
			lit, dt := l.readNumber()
			tok = Token{Type: dt, Literal: lit}
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(l.ch)}
		}
	}

	l.readChar()
	return tok
}


func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}


func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}


func (l *Lexer) readNumber() (string, TokenType) {
	start := l.position
	ttype := INT_LITERAL
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		ttype = FLOAT_LITERAL
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[start:l.position], ttype
}


func (l *Lexer) readString() string {
	l.readChar()
	start := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	lit := l.input[start:l.position]
	l.readChar()
	return lit
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) && ch >= 'a' && ch <= 'z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
