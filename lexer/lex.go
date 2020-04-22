package lexer

import (
	"unicode"

	"github.com/goccmack/gogll/token"
)

func New(input []rune) *Lexer {
	lex := &Lexer{
		I:      input,
		Tokens: make([]*token.Token, 0, 2048),
		pos:    0,
	}
	lex.scan()
	lex.add(token.EOF, len(input), len(input))
	return lex
}

func (l *Lexer) scan() {
	for l.pos < len(l.I) {
		l.skipWhiteSpace()
		if l.pos >= len(l.I) {
			return
		}

		lext := l.pos
		l.pos++
		switch l.I[lext] {
		case ':':
			l.add(token.Type0, lext, l.pos)
		case ';':
			l.add(token.Type1, lext, l.pos)
		case '"':
			l.scanStringLiteral(lext)
		case '|':
			l.add(token.Type6, lext, l.pos)
		default:
			if unicode.IsLetter(l.I[lext]) {
				l.scanIDOrReservedWord(lext)
			} else {
				l.add(token.Error, lext, l.pos)
			}
		}
	}
}

func (l *Lexer) scanNT(lext int) {
	for l.isIDChar() {
		l.pos++
	}
	l.add(token.Type3, lext, l.pos)
}

func (l *Lexer) scanIDOrReservedWord(lext int) {
	for l.isIDChar() {
		l.pos++
	}
	switch string(l.I[lext:l.pos]) {
	case "empty":
		l.add(token.Type2, lext, l.pos)
	case "package":
		l.add(token.Type4, lext, l.pos)
	default: // id
		l.add(token.Type3, lext, l.pos)
	}
}

func (l *Lexer) isIDChar() bool {
	if l.pos >= len(l.I) {
		return false
	}
	c := l.I[l.pos]
	return unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_'
}

func (l *Lexer) scanStringLiteral(lext int) {
	for l.pos < len(l.I) && l.I[l.pos] != '"' {
		l.pos++
	}
	if l.pos >= len(l.I) || l.I[l.pos] != '"' {
		l.add(token.Error, lext, l.pos)
	} else {
		l.pos++
		l.add(token.Type5, lext, l.pos)
	}
}
