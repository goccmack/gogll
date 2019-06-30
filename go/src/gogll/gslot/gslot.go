/*
Package gslot implements grammar slots
*/
package gslot

import (
	"bytes"
	"fmt"
	"gogll/ast"
	"sort"
)

type SlotLabel struct {
	Head      string
	Alternate int
	Pos       int
}

var slots = make(map[SlotLabel]Symbols)

type Symbols []string

type Slots []SlotLabel

func GetSlots() Slots {
	if len(slots) == 0 {
		genSlots()
	}
	res := make(Slots, 0, len(slots))
	for l, _ := range slots {
		res = append(res, l)
	}
	sort.Sort(res)
	return res
}

func (s SlotLabel) Label() string {
	return fmt.Sprintf("%s%dR%d", s.Head, s.Alternate+1, s.Pos+1)
}

func (s SlotLabel) IsEoR() bool {
	symbols := slots[s]
	return s.Pos == len(symbols)
}

func (s SlotLabel) IsFiR() bool {
	return s.Pos == 1
}

func (s SlotLabel) String() string {
	symbols := slots[s]
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s : ", s.Head)
	for i, sym := range symbols {
		if i == s.Pos {
			fmt.Fprintf(buf, "∙")
		}
		fmt.Fprintf(buf, "%s ", sym)
	}
	if s.Pos == len(symbols) {
		fmt.Fprintf(buf, "∙")
	}
	return buf.String()
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
		if ss[i].Pos < ss[j].Pos {
			return true
		}
	}
	return false
}

func (ss Slots) Swap(i, j int) {
	iTmp := ss[i]
	ss[i] = ss[j]
	ss[j] = iTmp
}

func genSlots() {
	for _, rule := range ast.GetRules() {
		genSlotsOfRule(rule)
	}
}

func genSlotsOfRule(r *ast.Rule) {
	for i, a := range r.Alternates {
		genSlotsOfAlternate(r.Head.Value(), i, a.Symbols()...)
	}
}

func genSlotsOfAlternate(nt string, altI int, symbols ...string) {
	for pos := range symbols {
		genSlot(nt, altI, pos+1, symbols...)
	}
}

func genSlot(nt string, altI, pos int, symbols ...string) {
	slot := SlotLabel{
		Head:      nt,
		Alternate: altI,
		Pos:       pos,
	}
	slots[slot] = symbols
}
