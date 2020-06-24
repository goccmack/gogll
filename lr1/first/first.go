package first

import (
	"bytes"
	"fmt"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/lr1/basicprod"
	"github.com/goccmack/gogll/symbols"
)

type First struct {
	// key: symbol
	symbol map[string]FirstSet

	// key: prod id
	ntMap  map[string][]*basicprod.Production
	ntList []string

	tList []string
}

func New(prods []*basicprod.Production) *First {
	first := &First{
		symbol: make(map[string]FirstSet),
		ntMap:  make(map[string][]*basicprod.Production),
		ntList: symbols.GetNonTerminalSymbols(),
		tList:  symbols.GetTerminalSymbols(),
	}
	first.AddTerminals()
	first.AddNonTerminals(prods)
	return first
}

func (this *First) AddNonTerminals(prods []*basicprod.Production) {
	for again := true; again; {
		again = false
		for _, prod := range prods {
			changed := false
			this.symbol[prod.Head], changed =
				this.symbol[prod.Head].Add(this.firstTerms(prod.Body)...)
			if changed {
				again = true
			}
		}
	}
}

func (this *First) AddTerminals() {
	for _, s := range this.tList {
		if _, exist := this.symbol[s]; exist {
			panic(fmt.Sprintf("Duplicate terminal: %s", s))
		}
		this.symbol[s] = FirstSet{s}
	}
}

func (this *First) FirstSymbol(sym string) FirstSet {
	return this.symbol[sym]
}

func (this *First) firstTerms(alt *ast.SyntaxAlternate) (first FirstSet) {
	for _, t := range alt.Symbols {
		switch term := t.(type) {
		case *ast.TokID:
			first, _ = first.Add(term.ID())
			return
		case *ast.NT:
			if f := this.symbol[term.ID()]; !f.Contain("ℇ") {
				first, _ = first.Add(f...)
				return
			} else {
				for _, sym := range f {
					if sym != "ℇ" {
						first, _ = first.Add(sym)
					}
				}
			}
		case *ast.StringLit:
			first, _ = first.Add(string(term.Value()))
			return
		default:
			panic("Invalid")
		}
	}
	first, _ = first.Add("ℇ")
	return
}

func (this *First) FirstString(s []string, context ...string) (fs FirstSet) {
	if len(s) == 0 {
		fs, _ = fs.Add(context...)
		return
	}
	frst := make(FirstSet, 0, 4)
	deriveEmpty := true
	for i := 0; deriveEmpty && i < len(s); deriveEmpty, i = frst.Contain("ℇ"), i+1 {
		frst = this.FirstSymbol(s[i])
		fs, _ = fs.Add(frst.Min("ℇ")...)
	}
	if deriveEmpty {
		fs, _ = fs.Add(context...)
	}
	return
}

func (this *First) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "Terminals:\n")
	for _, t := range this.tList {
		fmt.Fprintf(w, "\t%s\t%s\n", t, this.symbol[t])
	}
	fmt.Fprintf(w, "\nNon-terminals:\n")
	for _, nt := range this.ntList {
		fmt.Fprintf(w, "\t%s\t%s\n", nt, this.symbol[nt])
	}
	return w.String()
}
