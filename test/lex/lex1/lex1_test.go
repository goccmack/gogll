package lex1

import (
	"testing"

	"github.com/goccmack/gogll/v3/test/lex/lex1/lexer"
)

type Test struct {
	Input      []rune
	TokenTypes []string
}

// rule1: ['a'] 'b' ;
var test1 = []Test{
	{Input: []rune(" a "), TokenTypes: []string{"Error", "EOF"}},
	{Input: []rune(" ab "), TokenTypes: []string{"rule1", "EOF"}},
	{Input: []rune("    "), TokenTypes: []string{"EOF"}},
	{Input: []rune(" bb "), TokenTypes: []string{"rule1", "rule1", "EOF"}},
	{Input: []rune(" ab c "), TokenTypes: []string{"rule1", "Error", "EOF"}},
	{Input: []rune(" c c c "), TokenTypes: []string{"Error", "Error", "Error", "EOF"}},
}

func Test1(t *testing.T) {
	for i, tst := range test1 {
		l := lexer.New(tst.Input)
		if len(l.Tokens) != len(tst.TokenTypes) {
			t.Logf("%3d: Expected: %s", i, tst.TokenTypes)
			t.Logf("     But received: %s", l.Tokens)
		} else {
			for i := range l.Tokens {
				if tst.TokenTypes[i] != l.Tokens[i].TypeID() {
					t.Logf("%3d: Expected: %s", i, tst.TokenTypes)
					t.Logf("     But received: %s", l.Tokens)
				}
			}
		}
	}
}
