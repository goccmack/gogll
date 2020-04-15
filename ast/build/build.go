// Package build builds the Abstract Syntax Tree from a disambiguated parse BSR forest.
package build

import (
	"fmt"
	"os"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lexer"
	"github.com/goccmack/gogll/parser/bsr"
	"github.com/goccmack/gogll/parser/symbols"
	"github.com/goccmack/gogll/token"
	"github.com/goccmack/goutil/stringset"
)

type builder struct {
	lex *lexer.Lexer
}

// From builds an AST from the BSR root. `root` is the root of a disambiguated BSR forest
//
// GoGLL : Package Rules ;
func From(root bsr.BSR, l *lexer.Lexer) *ast.GoGLL {
	bld := &builder{lex: l}
	gogll := &ast.GoGLL{
		Package: bld.packge(root.GetNTChild(symbols.NT_Package, 0)),
		Rules:   bld.rules(root.GetNTChild(symbols.NT_Rules, 0)),
	}
	gogll.NonTerminals = bld.nonTerminals(gogll.Rules)
	gogll.Terminals = bld.terminals(gogll.Rules)
	return gogll
}

// Alternate
//     :   Symbols
//     |   "empty"
//     ;
func (bld *builder) alternate(b bsr.BSR) *ast.Alternate {
	alt := &ast.Alternate{}
	if b.Alternate() == 0 {
		alt.Symbols = bld.symbols(b.GetNTChildI(0))
	} // if alt = empty return alt with empty Symbols
	return alt
}

// Alternates
//     :   Alternate
//     |   Alternate "|" Alternates
//     ;
func (bld *builder) alternates(b bsr.BSR) []*ast.Alternate {
	alts := []*ast.Alternate{
		bld.alternate(b.GetNTChild(symbols.NT_Alternate, 0)),
	}
	if b.Alternate() == 1 {
		alts = append(alts, bld.alternates(b.GetNTChild(symbols.NT_Alternates, 0))...)
	}
	return alts
}

func (bld *builder) nonTerminals(rules []*ast.Rule) *stringset.StringSet {
	nts := stringset.New()
	for _, r := range rules {
		if nts.Contain(r.Head.Token()) {
			bld.fail(fmt.Errorf("Duplicate rule %s", r.Head.Token()), r.Head.Lext())
		} else {
			nts.Add(r.Head.Token())
		}
	}
	return nts
}

// NT : nt  ;
func (bld *builder) nt(b bsr.BSR) *ast.NT {
	tok := b.GetTChildI(0)
	if tok.Type != token.StringToType["nt"] {
		bld.fail(fmt.Errorf("expected non-terminal ID"), b.LeftExtent())
	}
	return &ast.NT{
		Tok: tok,
	}
}

// Package : "package" string_lit ;
func (bld *builder) packge(b bsr.BSR) *ast.Package {
	return &ast.Package{
		Tok: b.GetTChildI(1),
	}
}

// Rule : NT ":" Alternates ";"  ;
func (bld *builder) rule(b bsr.BSR) *ast.Rule {
	return &ast.Rule{
		Head:       bld.nt(b.GetNTChild(symbols.NT_NT, 0)),
		Alternates: bld.alternates(b.GetNTChild(symbols.NT_Alternates, 0)),
	}
}

// Rules
//     :   Rule
//     |   Rule Rules
//     ;
func (bld *builder) rules(b bsr.BSR) []*ast.Rule {
	rules := []*ast.Rule{bld.rule(b.GetNTChildI(0))}
	if b.Alternate() == 1 {
		rules = append(rules, bld.rules(b.GetNTChildI(1))...)
	}
	return rules
}

func (bld *builder) stringLit(tok *token.Token) *ast.StringLit {
	if tok.Type != token.StringToType["string_lit"] {
		bld.fail(
			fmt.Errorf("expected string_lit but got %s", token.TypeToString[tok.Type]),
			tok.Lext)
	}
	return &ast.StringLit{
		Tok: tok,
	}
}

// Symbol : NT | tokid | string_lit ;
func (bld *builder) symbol(b bsr.BSR) ast.Symbol {
	switch b.Alternate() {
	case 0:
		return bld.nt(b.GetNTChildI(0))
	case 1:
		return bld.tokID(b.GetTChildI(0))
	case 2:
		return bld.stringLit(b.GetTChildI(0))
	}
	panic(fmt.Sprintf("invalid alternate %d", b.Alternate()))
}

// Symbols
//     :   Symbol
//     |   Symbol Symbols
//     ;
func (bld *builder) symbols(b bsr.BSR) []ast.Symbol {
	symbols := []ast.Symbol{bld.symbol(b.GetNTChildI(0))}
	if b.Alternate() == 1 {
		symbols = append(symbols, bld.symbols(b.GetNTChildI(1))...)
	}
	return symbols
}

func (bld *builder) terminals(rules []*ast.Rule) *stringset.StringSet {
	terminals := stringset.New()
	for _, r := range rules {
		for _, a := range r.Alternates {
			for _, s := range a.Symbols {
				switch t := s.(type) {
				case *ast.TokID, *ast.StringLit:
					terminals.Add(t.Token())
				}
			}
		}
	}
	return terminals
}

func (bld *builder) tokID(tok *token.Token) *ast.TokID {
	if tok.Type != token.StringToType["tokid"] {
		bld.fail(
			fmt.Errorf("expected tokid but got %s", token.TypeToString[tok.Type]),
			tok.Lext)
	}
	return &ast.TokID{Tok: tok}
}

/*** Errors ***/

// i is the position of the failure in input slice of runes
func (bld *builder) fail(err error, i int) {
	ln, col := bld.lex.GetLineColumn(i)
	fmt.Printf("Error: %s at line %d col %d\n", err, ln, col)
	os.Exit(1)
}
