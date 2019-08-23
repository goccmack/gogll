/*
Package sa performs semantic analysis on the BSR and builds the AST
*/
package sa

import (
	"fmt"
	"gogll/ast"
	"gogll/goutil/bsr"
	"gogll/parser"
)

// Semantic analysis
type sa struct {
	grammar *ast.Grammar
	errors  []*SemanticError
}

type SemanticError struct {
	Line  int
	Col   int
	Error error
}

func (e *SemanticError) String() string {
	return fmt.Sprintf("Semantic Error: %s at line %d col %d", e.Error, e.Line, e.Col)
}

func Go() (*ast.Grammar, []*SemanticError) {
	sa := &sa{
		grammar: ast.New(),
	}
	sa.gogll(bsr.GetRoot())
	if len(sa.errors) != 0 {
		return nil, sa.errors
	}
	return sa.grammar, nil
}

/*
*GoGLL : SepE Package Sep Rules SepE ;
 */
func (a *sa) gogll(b bsr.BSR) {
	a.pckage(b.GetNTChild("Package", 0))
	a.rules(b.GetNTChild("Rules", 0))
	a.checkAllNTsDeclared()
}

/*
Package : "package" Sep String ;
*/
func (a *sa) pckage(b bsr.BSR) {
	pkg := b.GetNTChild("String", 0)
	a.grammar.AddPackage(pkg.GetString())
}

/*
Rules
    :   Rule
    |   Rule SepE Rules
    ;
*/
func (a *sa) rules(b bsr.BSR) {
	b1 := b
	for b1.Alternate() == 1 {
		a.rule(b1.GetNTChild("Rule", 0))
		b1 = b1.GetNTChild("Rules", 0)
	}
	a.rule(b1.GetNTChild("Rule", 0))
}

/*
Rule : Head SepE ":" SepE Alternates SepE ";" ;
*/
func (a *sa) rule(b bsr.BSR) {
	r := &ast.Rule{
		Head:       a.getHead(b.GetNTChild("Head", 0)),
		Alternates: a.getAlternates(b.GetNTChild("Alternates", 0)),
	}
	if err := a.grammar.AddRule(r); err != nil {
		a.Error(b, err)
	}
}

/*
Head : "*" NonTerminal | NonTerminal ;
*/
func (a *sa) getHead(b bsr.BSR) *ast.Head {
	nt := b.GetNTChild("NonTerminal", 0).GetString()
	if b.Alternate() == 0 {
		if err := a.grammar.SetStartSymbol(nt); err != nil {
			a.Error(b, err)
		}
	}
	return ast.NewHead(b, nt)
}

/*
Alternates
    :   Alternate
    |   Alternate SepE "|" SepE Alternates
    ;
*/
func (a *sa) getAlternates(b bsr.BSR) (alts []*ast.Alternate) {
	b1 := b
	for b1.Alternate() == 1 {
		alts = append(alts, a.getAlternate(b1.GetNTChild("Alternate", 0)))
		b1 = b1.GetNTChild("Alternates", 0)
	}
	alts = append(alts, a.getAlternate(b1.GetNTChild("Alternate", 0)))
	return
}

/*
Alternate
    :   Symbols
    |   "empty"
    ;

*/
func (a *sa) getAlternate(b bsr.BSR) *ast.Alternate {
	alt := &ast.Alternate{}
	if b.Alternate() == 1 {
		return alt
	}
	alt.Symbols = a.symbols(b.GetNTChild("Symbols", 0))
	return alt
}

/*
Symbols
    :   Symbol Sep Symbols
    |   Symbol
    ;
*/
func (a *sa) symbols(b bsr.BSR) (symbols []ast.Symbol) {
	// line, col := parser.GetLineColumn(b.LeftExtent())
	// fmt.Printf("sa.symbols(%s) line %d col %d\n", b, line, col)
	// for _, b1 := range b.GetNTChildren("Symbol", 0) {
	// 	fmt.Println("", b1)
	// }
	b1 := b
	for b1.Alternate() == 0 {
		symbols = append(symbols, a.symbol(b1.GetNTChild("Symbol", 0))...)
		b1 = b1.GetNTChild("Symbols", 0)
	}
	symbols = append(symbols, a.symbol(b1.GetNTChild("Symbol", 0))...)
	return
}

/*
Symbol : NonTerminal | Terminal ;
*/
func (a *sa) symbol(b bsr.BSR) []ast.Symbol {
	if b.Alternate() == 0 {
		return []ast.Symbol{ast.NewNonTerminal(b, b.GetNTChild("NonTerminal", 0).GetString())}
	}
	return a.terminal(b.GetNTChild("Terminal", 0))
}

/*
Terminal
    :   "any"
    |   "anyof" Sep String
    |   "letter"
    |   "number"
    |   "space"
    |   "upcase"
    |   "lowcase"
    |   "not" Sep String
    |   CharLiteral
    |   String
    ;
*/
func (a *sa) terminal(b bsr.BSR) []ast.Symbol {
	switch b.Alternate() {
	case 0:
		return []ast.Symbol{ast.NewAny(b)}
	case 1:
		return []ast.Symbol{a.grammar.NewAnyOf(b, b.GetNTChild("String", 0).GetString())}
	case 2:
		return []ast.Symbol{ast.NewLetter(b)}
	case 3:
		return []ast.Symbol{ast.NewNumber(b)}
	case 4:
		return []ast.Symbol{ast.NewSpace(b)}
	case 5:
		return []ast.Symbol{ast.NewUpcase(b)}
	case 6:
		return []ast.Symbol{ast.NewLowcase(b)}
	case 7:
		return []ast.Symbol{a.grammar.NewNot(b, b.GetNTChild("String", 0).GetString())}
	case 8:
		return []ast.Symbol{a.grammar.NewCharLiteral(b, b.GetNTChild("CharLiteral", 0).GetString())}
	case 9:
		return a.grammar.NewString(b, b.GetNTChild("String", 0).GetString())
	}
	panic(fmt.Sprintf("Invalid alternate %d", b.Alternate()))
}

/*** Error ***/

func (sa *sa) Error(b bsr.BSR, err error) {
	se := &SemanticError{
		Error: err,
	}
	se.Line, se.Col = parser.GetLineColumn(b.LeftExtent())
	sa.errors = append(sa.errors, se)
}

/*** checks ***/

func (a *sa) checkAllNTsDeclared() {
	usedNT := map[string]bool{}
	for _, rule := range a.grammar.Rules {
		for _, alt := range rule.Alternates {
			for _, sym := range alt.Symbols {
				if nt, ok := sym.(*ast.NonTerminal); ok {
					usedNT[nt.NT] = true
					if nil == a.grammar.GetRule(string(nt.NT)) {
						a.Error(nt.GetBSR(), fmt.Errorf("No grammar rule %s", nt.NT))
					}
				}
			}
		}
	}
	for _, rule := range a.grammar.Rules {
		if _, used := usedNT[rule.Head.NT]; !used && rule.Head.NT != a.grammar.StartSymbol {
			fmt.Printf("sa.checkAllNTsDeclared: NT=%s, grammer.S=%s\n", rule.Head.NT, a.grammar.StartSymbol)
			a.Error(rule.Head.BSR, fmt.Errorf("Rule %s is not used", rule.Head.NT))
		}
	}
}
