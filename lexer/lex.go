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

func (l *Lexer) add(t token.Type, lext, rext int) {
	l.Tokens = append(l.Tokens,
		token.New(t, lext, rext, l.I[lext:rext]))
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
		case '|':
			l.add(token.Type7, lext, l.pos)
		case '"':
			l.scanStringLiteral(lext)
		default:
			switch {
			case unicode.IsLower(l.I[lext]):
				l.scanTokIDOrReserevedWord(lext)
			case unicode.IsUpper((l.I[lext])):
				l.scanNT(lext)
			default:
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

func (l *Lexer) scanTokIDOrReserevedWord(lext int) {
	for l.isIDChar() {
		l.pos++
	}
	switch string(l.I[lext:l.pos]) {
	case "empty":
		l.add(token.Type2, lext, l.pos)
	case "package":
		l.add(token.Type4, lext, l.pos)
	default:
		l.add(token.Type6, lext, l.pos)
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

func (l *Lexer) skipWhiteSpace() {
	for l.pos < len(l.I) && unicode.IsSpace(l.I[l.pos]) {
		l.pos++
	}
}
