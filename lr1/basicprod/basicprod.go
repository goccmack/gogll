package basicprod

import (
	"fmt"
	"strings"

	"github.com/goccmack/gogll/ast"
)

type Production struct {
	Alternate int
	Head      string
	Body      *ast.SyntaxAlternate
}

func Get(rules []*ast.SyntaxRule) (prods []*Production) {
	prods = append(prods,
		&Production{
			Head: "G0",
			Body: &ast.SyntaxAlternate{
				Symbols: []ast.SyntaxSymbol{rules[0].Head},
			},
		})
	for _, r := range rules {
		for i, alt := range r.Alternates {
			prods = append(prods,
				&Production{
					Alternate: i,
					Head:      r.Head.ID(),
					Body:      alt,
				})
		}
	}
	return prods
}

func (p *Production) String() string {
	return fmt.Sprintf("%s : %s ;",
		p.Head,
		strings.Join(p.Body.GetSymbols(), " "))
}
