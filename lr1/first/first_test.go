package first

import (
	"fmt"
	"testing"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/lexer"
	"github.com/goccmack/gogll/v3/lr1/basicprod"
	"github.com/goccmack/gogll/v3/parser"
)

const src = `
package "testx"
A : B | "a" ;
B : "b" ; 
C : B "c" ;
`

func TestFirst1(t *testing.T) {
	lex := lexer.New([]rune(src))
	bsr, err := parser.Parse(lex)
	if err != nil {
		for _, e := range err {
			fmt.Printf("error: %#v\n", e.String())
		}
		panic(err)
	}
	g := ast.Build(bsr.GetRoot(), lex, "test.md")
	basicProds := basicprod.Get(g.SyntaxRules)
	first := New(basicProds)
	exp := FirstSet{"b"}
	f := first.FirstString([]string{"B"}, "x", "y")
	if !f.Equal(exp) {
		t.Errorf("Expected %s, got %s", exp, f)
	}

	exp = FirstSet{"x", "y"}
	f = first.FirstString(nil, "x", "y")
	if !f.Equal(exp) {
		t.Errorf("Expected %s, got %s", exp, f)
	}
}
