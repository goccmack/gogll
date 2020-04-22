// Package tokens is the intermediate representation of tokens for all code generators.
package tokens

import (
	"fmt"

	"github.com/goccmack/gogll/ast"
)

type Tokens struct {
	TypeToString    []string
	StringToType    map[string]int
	LiteralToString map[string]string
	TypeToLiteral   []string
}

func New(g *ast.GoGLL) *Tokens {
	tokens := &Tokens{
		TypeToString:    []string{},
		StringToType:    map[string]int{},
		LiteralToString: map[string]string{},
	}
	tokens.add("Error", "Error")
	tokens.add("EOF", "EOF")
	for i, tok := range g.Terminals.ElementsSorted() {
		tokens.add(fmt.Sprintf("Type%d", i), tok)
	}
	return tokens
}

func (t *Tokens) add(tokID, tokLit string) {
	tokType := len(t.TypeToString)
	t.StringToType[tokID] = tokType
	t.TypeToString = append(t.TypeToString, tokID)
	t.TypeToLiteral = append(t.TypeToLiteral, tokLit)
	t.LiteralToString[tokLit] = tokID
}
