package lex6

import (
	"testing"

	"github.com/goccmack/gogll/v3/test/lex/lex6/lexer"
)

var src = []string{"name"}
var tokType = []string{"id"}

func Test1(t *testing.T) {
	for i, s := range src {
		lxr := lexer.New([]rune(s))
		if lxr.Tokens[0].TypeID() != tokType[i] {
			t.Errorf("%s: %s", lxr.Tokens[0].LiteralString(), lxr.Tokens[0].TypeID())
		}
	}
}
