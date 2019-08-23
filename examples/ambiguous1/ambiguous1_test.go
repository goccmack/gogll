package ambiguous1

import (
	"fmt"
	"gogll/examples/ambiguous1/goutil/bsr"
	"gogll/examples/ambiguous1/parser"
	"gogll/examples/ambiguous1/parser/symbols"
	"testing"
)

func Test1(t *testing.T) {
	parser.Parse([]byte("aba"))
	da()
	bsr.DumpSlots()
}

func da() {
	rmBa()
	rmBereavedParents()
}

func rmBa() {
	for _, b := range bsr.GetAll() {
		fmt.Println(b)
		if b.GetString() == "a" && b.Label.Head() == "B" {
			fmt.Println("  ignore")
			b.Ignore()
		}
	}
}

func rmBereavedParents() {
	for again := true; again; {
		again = false
		for _, b := range bsr.GetAll() {
			for i, c := range b.Label.Symbols() {
				if symbols.IsNonTerminal(c) && nil == b.GetNTChildrenI(i) {
					b.Ignore()
					again = true
				}
			}
		}
	}
}
