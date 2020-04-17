// Package da implements disambiguation of the BSR forest
package da

import (
	"unicode"

	"github.com/goccmack/gogll/parser/bsr"
)

func Go() {
	goGLL(bsr.GetRoot())
}

// Alternate
//     :   Symbols
//     |   "empty"
//     ;
func alternate(b bsr.BSR) {
	if b.Alternate() == 0 {
		symbols(b.GetNTChildI(0))
	} // else do nothing
}

// Alternates
//     :   Alternate
//     |   Alternate "|" Alternates
//     ;
func alternates(b bsr.BSR) {
	alternate(b.GetNTChildI(0))
	if b.Alternate() == 1 {
		alternates(b.GetNTChildI(2))
	}
}

// GoGLL : Package Rules ;
func goGLL(b bsr.BSR) {
	rules(b.GetNTChildI(1))
}

// NT : id  ;
func nt(b bsr.BSR) {
	id := b.GetTChildI(0)
	if !unicode.IsUpper([]rune(id.Literal)[0]) {
		b.Ignore()
	}
}

// Rule : NT ":" Alternates ";"  ;
func rule(b bsr.BSR) {
	nt(b.GetNTChildI(0))
	alternates(b.GetNTChildI(2))
}

// Rules
//     :   Rule
//     |   Rule Rules
//     ;
func rules(b bsr.BSR) {
	rule(b.GetNTChildI(0))
	if b.Alternate() == 1 {
		rules(b.GetNTChildI(1))
	}
}

//  Symbol : NT | TokID | string_lit | "any" | "not" ;
func symbol(b bsr.BSR) {
	if b.Alternate() >= 2 {
		return
	}
	switch b.Alternate() {
	case 0:
		nt(b.GetNTChildI(0))
	case 1:
		tokID(b.GetNTChildI(0))
	}
	if len(b.GetNTChildrenI(0)) == 0 {
		b.Ignore()
	}
}

// Symbols
//     :   Symbol
//     |   Symbol Symbols
//     ;
func symbols(b bsr.BSR) {
	for _, sym := range b.GetNTChildrenI(0) {
		symbol(sym)
	}
	if b.Alternate() == 1 {
		symbols(b.GetNTChildI(1))
	}
}

// TokID : id ;
func tokID(b bsr.BSR) {
	id := b.GetTChildI(0)
	if !unicode.IsLower([]rune(id.Literal)[0]) {
		b.Ignore()
	}
}
