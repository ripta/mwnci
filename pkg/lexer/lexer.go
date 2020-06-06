package lexer

import (
	"github.com/ripta/mwnci/pkg/ast"
	"github.com/ripta/mwnci/pkg/token"
)

type Lexer struct {
	input []rune
	pos   int
	read  int
	rune  rune

	lnum int
	cnum int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: []rune(input),
	}

	l.readRune()
	return l
}

func (l *Lexer) GetLineCol() (int, int) {
	return l.lnum + 1, l.cnum
}

func (l *Lexer) GetLocation() ast.Location {
	return ast.Location{
		Line: l.lnum + 1,
		Col:  l.cnum,
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.rune {
	case '=':
		if l.peekRune() == '=' {
			r := l.rune
			l.readRune()
			tok = token.Token{
				Type:    token.EQ,
				Literal: string(r) + string(l.rune),
			}
		} else {
			tok = newToken(token.ASSIGN, l.rune)
		}
	case '+':
		tok = newToken(token.PLUS, l.rune)
	case '-':
		tok = newToken(token.MINUS, l.rune)
	case '!':
		if l.peekRune() == '=' {
			r := l.rune
			l.readRune()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: string(r) + string(l.rune),
			}
		} else {
			tok = newToken(token.BANG, l.rune)
		}
	case '/':
		tok = newToken(token.SLASH, l.rune)
	case '*':
		tok = newToken(token.ASTERISK, l.rune)
	case '<':
		tok = newToken(token.LT, l.rune)
	case '>':
		tok = newToken(token.GT, l.rune)
	case ';':
		tok = newToken(token.SEMICOLON, l.rune)
	case ':':
		tok = newToken(token.COLON, l.rune)
	case ',':
		tok = newToken(token.COMMA, l.rune)
	case '{':
		tok = newToken(token.LBRACE, l.rune)
	case '}':
		tok = newToken(token.RBRACE, l.rune)
	case '(':
		tok = newToken(token.LPAREN, l.rune)
	case ')':
		tok = newToken(token.RPAREN, l.rune)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.rune)
	case ']':
		tok = newToken(token.RBRACKET, l.rune)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.rune) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.Lookup(tok.Literal)
			return tok
		} else if isDigit(l.rune) {
			intPart := l.readNumber()
			if l.rune == '.' && isDigit(l.peekRune()) {
				l.readRune()
				fracPart := l.readNumber()
				tok.Type = token.FLOAT
				tok.Literal = intPart + "." + fracPart
			} else {
				tok.Type = token.INT
				tok.Literal = intPart
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.rune)
		}
	}

	l.readRune()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.rune == ' ' || l.rune == '\t' || l.rune == '\n' || l.rune == '\r' {
		l.readRune()
	}
}

func (l *Lexer) readRune() {
	if l.read >= len(l.input) {
		l.rune = 0
	} else {
		l.rune = l.input[l.read]
	}
	l.pos = l.read
	l.read += 1

	l.cnum += 1
	if l.rune == '\n' {
		l.cnum = 0
		l.lnum += 1
	}
}

func (l *Lexer) peekRune() rune {
	if l.read >= len(l.input) {
		return 0
	} else {
		return l.input[l.read]
	}
}

func (l *Lexer) readIdentifier() string {
	p := l.pos
	for isLetter(l.rune) {
		l.readRune()
	}
	return string(l.input[p:l.pos])
}

func (l *Lexer) readNumber() string {
	p := l.pos
	for isDigit(l.rune) {
		l.readRune()
	}
	return string(l.input[p:l.pos])
}

func (l *Lexer) readString() string {
	p := l.pos + 1
	for {
		l.readRune()
		if l.rune == '"' || l.rune == 0 {
			break
		}
	}
	return string(l.input[p:l.pos])
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func newToken(tokenType token.Type, r rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(r)}
}
