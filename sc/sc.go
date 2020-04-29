// Package sc implements semantic checks on the target grammar
package sc

import (
	"fmt"
	"os"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lexer"
)

func Go(g *ast.GoGLL, l *lexer.Lexer) {
	checkNTRefs(g, l)
}

func checkNTRefs(g *ast.GoGLL, l *lexer.Lexer) {
	for _, r := range g.SyntaxRules {
		for _, alt := range r.Alternates {
			for _, sym := range alt.Symbols {
				switch s := sym.(type) {
				case *ast.NT:
					if nil == g.GetSyntaxRule(s.ID()) {
						fail(l, sym.Lext(), "No declaration of syntax rule %s", s.ID())
					}
				case *ast.TokID:
					if nil == g.GetLexRule(s.ID()) {
						fail(l, sym.Lext(), "No declaration of lex rule %s", s.ID())
					}
				}
			}
		}
	}
}

func fail(l *lexer.Lexer, pos int, format string, params ...interface{}) {
	ln, col := l.GetLineColumn(pos)
	msg := fmt.Sprintf(format, params...)
	fmt.Printf("Semantic Error at line %d col %d: %s\n", ln, col, msg)
	os.Exit(1)
}
