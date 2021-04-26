//  Copyright 2020 Marius Ackerman
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

package parser

import (
	"github.com/goccmack/gogll/v3/ast"
	"github.com/goccmack/gogll/v3/gslot"
)

func (g *gen) getAlternates() (alts []*Alternate) {
	for _, nt := range g.g.NonTerminals.ElementsSorted() {
		rule := g.g.GetSyntaxRule(nt)
		for i, alt := range rule.Alternates {
			alts = append(alts, g.getAlternate(nt, alt, i))
		}
	}
	return
}

func (g *gen) getAlternate(nt string, alt *ast.SyntaxAlternate, altI int) *Alternate {
	// fmt.Printf("codex.getAltData %s[%d]\n", nt, altI)
	L := gslot.NewLabel(nt, altI, 0, g.gs, g.ff)
	d := &Alternate{
		NT:       nt,
		AltLabel: L.Label(),
		Comment:  L.String(),
		Empty:    alt.Empty(),
	}
	if !alt.Empty() {
		d.Slots = g.getSlotsData(nt, alt, altI)
		d.LastSlot = d.Slots[len(d.Slots)-1]
	}
	return d
}

func (g *gen) getSlotsData(nt string, alt *ast.SyntaxAlternate, altI int) (data []*SlotData) {
	for i, sym := range alt.Symbols {
		// fmt.Printf("getSlotsData(%s) %s\n", nt, getSlotData(nt, altI, sym, i))
		data = append(data, g.getSlotData(nt, altI, sym.String(), i))
	}
	return
}

func (g *gen) getSlotData(nt string, altI int, symbol string, pos int) *SlotData {
	preLabel := gslot.NewLabel(nt, altI, pos, g.gs, g.ff)
	postLabel := gslot.NewLabel(nt, altI, pos+1, g.gs, g.ff)
	sd := &SlotData{
		AltLabel:  gslot.NewLabel(nt, altI, 0, g.gs, g.ff).Label(),
		PreLabel:  preLabel.Label(),
		PostLabel: postLabel.Label(),
		Comment:   postLabel.String(),
		Head:      nt,
	}
	sd.IsNT = !g.g.Terminals.Contain(symbol)
	// fmt.Printf("getSlotData: altlabel:%s, pre:%s, post:%s\n",
	// 	sd.AltLabel, sd.PreLabel, sd.PostLabel)
	return sd
}
