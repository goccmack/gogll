package items

import (
	"fmt"
	"testing"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/lexer"
	"github.com/goccmack/gogll/v3/parser"
)

const src = `package "names"
qualifiedName : letter {letter|number|'_'} <'.' <letter|number|'_'>> ;
`

func Test1(t *testing.T) {
	lex := lexer.New([]rune(src))
	bsr, err := parser.Parse(lex)
	if err != nil {
		for _, e := range err {
			fmt.Printf("error: %#v\n", e.String())
		}
		panic(err)
	}
	g := ast.Build(bsr.GetRoot(), lex, "test.md")

	New(g)
}
