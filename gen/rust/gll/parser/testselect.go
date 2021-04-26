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
	"fmt"
	"os"
	"sort"

	"github.com/goccmack/gogll/v3/frstflw"
	"github.com/goccmack/gogll/v3/gslot"
	"github.com/goccmack/gogll/v3/symbols"
)

func (g *gen) getTSData() (data []*TSData) {
	for _, s := range g.gs.Slots() {
		data = append(data, g.getSlotTSData(s))
	}
	return
}

func (g *gen) getSlotTSData(l gslot.Label) *TSData {
	return &TSData{
		Label:   l.Label(),
		Comment: l.String(),
		Symbols: g.getFirst(l),
	}
}

func (g *gen) getFirst(l gslot.Label) (tokTypes []*Symbol) {
	// fmt.Printf("testSelect.getFirst(%s)\n", l)
	ss := l.Symbols()[l.Pos:]
	frst := g.ff.FirstOfString(ss.Strings())
	firstSymbols := frst.Elements()
	sort.Slice(
		firstSymbols,
		func(i, j int) bool { return firstSymbols[i] < firstSymbols[j] })
	// fmt.Printf("  first: %s\n", frst)
	for _, sym := range firstSymbols {
		if sym != frstflw.Empty {
			tokTypes = append(tokTypes,
				&Symbol{
					TokenType: symbols.TerminalLiteralToType(sym).TypeString(),
					Comment:   sym,
				})
		}
	}
	if frst.Contain(frstflw.Empty) {
		tokTypes = append(tokTypes, g.getFollowConditions(l.Head)...)
	}
	return
}

func (g *gen) getFollowData() (data []*TSData) {
	for _, nt := range g.g.NonTerminals.ElementsSorted() {
		data = append(data, g.getFollowDataForNT(nt))
	}
	return
}

func (g *gen) getFollowDataForNT(nt string) *TSData {
	d := &TSData{
		Label:   nt,
		Comment: nt,
		Symbols: g.getFollowConditions(nt),
	}
	return d
}

func (g *gen) getFollowConditions(nt string) (tokens []*Symbol) {
	// fmt.Printf("testselect.getFollowConditions(%s)\n", nt)
	flw := g.ff.Follow(nt)
	if flw.Len() == 0 {
		fmt.Printf("Production %s has empty follow set. It is never called\n", nt)
		os.Exit(1)
	}
	for _, sym := range flw.ElementsSorted() {
		tokens = append(tokens,
			&Symbol{
				TokenType: symbols.TerminalLiteralToType(sym).TypeString(),
				Comment:   sym,
			})
	}
	return
}
