package items

import (
	"testing"

	"github.com/goccmack/gogll/v3/ast"
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

	New(g)
}
