//  Copyright 2019 Marius Ackerman
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

// Package gslot implements grammar slots
package gslot

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/goccmack/gogll/ast"
	"github.com/goccmack/gogll/frstflw"
	"github.com/goccmack/gogll/symbols"
)

type Label struct {
	Head      string
	Alternate int
	Pos       int
	gs        *GSlot
	ff        *frstflw.FF
}

type Slots []Label

type GSlot struct {
	g     *ast.GoGLL
	ff    *frstflw.FF
	slots map[Label]symbols.Symbols
}

func New(g *ast.GoGLL, ff *frstflw.FF) *GSlot {
	gs := &GSlot{
		g:     g,
		ff:    ff,
		slots: make(map[Label]symbols.Symbols),
	}
	gs.genSlots()
	return gs
}

func NewLabel(head string, alt, pos int, gs *GSlot, ff *frstflw.FF) *Label {
	return &Label{
		Head:      head,
		Alternate: alt,
		Pos:       pos,
		gs:        gs,
		ff:        ff,
	}
}

func (gs *GSlot) Slots() Slots {
	res := make(Slots, 0, len(gs.slots))
	for l, _ := range gs.slots {
		res = append(res, l)
	}
	sort.Sort(res)
	return res
}

func (s Label) Label() string {
	return fmt.Sprintf("%s%dR%d", s.Head, s.Alternate, s.Pos)
}

func (s Label) IsEoR() bool {
	symbols := s.gs.slots[s]
	return s.Pos >= len(symbols)
}

func (s Label) IsFiR() bool {
	symbols := s.gs.slots[s]
	if s.Pos > 1 || len(symbols) <= 1 {
		return false
	}
	if s.ff.FirstOfSymbol(symbols[0].Literal()).Contain(frstflw.Empty) &&
		symbols[0] == symbols[1] {
		return false
	}
	return true
}

func (s Label) String() string {
	symbols := s.gs.slots[s]
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s : ", s.Head)
	for i, sym := range symbols {
		if i == s.Pos {
			fmt.Fprintf(buf, "∙")
		}
		fmt.Fprintf(buf, "%s ", sym)
	}
	// fmt.Printf("slotLabel.String(): %s pos=%d len(symbols)=%d\n", s.Head, s.Pos, len(symbols))
	if s.Pos >= len(symbols) {
		fmt.Fprintf(buf, "∙")
	}
	// fmt.Println("  ", buf.String())
	return buf.String()
}

func (s Label) Symbols() symbols.Symbols {
	return s.gs.slots[s]
}

func (ss Slots) Labels() (labels []string) {
	for _, s := range ss {
		labels = append(labels, s.Label())
	}
	return
}

func (ss Slots) Len() int {
	return len(ss)
}

func (ss Slots) Less(i, j int) bool {
	if ss[i].Head < ss[j].Head {
		return true
	}

	if ss[i].Head == ss[j].Head {
		if ss[i].Alternate < ss[j].Alternate {
			return true
		}
		if ss[i].Alternate == ss[j].Alternate {
			if ss[i].Pos < ss[j].Pos {
				return true
			}
		}
	}
	return false
}

func (ss Slots) Swap(i, j int) {
	iTmp := ss[i]
	ss[i] = ss[j]
	ss[j] = iTmp
}

func (gs *GSlot) genSlots() {
	for _, rule := range gs.g.SyntaxRules {
		gs.genSlotsOfRule(rule)
	}
}

func (gs *GSlot) genSlotsOfRule(r *ast.SyntaxRule) {
	for i, a := range r.Alternates {
		gs.genSlotsOfAlternate(r.Head.ID(), i, getSymbols(a.GetSymbols())...)
	}
}

func (gs *GSlot) genSlotsOfAlternate(nt string, altI int, symbls ...symbols.Symbol) {
	if len(symbls) == 0 {
		gs.genSlot(nt, altI, 0, []symbols.Symbol{}...)
	} else {
		for pos := 0; pos <= len(symbls); pos++ {
			gs.genSlot(nt, altI, pos, symbls...)
		}
	}
}

func (gs *GSlot) genSlot(nt string, altI, pos int, symbols ...symbols.Symbol) {
	slot := Label{
		Head:      nt,
		Alternate: altI,
		Pos:       pos,
		gs:        gs,
		ff:        gs.ff,
	}
	gs.slots[slot] = symbols
}

// getSymbols translates AST symbol strings to symbols.Symbol
func getSymbols(astSymbols []string) []symbols.Symbol {
	symbls := make([]symbols.Symbol, len(astSymbols))
	for i, s := range astSymbols {
		symbls[i] = symbols.FromASTString(s)
	}
	return symbls
}
