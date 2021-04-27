package item

import (
	"fmt"
	"testing"

	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/lex/item/pos"
	"github.com/goccmack/gogll/v3/lexer"
	"github.com/goccmack/gogll/v3/parser"
	"github.com/goccmack/gogll/v3/parser/bsr"
)

const src = `package "names"
qualifiedName : letter {letter|number|'_'} <'.' <letter|number|'_'>> ;
`

func Test1(t *testing.T) {
	lex := lexer.New([]rune(src))
	parser.Parse(lex)
	g := ast.Build(bsr.GetRoot(), lex)

	it := &Item{g.LexRules[0], pos.From([]int{2, 0, 1, 0, 1})}
	// it := &Item{g.LexRules[0], pos.From([]int{2, 0, 2})}
	fmt.Println(it)
	fmt.Println("Closure:")
	set := it.Emoves()
	for _, it := range set {
		fmt.Println(it)
	}
}
