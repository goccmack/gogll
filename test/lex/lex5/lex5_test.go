package lex5

import (
	"testing"

	"github.com/goccmack/gogll/test/lex/lex5/lexer"
	"github.com/goccmack/gogll/test/lex/lex5/parser"
)

var src = []string{"\"", "\\"}

func Test1(t *testing.T) {
	for _, s := range src {
		lxr := lexer.New([]rune(s))
		_, errs := parser.Parse(lxr)
		if errs != nil {
			t.Fatalf("Error: parsing \"%s\" - %s", s, errs[0])
		}
	}
}
